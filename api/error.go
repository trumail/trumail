package api

import (
	"net/http"

	"github.com/labstack/echo"
)

// ErrorHandler is a custom error handler that will
// redirect all 404s back to https://trumail.io
func ErrorHandler(err error, c echo.Context) {
	if e, ok := err.(*echo.HTTPError); ok &&
		e.Code == http.StatusNotFound {
		c.Redirect(http.StatusFound, "https://trumail.io")
		return
	}
	c.Echo().DefaultHTTPErrorHandler(err, c)
}
