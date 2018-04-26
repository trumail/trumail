package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/sdwolfe32/trumail/verifier"
)

// LookupHandler performs a single email verification and returns
// a fully populated lookup or an error
func LookupHandler(v *verifier.Verifier,
	timeout time.Duration) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Perform the full email verification and return the
		// response to the requestor
		lookup, err := v.VerifyTimeout(c.Param("email"), timeout)
		if err != nil {
			return FormatEncoder(c, http.StatusInternalServerError, err)
		}
		return FormatEncoder(c, http.StatusOK, lookup)
	}
}
