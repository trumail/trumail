package api

import (
	"github.com/sdwolfe32/trumail/verifier"
	"github.com/sirupsen/logrus"
)

// TrumailAPI contains all dependencies for the Trumail API
type TrumailAPI struct {
	log    *logrus.Entry
	verify *verifier.Verifier
}

// NewTrumailAPI generates a new, fully populated Trumail reference
func NewTrumailAPI(log *logrus.Logger, v *verifier.Verifier) *TrumailAPI {
	return &TrumailAPI{log: log.WithField("service", "api"), verify: v}
}
