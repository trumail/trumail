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
func parseSTDErr(err error) (error, error) {
	if err == nil {
		return nil, nil
	}
	errStr := strings.ToLower(err.Error())

	// Return a friendly error that
	switch {
	case strings.Contains(errStr, "block") || strings.Contains(errStr, "blacklist") || strings.Contains(errStr, "spamhaus"):
		return ErrBlocked, err
	case strings.Contains(errStr, "timeout"):
		return ErrTimeout, err
	case strings.Contains(errStr, "no such host"):
		return ErrNoSuchHost, err
	case strings.Contains(errStr, "unavailable"):
		return ErrServerUnavailable, err
	default:
		return err, err
	}
}

// parseRCPTErr receives an MX Servers RCPT response message and generates the
// cooresponding MX error
func parseRCPTErr(err error) (error, error) {
	if err == nil {
		return nil, err
	}
	errStr := strings.ToLower(err.Error())

	// Verify the length of the error before reading nil indexes
	if len(errStr) < 3 {
		return ErrNoStatusCode, err
	}

	// Strips out the status code string and converts to an integer for parsing
	status, err := strconv.Atoi(string([]rune(errStr)[0:3]))
	if err != nil {
		return ErrInvalidStatusCode, err
	}

	// Don't return an error if the error contains anything about the address
	// being undeliverable
	if strings.Contains(errStr, "undeliverable") {
		return nil, nil
	}

	// If the status code is above 400 there was an error and we should return it
	if status > 400 {
		switch status {
		case 421:
			return ErrTryAgainLater, err
		case 450:
			return ErrMailboxBusy, err
		case 452:
			if strings.Contains(errStr, "full") || strings.Contains(errStr, "space") {
				return ErrFullInbox, err
			}
			return ErrTooManyRCPT, err
		case 503:
			return ErrNeedMAILBeforeRCPT, err
		case 550:
			if strings.Contains(errStr, "spamhaus") {
				return ErrBlocked, err
			}
			return ErrMailboxUnavailable, err
		case 551:
			return ErrRCPTHasMoved, err
		case 552:
			return ErrFullInbox, err
		case 553:
			return ErrNoRelay, err
		default:
			return parseSTDErr(err)
		}
	}
	return nil, nil
}

// errStr returns the string representation of an error, returning
// an empty string if the error passed is nil
func errStr(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
