package api

import (
	"net/http"

	"github.com/labstack/echo"
)

// Health is a healthcheck response body
type Health struct {
	Status string `json:"status"`
}

// HealthHandler returns a HealthResponse indicating the
// current health state of the service
func HealthHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, &Health{"OK"})
	}
}
