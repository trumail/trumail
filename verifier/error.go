package verifier

import (
	"errors"
	"strconv"
	"strings"
)

var (
	// Standard Errors
	ErrNoStatusCode      = errors.New("No status code")
	ErrInvalidStatusCode = errors.New("Invalid status code")
	ErrTimeout           = errors.New("The connection to the mail server has timed out")
	ErrNoSuchHost        = errors.New("Mail server does not exist")
	ErrServerUnavailable = errors.New("Mail server is unavailable")
	ErrBlocked           = errors.New("Blocked by mail server")

	// RCPT Errors
	ErrTryAgainLater      = errors.New("Try again later")
	ErrMailboxUnavailable = errors.New("Mailbox Unavailable")
	ErrFullInbox          = errors.New("Recipient out of disk space")
	ErrTooManyRCPT        = errors.New("Too many recipients")
	ErrNoRelay            = errors.New("Not an open relay")
	ErrMailboxBusy        = errors.New("Mailbox busy")
	ErrNeedMAILBeforeRCPT = errors.New("Need MAIL before RCPT")
	ErrRCPTHasMoved       = errors.New("Recipient has moved")
)

// parseSTDErr parses a standard error in order to return a more user
// friendly version of the error
func parseSTDErr(err error) error {
	if err == nil {
		return nil
	}
	errStr := err.Error()

	// Return a friendly error that
	switch {
	case strings.Contains(errStr, "timeout"):
		return ErrTimeout
	case strings.Contains(errStr, "no such host"):
		return ErrNoSuchHost
	case strings.Contains(errStr, "unavailable"):
		return ErrServerUnavailable
	case strings.Contains(errStr, "block") || strings.Contains(errStr, "spamhaus"):
		return ErrBlocked
	default:
		return err
	}
}

// parseRCPTErr receives an MX Servers RCPT response message and generates the
// cooresponding MX error
func parseRCPTErr(err error) error {
	if err == nil {
		return nil
	}
	errStr := err.Error()

	// Verify the length of the error before reading nil indexes
	if len(errStr) < 3 {
		return ErrNoStatusCode
	}

	// Strips out the status code string and converts to an integer for parsing
	status, err := strconv.Atoi(string([]rune(errStr)[0:3]))
	if err != nil {
		return ErrInvalidStatusCode
	}

	// If the status code is above 400 there was an error and we should return it
	if status > 400 {
		switch status {
		case 421:
			return ErrTryAgainLater
		case 450:
			return ErrMailboxBusy
		case 452:
			if strings.Contains(errStr, "full") || strings.Contains(errStr, "space") {
				return ErrFullInbox
			}
			return ErrTooManyRCPT
		case 503:
			return ErrNeedMAILBeforeRCPT
		case 550:
			if strings.Contains(errStr, "spamhaus") {
				return ErrBlocked
			}
			return ErrMailboxUnavailable
		case 551:
			return ErrRCPTHasMoved
		case 552:
			return ErrFullInbox
		case 553:
			return ErrNoRelay
		default:
			return parseSTDErr(err)
		}
	}
	return nil
}

// errStr returns the string representation of an error, returning
// an empty string if the error passed is nil
func errStr(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
