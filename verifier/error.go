package verifier

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ErrUnexpectedResponse = "Unexpected response from deliverabler"
	ErrEmailParseFailure  = "Failed to parse email address"

	// Standard Errors
	ErrNoStatusCode      = "No status code"
	ErrInvalidStatusCode = "Invalid status code"
	ErrTimeout           = "The connection to the mail server has timed out"
	ErrNoSuchHost        = "Mail server does not exist"
	ErrServerUnavailable = "Mail server is unavailable"
	ErrBlocked           = "Blocked by mail server"

	// RCPT Errors
	ErrTryAgainLater           = "Try again later"
	ErrFullInbox               = "Recipient out of disk space"
	ErrTooManyRCPT             = "Too many recipients"
	ErrNoRelay                 = "Not an open relay"
	ErrMailboxBusy             = "Mailbox busy"
	ErrExceededMessagingLimits = "Messaging limits have been exceeded"
	ErrNotAllowed              = "Not Allowed"
	ErrNeedMAILBeforeRCPT      = "Need MAIL before RCPT"
	ErrRCPTHasMoved            = "Recipient has moved"
)

// LookupError is an error
type LookupError struct {
	Message string `json:"message" json:"message"`
	Details string `json:"details" json:"details"`
	Report  bool   `json:"-" xml:"-"`
}

// newLookupError creates a new LookupError reference and
// returns it
func newLookupError(message, details string, report bool) *LookupError {
	return &LookupError{message, details, report}
}

// Error satisfies the error interface
func (e *LookupError) Error() string {
	return fmt.Sprintf("%s : %s", e.Message, e.Details)
}

// parseSTDErr parses a standard error in order to return a
// more user friendly version of the error
func parseSTDErr(err error) *LookupError {
	if err == nil {
		return nil
	}
	errStr := err.Error()

	// Return a more understandable error
	switch {
	case insContains(errStr,
		"spamhaus",
		"proofpoint",
		"cloudmark",
		"banned",
		"blocked",
		"denied"):
		return newLookupError(ErrBlocked, errStr, true)
	case insContains(errStr, "timeout"):
		return newLookupError(ErrTimeout, errStr, false)
	case insContains(errStr, "no such host"):
		return newLookupError(ErrNoSuchHost, errStr, false)
	case insContains(errStr, "unavailable"):
		return newLookupError(ErrServerUnavailable, errStr, false)
	default:
		return newLookupError(errStr, errStr, true)
	}
}

// parseRCPTErr receives an MX Servers RCPT response message
// and generates the cooresponding MX error
func parseRCPTErr(err error) *LookupError {
	if err == nil {
		return nil
	}
	errStr := err.Error()

	// Verify the length of the error before reading nil indexes
	if len(errStr) < 3 {
		return newLookupError(ErrNoStatusCode, errStr, true)
	}

	// Strips out the status code string and converts to an integer for parsing
	status, convErr := strconv.Atoi(string([]rune(errStr)[0:3]))
	if convErr != nil {
		return newLookupError(ErrInvalidStatusCode, errStr, true)
	}

	// If the status code is above 400 there was an error and we should return it
	if status > 400 {
		// Don't return an error if the error contains anything about the address
		// being undeliverable
		if insContains(errStr,
			"undeliverable",
			"does not exist",
			"may not exist",
			"invalid address",
			"recipient invalid",
			"recipient rejected",
			"no mailbox") {
			return nil
		}

		switch status {
		case 421:
			return newLookupError(ErrTryAgainLater, errStr, true)
		case 450:
			return newLookupError(ErrMailboxBusy, errStr, true)
		case 451:
			return newLookupError(ErrExceededMessagingLimits, errStr, true)
		case 452:
			if insContains(errStr,
				"full",
				"space",
				"over quota") {
				return newLookupError(ErrFullInbox, errStr, true)
			}
			return newLookupError(ErrTooManyRCPT, errStr, true)
		case 503:
			return newLookupError(ErrNeedMAILBeforeRCPT, errStr, true)
		case 550: // 550 is Mailbox Unavailable - usually undeliverable
			if insContains(errStr,
				"spamhaus",
				"proofpoint",
				"cloudmark",
				"banned",
				"blocked",
				"denied") {
				return newLookupError(ErrBlocked, errStr, true)
			}
			return nil
		case 551:
			return newLookupError(ErrRCPTHasMoved, errStr, true)
		case 552:
			return newLookupError(ErrFullInbox, errStr, false)
		case 553:
			return newLookupError(ErrNoRelay, errStr, true)
		case 554:
			return newLookupError(ErrNotAllowed, errStr, true)
		default:
			return parseSTDErr(err)
		}
	}
	return nil
}

// insContains returns true if any of the substrings
// are found in the passed string. This method of checking
// contains is case insensitive
func insContains(str string, substr ...string) bool {
	lowStr := strings.ToLower(str)
	for _, sub := range substr {
		lowSub := strings.ToLower(sub)
		if strings.Contains(lowStr, lowSub) {
			return true
		}
	}
	return false
}
