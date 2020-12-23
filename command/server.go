package command

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	steehttp "github.com/milanrodriguez/stee/internal/http"
	"github.com/milanrodriguez/stee/internal/stee"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCommand.AddCommand(serverCommand)
}

var serverCommand = &cobra.Command{
	Use:     "server",
	Short:   "Start the Stee server",
	Long:    `Start the Stee server`,
	Run:     serverRun,
	Aliases: []string{"serve", "srv"},
}

// ServerConfig is the configuration for the Stee server
// struct members need to be exported for Viper to unmarshal configuration
type serverConfig struct {
	Address string
	Port    int
	TLS     struct {
		Enable   bool
		CertPath string
		KeyPath  string
	}
	API struct {
		Enable        bool
		URLPathPrefix string
		SimpleAPI     struct {
			Enable bool
		}
	}
	UI struct {
		Enable        bool
		URLPathPrefix string
	}
}

func serverRun(cmd *cobra.Command, args []string) {
	// Create configuration
	globalConfig, err := loadConfig()
	config := globalConfig.Server
	if err != nil {
		panic(fmt.Errorf("cannot load configuration: %v", err))
	}

	// Create core
	core, err := stee.NewCore(
		stee.Store(viper.Sub("storage")),
	)
	if err != nil {
		panic(fmt.Errorf("cannot initialize Stee: %v", err))
	}
	_ = core.AddRedirectionWithKey("_stee", "https://github.com/milanrodriguez/stee")

	httpHandler := steehttp.HandleRoot(core,
		steehttp.EnableAPI(config.API.Enable, config.API.URLPathPrefix),
		steehttp.EnableSimpleAPI(config.API.SimpleAPI.Enable),
		steehttp.EnableUI(config.UI.Enable, config.UI.URLPathPrefix),
	)

	srv := steehttp.NewServer(
		steehttp.ServerConfig{
			ListenAddress: config.Address + ":" + strconv.Itoa(config.Port),
			Handler:       httpHandler,
		},
	)
	// Start to listen.
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.ParseIP(config.Address),
		Port: config.Port,
	})
	if err != nil {
		panic(fmt.Errorf("could not start listening at %v: %v", srv.Addr, err))
	}
	fmt.Printf("✔️ Listening at %s\n", srv.Addr)
	var scheme string
	var serve func() error
	if config.TLS.Enable {
		serve = func() error { return srv.ServeTLS(listener, config.TLS.CertPath, config.TLS.KeyPath) }
		scheme = "https"
	} else {
		serve = func() error { return srv.Serve(listener) }
		scheme = "http"
		fmt.Printf("\n⚠️ You are running the server without TLS encryption!\n⚠️ You should consider setting up HTTPS for production use.\n\n")
	}
	// Start to serve.
	go func() {
		fmt.Printf("✔️ %s://%s/\n", scheme, srv.Addr)
		err := serve()
		if err != http.ErrServerClosed {
			panic(fmt.Errorf("server closed unexpectedly: %v", err))
		}
	}()

	//////////////////////////////////////////////////////////////////
	// At this point initialization is done, the server is running. //
	// We're now listening for OS sigint signal                     //
	//////////////////////////////////////////////////////////////////

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for range c {
		fmt.Printf("\n🛑 Interruption requested. We're gonna perform a clean shutdown...\n")

		// Shutting down http servers
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		err := srv.Shutdown(ctx)
		cancel()
		if err != nil {
			fmt.Printf("❌ problem while shutting down the http server: %v", err)
		}

		listener.Close()

		// Shutting down the core.
		err = core.Close()
		if err != nil {
			fmt.Printf("❌ problem while shutting down Stee: %v", err)
		}
		fmt.Printf("Bye!\n")

		close(c)
		os.Exit(0)
	}
}
