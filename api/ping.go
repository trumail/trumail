package api

import (
	"net/http"

	"github.com/labstack/echo"
)

// PingRes is an overall stats summary
type PingRes struct {
	Result string `json:"result"`
}

// Ping retrieves and returns general Trumail statistics
func (t *TrumailAPI) Ping(c echo.Context) error {
	var s PingRes
	s.Result = "PONG"
	return c.JSON(http.StatusOK, s)
}
