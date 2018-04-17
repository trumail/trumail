package api

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/technosolutionscl/trumail/verifier"
)

// TrumailAPI contains all dependencies for the Trumail API
type TrumailAPI struct {
	log      *logrus.Entry
	timeout  time.Duration
	verifier *verifier.Verifier
}

// NewTrumailAPI generates a new, fully populated Trumail reference
func NewTrumailAPI(log *logrus.Logger, timeout time.Duration,
	verifier *verifier.Verifier) *TrumailAPI {
	return &TrumailAPI{
		log:      log.WithField("service", "api"),
		timeout:  timeout,
		verifier: verifier,
	}
}
