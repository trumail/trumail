package api

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/sdwolfe32/trumail/verifier"
)

const maxWorkerCount = 20

// Lookuper defines all functionality for an email verification
// lookup API service
type Lookuper interface {
	Lookup(r *http.Request) (interface{}, error)
	Healthcheck(r *http.Request) (interface{}, error)
}

// lookuper contains all dependencies for a Lookuper
type lookuper struct {
	log      *logrus.Entry
	hostname string
	ever     verifier.Verifier
}

// NewLookuper generates a new email verification lookup API service
func NewLookuper(log *logrus.Logger, hostname, sourceAddr string) Lookuper {
	return &lookuper{
		log:      log.WithField("service", "lookup"),
		hostname: hostname,
		ever:     verifier.NewVerifier(maxWorkerCount, hostname, sourceAddr),
	}
}

// Lookup performs a single email validation
func (s *lookuper) Lookup(r *http.Request) (interface{}, error) {
	l := s.log.WithField("handler", "Lookup")
	l.Info("New Lookup request received")

	// Decode the request
	l.Info("Decoding the request")
	email := mux.Vars(r)["email"]
	l = l.WithField("email", email)

	// Verify required fields exist
	l.Info("Verifying request has all required fields")
	if email == "" {
		return nil, NewError("No email found on request", http.StatusBadRequest, nil).Log(l)
	}

	// Performs the full email validation
	l.Info("Performing new validation lookup")
	lookups := s.ever.Verify(email)
	if len(lookups) == 0 {
		return nil, NewError("Error validating email", http.StatusInternalServerError, nil).Log(l)
	}
	lookup := lookups[0]

	// Returns the email validation lookup to the requestor
	l.WithField("lookup", lookup).Info("Returning Email Lookup")
	return lookup, nil
}

// healthcheck represents the response to a healthcheck request
type healthcheck struct {
	Status   string `json:"status" xml:"status"`
	Hostname string `json:"hostname" xml:"hostname"`
}

// GetHealthcheck handles and returns a 200 and our hostname
func (s *lookuper) Healthcheck(r *http.Request) (interface{}, error) {
	l := s.log.WithField("handler", "Healthcheck")
	l.Info("New Healthcheck request received")
	l.Info("Returning newly generated Healthcheck")
	return &healthcheck{Status: "OK", Hostname: s.hostname}, nil
}
