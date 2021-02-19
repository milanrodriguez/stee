package http

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/milanrodriguez/stee/internal/stee"
)

func handleSimpleAdd(core *stee.Core) echo.HandlerFunc {
	return func(c echo.Context) error {
		targetBytes, err := base64.URLEncoding.DecodeString(c.Param("base64target"))
		if err != nil {
			// If URLEncoding doesn't, we use RawURLEncoding (omits padding characters)
			targetBytes, err = base64.RawURLEncoding.DecodeString(c.Param("base64target"))
		}
		if err != nil {
			c.Response().WriteHeader(http.StatusBadRequest)
			_, err = fmt.Fprintf(c.Response().Writer, "Error: %s", err.Error())
			return err
		}
		target := string(targetBytes)

		key := c.Param("key")
		if key != "" {
			err = core.AddRedirectionWithKey(key, target)
		} else {
			key, err = core.AddRedirectionWithoutKey(target)
		}
		if err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			_, err = fmt.Fprintf(c.Response().Writer, "Error: %s", err.Error())
			return err
		}
		_, err = fmt.Fprintf(c.Response().Writer, "✔️ Added redirection: '%s' -> %s", key, target)
		return err
	}
}

func handleSimpleGet(core *stee.Core) echo.HandlerFunc {
	return func(c echo.Context) error {
		key := c.Param("key")
		target, err := core.GetRedirection(key)
		if err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			_, err = fmt.Fprintf(c.Response().Writer, "Impossible to read key: %s", err)
			return err
		}
		_, err = fmt.Fprintf(c.Response().Writer, "Key \"%s\" is pointing to \"%s\"", key, target)
		return err
	}
}

func handleSimpleDel(core *stee.Core) echo.HandlerFunc {
	return func(c echo.Context) error {
		key := c.Param("key")
		err := core.DeleteRedirection(key)
		if err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			_, err = fmt.Fprintf(c.Response().Writer, "Error: %s", err.Error())
			return err
		}
		_, err = fmt.Fprintf(c.Response().Writer, "Deleted key \"%s\"", key)
		return err
	}
}
