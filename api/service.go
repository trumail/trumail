package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sdwolfe32/slimhttp"
	"github.com/sdwolfe32/trumail/verifier"
	"github.com/sirupsen/logrus"
)

// maxWorkerCount specifies a maximum number of goroutines allowed
// when processing bulk email lists (not a public endpoint yet)
const maxWorkerCount = 20

// TrumailService defines all functionality for the Trumail email
// verification API
type TrumailService interface {
	Lookup(r *http.Request) (interface{}, error)
}

// trumail contains all dependencies for a TrumailService
type trumail struct {
	log      *logrus.Entry
	hostname string
	verify   verifier.Verifier
}

// NewTrumailService generates a new NewTrumailService
func NewTrumailService(log *logrus.Logger, hostname, sourceAddr string) TrumailService {
	return &trumail{
		log:      log.WithField("service", "lookup"),
		hostname: hostname,
		verify:   verifier.NewVerifier(maxWorkerCount, hostname, sourceAddr),
	}
}

// Lookup performs a single email validation and returns a fully
// populated lookup or an error
func (t *trumail) Lookup(r *http.Request) (interface{}, error) {
	l := t.log.WithField("handler", "Lookup")
	l.Debug("New Lookup request received")

	// Decode the request
	l.Debug("Decoding the request")
	email := mux.Vars(r)["email"]
	l = l.WithField("email", email)

	// Verify required fields exist
	l.Debug("Verifying request has all required fields")
	if email == "" {
		return nil, slimhttp.NewError("No email found on request", http.StatusBadRequest, nil).Log(l)
	}

	// Performs the full email validation
	l.Debug("Performing new validation lookup")
	lookups := t.verify.Verify(email)
	if len(lookups) == 0 {
		return nil, slimhttp.NewError("Error validating email", http.StatusInternalServerError, nil).Log(l)
	}
	lookup := lookups[0]

	// Returns the email validation lookup to the requestor
	l.WithField("lookup", lookup).Debug("Returning Email Lookup")
	return lookup, nil
}
