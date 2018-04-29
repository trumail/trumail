package api

import (
	"net/http"

	"github.com/labstack/echo"
)

// Health is a healthcheck response body
type Health struct {
	Status   string `json:"status"`
	Hostname string `json:"hostname"`
}

// HealthHandler returns a HealthResponse indicating the
// current health state of the service
func HealthHandler(hostname string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, &Health{"OK", hostname})
	}
}
