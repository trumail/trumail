package api

import (
	"net/http"
	"strings"

	raven "github.com/getsentry/raven-go"
	"github.com/labstack/echo"
	tinystat "github.com/sdwolfe32/tinystat/client"
	"github.com/sdwolfe32/trumail/heroku"
	"github.com/sdwolfe32/trumail/verifier"
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
	// ErrVerificationFailure is thrown when there is error while validating an email
	ErrVerificationFailure = echo.NewHTTPError(http.StatusInternalServerError, "Failed to perform email verification lookup")
	// ErrUnsupportedFormat is thrown when the requestor has defined an unsupported response format
	ErrUnsupportedFormat = echo.NewHTTPError(http.StatusBadRequest, "Unsupported format")
	// ErrInvalidCallback is thrown when the request is missing the callback queryparam
	ErrInvalidCallback = echo.NewHTTPError(http.StatusBadRequest, "Invalid callback query param provided")
)

// Lookup performs a single email verification and returns a fully
// populated lookup or an error
func (s *Service) Lookup(c echo.Context) error {
	l := s.log.WithField("handler", "Lookup")
	l.Debug("New Lookup request received")

	// Decode the email from the request
	l.Debug("Decoding the request")
	email := c.Param("email")
	l = l.WithField("email", email)

	// Performs the full email verification
	l.Debug("Performing new email verification")
	lookup, err := s.verifier.VerifyTimeout(email, s.timeout)
	if err != nil {
		if strings.Contains(err.Error(), verifier.ErrBlocked) {
			// Restart Dyno if officially confirmed blacklisted
			if err := s.verifier.Blacklisted(); err != nil {
				l.WithError(err).Warn("Confirmed Blacklisted! - Restarting Dyno")
				go l.Info(heroku.RestartDyno())
			}
		}
		l.WithError(err).Error("Failed to perform verification")
		return countAndRespond(c, http.StatusInternalServerError, err)
	}
	l = l.WithField("lookup", lookup)

	// Returns the email validation lookup to the requestor
	l.Debug("Returning Email Lookup")
	return countAndRespond(c, http.StatusOK, lookup)
}

// countAndRespond encodes the passed response using the "format" and
// "callback" parameters on the passed echo.Context
func countAndRespond(c echo.Context, code int, res interface{}) error {
	count(res)                   // Submit metrics data
	return respond(c, code, res) // Encode the response
}

// count calls out to the various metrics APIs we have set up in order
// to submit metrics data based on the response
func count(res interface{}) {
	switch r := res.(type) {
	case *verifier.Lookup:
		if r.Deliverable {
			tinystat.CreateAction("deliverable")
		} else {
			tinystat.CreateAction("undeliverable")
		}
	case error:
		raven.CaptureError(r, nil) // Sentry metrics
		tinystat.CreateAction("error")
	}
}

// respond writes the status code and response in the desired
// format to the ResponseWriter using the passed echo.Context
func respond(c echo.Context, code int, res interface{}) error {
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
