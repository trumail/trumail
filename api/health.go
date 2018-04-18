package api

import (
	"net/http"

	"github.com/labstack/echo"
)

// Status is a type reserved for holding service healths state
type Status string

// StatusOK indicates the service is running and is healthy
const StatusOK Status = "OK"

// Health is a healthcheck response body
type Health struct {
	Status Status `json:"status"`
}

// Health returns a Health check response indicating the
// health state of the service
func (t *TrumailAPI) Health(c echo.Context) error {
	l := t.log.WithField("handler", "health")
	l.Debug("New Health check request received")

	// Return a new Health check reference
	l.Debug("Returning Health check Response")
	return c.JSON(http.StatusOK, &Health{
		Status: StatusOK,
	})
}
