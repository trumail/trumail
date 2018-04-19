package api

import (
	"net/http"

	raven "github.com/getsentry/raven-go"
	"github.com/labstack/echo"
	tinystat "github.com/sdwolfe32/tinystat/client"
	"github.com/sdwolfe32/trumail/verifier"
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
func (t *TrumailAPI) Lookup(c echo.Context) error {
	l := t.log.WithField("handler", "Lookup")
	l.Debug("New Lookup request received")

	// Decode the request
	l.Debug("Decoding the request")
	email := c.Param("email")
	l = l.WithField("email", email)

	// Check cache for a successful lookup
	if lookup, ok := t.lookupCache.Get(email); ok {
		return t.encodeResponse(c, http.StatusOK, lookup)
	}

	// Performs the full email verification
	l.Debug("Performing new email verification")
	lookup, err := t.verifier.VerifyTimeout(email, t.timeout)
	if err != nil {
		l.WithError(err).Error("Failed to perform verification")
		if le, ok := err.(*verifier.LookupError); ok {
			return t.encodeResponse(c, http.StatusInternalServerError, le)
		}
		if err.Error() == verifier.ErrEmailParseFailure {
			return t.encodeResponse(c, http.StatusBadRequest, err)
		}
		return t.encodeResponse(c, http.StatusInternalServerError, err)
	}
	l = l.WithField("lookup", lookup)

	// Store the lookup in cache
	t.lookupCache.SetDefault(email, lookup)

	// Returns the email validation lookup to the requestor
	l.Debug("Returning Email Lookup")
	return t.encodeResponse(c, http.StatusOK, lookup)
}

// encodeResponse encodes the passed response using the "format" and
// "callback" parameters on the passed echo.Context
func (t *TrumailAPI) encodeResponse(c echo.Context, code int, res interface{}) error {
	// Send metrics of successful response
	if le, ok := res.(*verifier.Lookup); ok {
		if le.Deliverable {
			tinystat.CreateAction("deliverable")
		} else {
			tinystat.CreateAction("undeliverable")
		}
	}

	// Send metrics of error response
	if e, ok := res.(error); ok {
		if le, ok := e.(*verifier.LookupError); ok {
			// LookupError with report == true
			if le.Report {
				raven.CaptureError(e, nil) // Sentry metrics
				tinystat.CreateAction("error")
			}
		} else {
			// Standard error
			raven.CaptureError(e, nil) // Sentry metrics
			tinystat.CreateAction("error")
		}
	}

	// Encode the in requested format
	switch c.Param("format") {
	case "json":
		return c.JSON(code, res)
	case "jsonp":
		callback := c.QueryParam("callback")
		if callback == "" {
			return ErrInvalidCallback
		}
		return c.JSONP(code, callback, res)
	case "xml":
		return c.XML(code, res)
	default:
		return ErrUnsupportedFormat
	}
}
