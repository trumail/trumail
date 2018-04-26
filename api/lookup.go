package api

import (
	"net/http"

	"github.com/labstack/echo"
)

// ErrVerificationFailure is thrown when there is error while
// validating an email
var ErrVerificationFailure = echo.NewHTTPError(http.StatusInternalServerError,
	"Failed to perform email verification lookup")

// Lookup performs a single email verification and returns a fully
// populated lookup or an error
func (s *Service) Lookup(c echo.Context) error {
	// Performs the full email verification
	lookup, err := s.Verifier.VerifyTimeout(c.Param("email"), s.Timeout)
	if err != nil {
		return s.Encode(c, http.StatusInternalServerError, err)
	}

	// Returns the email validation lookup to the requestor
	return s.Encode(c, http.StatusOK, lookup)
}
