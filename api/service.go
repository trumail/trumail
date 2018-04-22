package api

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/sdwolfe32/trumail/heroku"
	"github.com/sdwolfe32/trumail/verifier"
	"github.com/sirupsen/logrus"
)

// Service contains all dependencies for the Trumail API
type Service struct {
	log      *logrus.Entry
	timeout  time.Duration
	verifier *verifier.Verifier
}

// NewService generates a new, fully populated Trumail reference
func NewService(l *logrus.Logger, sourceAddr string, timeout int) *Service {
	// Create a new verifier that will be used in the service
	v := verifier.NewVerifier(retrievePTR(), sourceAddr)

	// Restart Dyno if officially confirmed blacklisted
	if err := v.Blacklisted(); err != nil {
		l.WithError(err).Warn("Confirmed Blacklisted! - Restarting Dyno")
		go l.Info(heroku.RestartDyno())
	}

	// Return the fully populated API Service
	return &Service{
		log:      l.WithField("service", "api"),
		timeout:  time.Duration(timeout) * time.Second,
		verifier: v,
	}
}

// retrievePTR attempts to retrieve the PTR record for the IP
// address retrieved via an API call on api.ipify.org
func retrievePTR() string {
	// Request the IP from ipify
	resp, err := http.Get("https://api.ipify.org/")
	if err != nil {
		log.Fatal("Failed to retrieve IP from api.ipify.org")
	}
	defer resp.Body.Close()

	// Decodes the IP response body
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read IP response body")
	}

	// Retrieve the PTR record for our IP and return without a trailing dot
	names, err := net.LookupAddr(string(data))
	if err != nil {
		return string(data)
	}
	return strings.TrimSuffix(names[0], ".")
}
