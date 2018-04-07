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
	case strings.Contains(errStr, "block") ||
		strings.Contains(errStr, "blacklist") ||
		strings.Contains(errStr, "spamhaus"):
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
		return nil, nil
	}
	errStr := strings.ToLower(err.Error())

	// Verify the length of the error before reading nil indexes
	if len(errStr) < 3 {
		return ErrNoStatusCode, err
	}

	// Strips out the status code string and converts to an integer for parsing
	status, convErr := strconv.Atoi(string([]rune(errStr)[0:3]))
	if convErr != nil {
		return ErrInvalidStatusCode, err
	}

	// If the status code is above 400 there was an error and we should return it
	if status > 400 {
		// Don't return an error if the error contains anything about the address
		// being undeliverable
		if strings.Contains(errStr, "undeliverable") ||
			strings.Contains(errStr, "recipient invalid") ||
			strings.Contains(errStr, "recipient rejected") {
			return nil, nil
		}

		switch status {
		case 421:
			return ErrTryAgainLater, err
		case 450:
			return ErrMailboxBusy, err
		case 452:
			if strings.Contains(errStr, "full") ||
				strings.Contains(errStr, "space") {
				return ErrFullInbox, err
			}
			return ErrTooManyRCPT, err
		case 503:
			return ErrNeedMAILBeforeRCPT, err
		case 550: // 550 is Mailbox Unavailable - usually undeliverable
			if strings.Contains(errStr, "spamhaus") ||
				strings.Contains(errStr, "banned") ||
				strings.Contains(errStr, "blocked") ||
				strings.Contains(errStr, "denied") {
				return ErrBlocked, err
			}
			return nil, nil
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
