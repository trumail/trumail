package verifier

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse550RCPTError(t *testing.T) {
	err := errors.New("550 This mailbox does not exist")
	le := ParseSMTPError(err)
	assert.Equal(t, (*LookupError)(nil), le)
}

func TestParse550BlockedRCPTError(t *testing.T) {
	err := errors.New("550 spamhaus")
	le := ParseSMTPError(err)
	assert.Equal(t, ErrBlocked, le.Message)
	assert.Equal(t, err.Error(), le.Details)
}
