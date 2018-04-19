package api

import (
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/sdwolfe32/trumail/verifier"
	"github.com/sirupsen/logrus"
)

// TrumailAPI contains all dependencies for the Trumail API
type TrumailAPI struct {
	log         *logrus.Entry
	timeout     time.Duration
	lookupCache *cache.Cache
	verifier    *verifier.Verifier
}

// NewTrumailAPI generates a new, fully populated Trumail reference
func NewTrumailAPI(log *logrus.Logger, timeout time.Duration,
	verifier *verifier.Verifier) *TrumailAPI {
	return &TrumailAPI{
		log:         log.WithField("service", "api"),
		timeout:     timeout,
		lookupCache: cache.New(12*time.Hour, time.Hour),
		verifier:    verifier,
	}
}
