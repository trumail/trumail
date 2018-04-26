package api

import (
	"time"

	"github.com/sdwolfe32/trumail/verifier"
)

// Service contains all dependencies for the Trumail API
type Service struct {
	Encode   Encoder
	Timeout  time.Duration
	Verifier *verifier.Verifier
}

// NewService generates a new, fully populated Trumail reference
func NewService(timeout int, v *verifier.Verifier) *Service {
	return &Service{FormatEncoder, time.Duration(timeout) * time.Second, v}
}
