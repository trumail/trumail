package api

import (
	"net/http"

	"github.com/labstack/echo"
)

// ErrVerificationFailure is thrown when there is error while validating an email
var ErrVerificationFailure = echo.NewHTTPError(http.StatusInternalServerError, "Failed to perform email verification lookup")

// Lookup performs a single email verification and returns a fully
// populated lookup or an error
func (s *Service) Lookup(c echo.Context) error {
	l := s.Logger.WithField("handler", "Lookup")
	l.Debug("New Lookup request received")

	// Decode the email from the request
	l.Debug("Decoding the request")
	email := c.Param("email")
	l = l.WithField("email", email)

	// Performs the full email verification
	l.Debug("Performing new email verification")
	lookup, err := s.Verifier.VerifyTimeout(email, s.Timeout)
	if err != nil {
		l.WithError(err).Error("Failed to perform verification")
		return s.Encode(c, http.StatusInternalServerError, err)
	}
	l = l.WithField("lookup", lookup)

	// Returns the email validation lookup to the requestor
	l.Debug("Returning Email Lookup")
	return s.Encode(c, http.StatusOK, lookup)
}
