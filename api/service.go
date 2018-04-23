package api

import (
	"time"

	"github.com/sdwolfe32/trumail/verifier"
	"github.com/sirupsen/logrus"
)

// Service contains all dependencies for the Trumail API
type Service struct {
	Logger   *logrus.Entry
	Encode   Encoder
	Timeout  time.Duration
	Verifier *verifier.Verifier
}

// NewService generates a new, fully populated Trumail reference
func NewService(l *logrus.Logger, timeout int, v *verifier.Verifier) *Service {
	// Return the fully populated API Service
	return &Service{
		Logger:   l.WithField("service", "api"),
		Encode:   DefaultEncoder,
		Timeout:  time.Duration(timeout) * time.Second,
		Verifier: v,
	}
}
