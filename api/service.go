package api

import (
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/sdwolfe32/trumail/verifier"
	"github.com/sirupsen/logrus"
)

// Service contains all dependencies for the Trumail API
type Service struct {
	log         *logrus.Entry
	timeout     time.Duration
	lookupCache *cache.Cache
	verifier    *verifier.Verifier
}

// NewService generates a new, fully populated Trumail reference
func NewService(log *logrus.Logger, timeout time.Duration,
	verifier *verifier.Verifier) *Service {
	return &Service{
		log:         log.WithField("service", "api"),
		timeout:     timeout,
		lookupCache: cache.New(12*time.Hour, time.Hour),
		verifier:    verifier,
	}
}
