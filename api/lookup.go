package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/sdwolfe32/trumail/verifier"
)

// LookupHandler performs a single email verification and returns
// a fully populated lookup or an error
func LookupHandler(v *verifier.Verifier) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract the request variables
		email := c.Param("email")
		ts, _ := strconv.Atoi(c.QueryParam("timeout"))

		// Perform the unlimited verification
		if ts == 0 {
			lookup, err := v.Verify(email)
			if err != nil {
				return FormatEncoder(c, http.StatusInternalServerError, err)
			}
			return FormatEncoder(c, http.StatusOK, lookup)
		}

		// Parse the timeout and perform the limited lookup
		timeout := time.Duration(ts) * time.Second
		lookup, err := v.VerifyTimeout(email, timeout)
		if err != nil {
			return FormatEncoder(c, http.StatusInternalServerError, err)
		}
		return FormatEncoder(c, http.StatusOK, lookup)
	}
}
