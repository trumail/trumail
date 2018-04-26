package api

import (
	"net/http"

	"github.com/labstack/echo"
)

// Health is a healthcheck response body
type Health struct {
	Status string `json:"status"`
}

// Healthcheck returns a Health response indicating the
// health state of the service
func Healthcheck(c echo.Context) error {
	return c.JSON(http.StatusOK, &Health{"OK"})
}
