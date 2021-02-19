package command

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"

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

	srvConf := steehttp.ServerConfig{
		Main: struct{}{},
		API: struct {
			Enable    bool
			Prefix    string
			SimpleAPI struct{ Enable bool }
		}{
			Enable: config.API.Enable,
			Prefix: config.API.URLPathPrefix,
			SimpleAPI: struct{ Enable bool }{
				Enable: config.API.SimpleAPI.Enable,
			},
		},
		UI: struct {
			Enable bool
			Prefix string
		}{
			Enable: config.UI.Enable,
			Prefix: config.UI.URLPathPrefix,
		},
	}
	srv := steehttp.NewServer(core, srvConf)
	srv.Start(config.Address + ":" + strconv.Itoa(config.Port))

	//////////////////////////////////////////////////////////////////
	// At this point initialization is done, the server is running. //
	// We're now listening for OS sigint signal                     //
	//////////////////////////////////////////////////////////////////

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for range c {
		fmt.Printf("\n🛑 Interruption requested. We're gonna perform a clean shutdown...\n")

		// Shutting down http servers
		err := srv.Shutdown()
		if err != nil {
			fmt.Printf("❌ problem while shutting down the http server: %v", err)
		}

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
