package verifier

import (
	"errors"
	"testing"
)

func TestParse550RCPTError(t *testing.T) {
	err := errors.New("550 This mailbox does not exist")
	basicErr, detailedErr := parseRCPTErr(err)
	equal(t, basicErr, nil)
	equal(t, detailedErr, nil)
}

func TestParse550BlockedRCPTError(t *testing.T) {
	err := errors.New("550 spamhaus")
	basicErr, detailedErr := parseRCPTErr(err)
	equal(t, basicErr, ErrBlocked)
	equal(t, detailedErr, err)
}

// equal is a assertion convenience function used to verify that
// two values equal each other when validating test results
func equal(t *testing.T, actual interface{}, expected interface{}) {
	if actual != expected {
		t.Logf("%v does not equal %v", actual, expected)
		t.Fail()
	}
}
