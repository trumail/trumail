package api

import (
	"net/http"

	"github.com/labstack/echo"
)

// Health is a healthcheck response body
type Health struct {
	Status string `json:"status"`
}

// Health returns a Health check response indicating the
// health state of the service
func (s *Service) Health(c echo.Context) error {
	l := s.log.WithField("handler", "health")
	l.Debug("New Health check request received")

	// Return a new Health check reference
	l.Debug("Returning Health check Response")
	return c.JSON(http.StatusOK, &Health{"OK"})
}
