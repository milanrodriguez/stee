package http

import (
	"github.com/labstack/echo/v4"
	"github.com/milanrodriguez/stee/internal/stee"
)

type Server interface {
	Start(address string) error
	Shutdown() error
}

type echoWrapper struct {
	echo *echo.Echo
}

type ServerConfig struct {
	Main struct{}
	API  struct {
		Enable    bool
		Prefix    string
		SimpleAPI struct {
			Enable bool
		}
	}
	UI struct {
		Enable bool
		Prefix string
	}
}

// NewServer returns a http.Server
func NewServer(core *stee.Core, config ServerConfig) Server {
	e := echo.New()
	e.HideBanner = true

	//route registration
	e.GET("/*", handleMain(core))

	e.GET(config.API.Prefix+"/*", handleAPI(core))
	e.GET(config.API.Prefix+"/simple/add/:base64target", handleSimpleAdd(core))
	e.GET(config.API.Prefix+"/simple/add/:base64target/:key", handleSimpleAdd(core))
	e.GET(config.API.Prefix+"/simple/get/:key", handleSimpleGet(core))
	e.GET(config.API.Prefix+"/simple/del/:key", handleSimpleDel(core))

	e.GET(config.UI.Prefix+"/*", handleUI(core))

	return &echoWrapper{e}
}

func (e *echoWrapper) Start(address string) error {
	return e.echo.Start(address)
}

func (e *echoWrapper) Shutdown() error {
	return e.echo.Close()
}
