package api

import (
	"net/http"
	"strings"

	raven "github.com/getsentry/raven-go"
	"github.com/labstack/echo"
	tinystat "github.com/sdwolfe32/tinystat/client"
	"github.com/sdwolfe32/trumail/verifier"
)

const (
	FormatJSON  = "JSON"
	FormatJSONP = "JSONP"
	FormatXML   = "XML"
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

	// Decode the email from the request
	l.Debug("Decoding the request")
	email := c.Param("email")
	l = l.WithField("email", email)

	// Parse the address passed
	address, err := verifier.ParseAddress(email)
	if err != nil {
		return t.encodeRes(c, http.StatusBadRequest, err)
	}

	// Check cache for a successful lookup
	if lookup, ok := t.lookupCache.Get(address.MD5Hash); ok {
		return t.encodeRes(c, http.StatusOK, lookup)
	}

	// Performs the full email verification
	l.Debug("Performing new email verification")
	lookup, err := t.verifier.VerifyAddressTimeout(address, t.timeout)
	if err != nil {
		l.WithError(err).Error("Failed to perform verification")
		return t.encodeRes(c, http.StatusInternalServerError, err)
	}
	l = l.WithField("lookup", lookup)

	// Store the lookup in cache
	t.lookupCache.SetDefault(address.MD5Hash, lookup)

	// Returns the email validation lookup to the requestor
	l.Debug("Returning Email Lookup")
	return t.encodeRes(c, http.StatusOK, lookup)
}

// encodeRes encodes the passed response using the "format" and
// "callback" parameters on the passed echo.Context
func (t *TrumailAPI) encodeRes(c echo.Context, code int, res interface{}) error {
	// Submit metrics data
	count(res)

	// Encode the in requested format
	switch strings.ToUpper(c.Param("format")) {
	case FormatJSON:
		return c.JSON(code, res)
	case FormatJSONP:
		callback := c.QueryParam("callback")
		if callback == "" {
			return ErrInvalidCallback
		}
		return c.JSONP(code, callback, res)
	case FormatXML:
		return c.XML(code, res)
	default:
		return ErrUnsupportedFormat
	}
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
		if le, ok := r.(*verifier.LookupError); ok {
			// LookupError with report == true
			if le.Report {
				raven.CaptureError(r, nil) // Sentry metrics
				tinystat.CreateAction("error")
			}
		} else {
			// Standard error
			raven.CaptureError(r, nil) // Sentry metrics
			tinystat.CreateAction("error")
		}
	}
}
