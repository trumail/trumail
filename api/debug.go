package api

import (
	"bytes"
	"net/http"
	"runtime/pprof"

	"github.com/labstack/echo"
	"github.com/sdwolfe32/trumail/config"
)

// Debug is an endpoint for debugging runaway goroutines
func Debug(c echo.Context) error {
	if c.Request().Header.Get("X-Auth-Token") != config.Token {
		return c.JSON(http.StatusUnauthorized, nil)
	}
	var buf bytes.Buffer
	pprof.Lookup("goroutine").WriteTo(&buf, 1)
	return c.String(http.StatusOK, buf.String())
}
