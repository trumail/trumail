package api

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/sdwolfe32/trumail/verifier"
)

// LookupHandler performs a single email verification and returns
// a fully populated lookup or an error
func LookupHandler(v *verifier.Verifier) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Perform the unlimited verification
		lookup, err := v.Verify(c.Param("email"))
		if err != nil {
			return FormatEncoder(c, http.StatusInternalServerError, err)
		}
		return FormatEncoder(c, http.StatusOK, lookup)
	}
}
