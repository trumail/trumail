package verifier

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

var (
	ErrTryAgainLater      = errors.New("Try again later")
	ErrNoSuchRCPT         = errors.New("No such recipient")
	ErrFullInbox          = errors.New("Recipient out of disk space")
	ErrTooManyRCPT        = errors.New("Too many recipients")
	ErrNoRelay            = errors.New("Not an open relay")
	ErrMailboxBusy        = errors.New("Mailbox busy")
	ErrNeedMAILBeforeRCPT = errors.New("Need MAIL before RCPT")
	ErrRCPTHasMoved       = errors.New("Recipient has moved")
)

// Deliverabler defines all functionality for checking an email addresses
// deliverability
type Deliverabler interface {
	IsDeliverable(email string, retry int) error
	HasCatchAll(domain string, retry int) bool
	Close()
}

// deliverabler contains the context and smtp.Client needed to check email
// address deliverability
type deliverabler struct {
	client                   *smtp.Client
	domain, host, sourceAddr string
}

// NewDeliverabler generates a new Deliverabler
func NewDeliverabler(domain, host, sourceAddr string) (Deliverabler, error) {
	// Looks up all MX records
	records, err := net.LookupMX(domain)
	if err != nil {
		return nil, err
	}

	// Verify that at least 1 MX record is found
	if len(records) == 0 {
		return nil, errors.New("No MX records found")
	}

	// Dials the tcp connection
	conn, err := net.DialTimeout("tcp", records[0].Host+":25", 3*time.Second)
	if err != nil {
		return nil, err
	}

	// Connect to the SMTP server
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return nil, err
	}

	// Sets the HELO hostname
	if err := client.Hello(host); err != nil {
		return nil, err
	}

	// Sets a source address
	if err := client.Mail(sourceAddr); err != nil {
		return nil, err
	}
	return &deliverabler{
		client:     client,
		domain:     domain,
		host:       host,
		sourceAddr: sourceAddr,
	}, nil
}

// IsDeliverable takes an email address and performs the operation of adding
// the email to the envelope. It also receives a number of retries to reconnect
// to the MX server before erring out. If a 250 is received the email is valid
func (d *deliverabler) IsDeliverable(email string, retry int) error {
	if err := parseRCPTErr(d.client.Rcpt(email)); err != nil {
		// In the case of a timeout on the MX connection we need to re-establish and
		// retry the deliverability check
		if shouldReconnect(err) && retry > 0 {
			d.Close()
			time.Sleep(time.Second)                                    // Sleep for 1s as a backoff
			d2, err := NewDeliverabler(d.domain, d.host, d.sourceAddr) // Generate a new client
			if err != nil {
				return err
			}
			return d2.IsDeliverable(email, retry-1) // Retry deliverability check
		}
		return err
	}
	return nil
}

// HasCatchAll checks the deliverability of a randomly generated address in
// order to verify the existence of a catch-all
func (d *deliverabler) HasCatchAll(domain string, retry int) bool {
	return d.IsDeliverable(randomEmail(domain), retry) == nil
}

// Close closes the Deliverablers smtp client connection
func (d *deliverabler) Close() {
	d.client.Quit()
	d.client.Close()
}

// shouldReconnect determines whether or not we should retry connecting to the
// smtp server based on the response received
func shouldReconnect(err error) bool {
	errStr := err.Error()
	if strings.Contains(errStr, "i/o timeout") ||
		strings.Contains(errStr, "broken pipe") ||
		strings.Contains(errStr, "use of closed network connection") ||
		strings.Contains(errStr, "connection reset by peer") ||
		strings.Contains(errStr, "EOF") ||
		err == ErrTooManyRCPT || err == ErrTryAgainLater {
		return true
	}
	return false
}

// parseRCPTErr receives an MX Servers RCPT response message and generates the
// cooresponding XM error
func parseRCPTErr(err error) error {
	if err == nil {
		return nil
	}
	response := err.Error()

	// Strips out the status code string and converts to an integer for parsing
	status, err := strconv.Atoi(string([]rune(response)[0:3]))
	if err != nil {
		return err
	}
	message := string([]rune(response)[3:])

	// If the status code is above 400 there was an error and we should return it
	if status > 400 {
		switch status {
		case 421:
			return ErrTryAgainLater
		case 450:
			return ErrMailboxBusy
		case 452:
			if strings.Contains(message, "full") || strings.Contains(message, "space") {
				return ErrFullInbox
			}
			return ErrTooManyRCPT
		case 503:
			return ErrNeedMAILBeforeRCPT
		case 550:
			return ErrNoSuchRCPT
		case 551:
			return ErrRCPTHasMoved
		case 552:
			return ErrFullInbox
		case 553:
			return ErrNoRelay
		default:
			return errors.New(response)
		}
	}
	return nil
}

// randomEmail generates a random email address using the domain passed. Used
// primarily for checking the existence of a catch-all address
func randomEmail(domain string) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 20)
	for i := 0; i < 20; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return fmt.Sprintf("%s@%s", string(result), domain)
}
