package api

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

const (
	// FormatJSON is the format constant for a JSON output
	FormatJSON = "json"
	// FormatJSONP is the format constant for a JSONP output
	FormatJSONP = "jsonp"
	// FormatXML is the format constant for a XML output
	FormatXML = "xml"
)

var (
	// ErrInvalidCallback is thrown when the request is missing the callback queryparam
	ErrInvalidCallback = echo.NewHTTPError(http.StatusBadRequest, "Invalid callback query param provided")
	// ErrUnsupportedFormat is thrown when the requestor has defined an unsupported response format
	ErrUnsupportedFormat = echo.NewHTTPError(http.StatusBadRequest, "Unsupported format")
)

// Encoder is a function type that encodes a response given a
// context, a status code and a response
type Encoder func(c echo.Context, code int, res interface{}) error

// DefaultEncoder is an encoder that reads the format from the
// passed echo context and writes the status code and response
// based on that format
func DefaultEncoder(c echo.Context, code int, res interface{}) error {
	// Encode the in requested format
	switch strings.ToLower(c.Param("format")) {
	case FormatXML:
		return c.XML(code, res)
	case FormatJSON:
		return c.JSON(code, res)
	case FormatJSONP:
		callback := c.QueryParam("callback")
		if callback == "" {
			return ErrInvalidCallback
		}
		return c.JSONP(code, callback, res)
	default:
		return ErrUnsupportedFormat
	}
}
