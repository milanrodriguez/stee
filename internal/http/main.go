package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/milanrodriguez/stee/internal/stee"
)

func handleMain(core *stee.Core) echo.HandlerFunc {
	return func(c echo.Context) error {
		key := c.Request().URL.Path

		key = strings.TrimPrefix(key, "/")
		key = strings.TrimSuffix(key, "/")

		target, err := core.GetRedirection(key)
		if err != nil {
			c.Response().WriteHeader(http.StatusNotFound)
			_, err = fmt.Fprintf(c.Response().Writer, "Error 404: No redirection found for key \"%s\"", key)
			return err
		}
		return c.Redirect(http.StatusFound, target)
	}
}
