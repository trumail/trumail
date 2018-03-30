package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo"

	"github.com/sdwolfe32/trumail/verifier"
	"github.com/sirupsen/logrus"
)

// maxWorkerCount specifies a maximum number of goroutines allowed
// when processing bulk email lists (not a public endpoint yet)
const maxWorkerCount = 20

var (
	// ErrValidationFailure indicates that there was an error while validating an email
	ErrValidationFailure = echo.NewHTTPError(http.StatusInternalServerError, "Error validating email")
	// ErrUnsupportedFormat indicates that the requestor has defined an unsupported response format
	ErrUnsupportedFormat = echo.NewHTTPError(http.StatusBadRequest, "Unsupported format")
	// ErrInvalidCallback indicates that the request is missing the callback queryparam
	ErrInvalidCallback = echo.NewHTTPError(http.StatusBadRequest, "Invalid callback query param provided")
)

// TrumailAPI contains all dependencies for the Trumail API
type TrumailAPI struct {
	log      *logrus.Entry
	hostname string
	verify   *verifier.Verifier
}

// NewTrumailAPI generates a new Trumail reference
func NewTrumailAPI(log *logrus.Logger, hostname, sourceAddr string) *TrumailAPI {
	return &TrumailAPI{
		log:      log.WithField("service", "lookup"),
		hostname: hostname,
		verify: verifier.NewVerifier(&http.Client{
			Timeout: time.Second,
		}, maxWorkerCount, hostname, sourceAddr),
	}
}

// Lookup performs a single email validation and returns a fully
// populated lookup or an error
func (t *TrumailAPI) Lookup(c echo.Context) error {
	l := t.log.WithField("handler", "Lookup")
	l.Debug("New Lookup request received")

	// Decode the request
	l.Debug("Decoding the request")
	email := c.Param("email")
	l = l.WithField("email", email)

	// Performs the full email validation
	l.Debug("Performing new validation lookup")
	lookups := t.verify.Verify(email)
	if len(lookups) == 0 {
		l.WithError(ErrValidationFailure).Error("Failed to validate email")
		return ErrValidationFailure
	}
	lookup := lookups[0]

	// Returns the email validation lookup to the requestor
	l.WithField("lookup", lookup).Debug("Returning Email Lookup")
	return t.encodeLookup(c, lookup)
}

// encodeLookup encodes the passed response using the "format" and
// "callback" parameters on the passed echo.Context
func (t *TrumailAPI) encodeLookup(c echo.Context, res interface{}) error {
	switch c.Param("format") {
	case "json":
		return c.JSON(http.StatusOK, res)
	case "jsonp":
		callback := c.QueryParam("callback")
		if callback == "" {
			return ErrInvalidCallback
		}
		return c.JSONP(http.StatusOK, callback, res)
	case "xml":
		return c.XML(http.StatusOK, res)
	default:
		return ErrUnsupportedFormat
	}
}
