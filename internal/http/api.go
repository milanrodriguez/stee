package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/milanrodriguez/stee/internal/stee"
)

func handleAPI(core *stee.Core) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().WriteHeader(http.StatusNotImplemented)
		_, err := fmt.Fprintf(c.Response().Writer, "Error 501: API not implemented in this version. Try to upgrade Stee.")
		return err
	}
}
