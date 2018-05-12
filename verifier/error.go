package verifier

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ErrUnexpectedResponse = "Unexpected response from deliverabler"

	// Standard Errors
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
	Message string `json:"message" xml:"message"`
	Details string `json:"details" xml:"details"`
}

// NewLookupError creates a new LookupError reference and
// returns it
func NewLookupError(message, details string) *LookupError {
	return &LookupError{message, details}
}

// Error satisfies the error interface
func (e *LookupError) Error() string {
	return fmt.Sprintf("%s : %s", e.Message, e.Details)
}

// parseSMTPError receives an MX Servers response message
// and generates the cooresponding MX error
func parseSMTPError(err error) *LookupError {
	if err == nil {
		return nil
	}
	errStr := err.Error()

	// Verify the length of the error before reading nil indexes
	if len(errStr) < 3 {
		return parseBasicErr(err)
	}

	// Strips out the status code string and converts to an integer for parsing
	status, convErr := strconv.Atoi(string([]rune(errStr)[0:3]))
	if convErr != nil {
		return parseBasicErr(err)
	}

	// If the status code is above 400 there was an error and we should return it
	if status > 400 {
		// Don't return an error if the error contains anything about the address
		// being undeliverable
		if insContains(errStr,
			"undeliverable",
			"does not exist",
			"may not exist",
			"user unknown",
			"user not found",
			"invalid address",
			"recipient invalid",
			"recipient rejected",
			"no mailbox") {
			return nil
		}

		switch status {
		case 421:
			return NewLookupError(ErrTryAgainLater, errStr)
		case 450:
			return NewLookupError(ErrMailboxBusy, errStr)
		case 451:
			return NewLookupError(ErrExceededMessagingLimits, errStr)
		case 452:
			if insContains(errStr,
				"full",
				"space",
				"over quota",
				"insufficient",
			) {
				return NewLookupError(ErrFullInbox, errStr)
			}
			return NewLookupError(ErrTooManyRCPT, errStr)
		case 503:
			return NewLookupError(ErrNeedMAILBeforeRCPT, errStr)
		case 550: // 550 is Mailbox Unavailable - usually undeliverable
			if insContains(errStr,
				"spamhaus",
				"proofpoint",
				"cloudmark",
				"banned",
				"blacklisted",
				"blocked",
				"block list",
				"denied") {
				return NewLookupError(ErrBlocked, errStr)
			}
			return nil
		case 551:
			return NewLookupError(ErrRCPTHasMoved, errStr)
		case 552:
			return NewLookupError(ErrFullInbox, errStr)
		case 553:
			return NewLookupError(ErrNoRelay, errStr)
		case 554:
			return NewLookupError(ErrNotAllowed, errStr)
		default:
			return parseBasicErr(err)
		}
	}
	return nil
}

// parseBasicErr parses a basic MX record response and returns
// a more understandable LookupError
func parseBasicErr(err error) *LookupError {
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
		return NewLookupError(ErrBlocked, errStr)
	case insContains(errStr, "timeout"):
		return NewLookupError(ErrTimeout, errStr)
	case insContains(errStr, "no such host"):
		return NewLookupError(ErrNoSuchHost, errStr)
	case insContains(errStr, "unavailable"):
		return NewLookupError(ErrServerUnavailable, errStr)
	default:
		return NewLookupError(errStr, errStr)
	}
}

// insContains returns true if any of the substrings
// are found in the passed string. This method of checking
// contains is case insensitive
func insContains(str string, subStrs ...string) bool {
	for _, subStr := range subStrs {
		if strings.Contains(strings.ToLower(str),
			strings.ToLower(subStr)) {
			return true
		}
	}
	return false
}
