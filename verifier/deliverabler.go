package verifier

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/smtp"
	"strings"
	"time"

	"golang.org/x/net/idna"
)

// Deliverabler contains the context and smtp.Client needed to check email
// address deliverability
type Deliverabler struct {
	client                   *smtp.Client
	timeout                  time.Duration
	domain, host, sourceAddr string
}

// NewDeliverabler generates a new Deliverabler reference
func NewDeliverabler(domain, host, sourceAddr string, timeout time.Duration) (*Deliverabler, error) {
	// Convert any internationalized domain names to ascii
	asciiDomain, err := idna.ToASCII(domain)
	if err != nil {
		asciiDomain = domain
	}

	// Lookup all MX records
	records, err := net.LookupMX(asciiDomain)
	if err != nil {
		return nil, err
	}

	// Verify that at least 1 MX record is found
	if len(records) == 0 {
		return nil, errors.New("No MX records found")
	}

	// Dial the SMTP server with the provided timeout
	client, err := dialSMTPTimeout(records[0].Host+":25", timeout)
	if err != nil {
		return nil, err
	}

	// Sets the HELO/EHLO hostname
	if err := client.Hello(host); err != nil {
		return nil, err
	}

	// Sets a source address
	if err := client.Mail(sourceAddr); err != nil {
		return nil, err
	}
	return &Deliverabler{
		client:     client,
		timeout:    timeout,
		domain:     domain,
		host:       host,
		sourceAddr: sourceAddr,
	}, nil
}

// dialSMTPTimeout is a timeout wrapper for smtp.Dial. It attempts to dial an
// SMTP server and fails with a timeout if the passed timeout is reached while
// attempting to establish a new connection
func dialSMTPTimeout(addr string, timeout time.Duration) (*smtp.Client, error) {
	ch := make(chan *smtp.Client, 1) // Channel holding the new smtp.Client

	// Dial the new smtp connection
	go func() {
		if client, err := smtp.Dial(addr); err == nil {
			ch <- client
		}
	}()

	// Retrieve the smtp client from our client channel or timeout
	select {
	case c := <-ch:
		return c, nil
	case <-time.After(timeout):
		return nil, errors.New("Timeout connecting to mail-exchanger")
	}
}

// IsDeliverable takes an email address and performs the operation of adding
// the email to the envelope. It also receives a number of retries to reconnect
// to the MX server before erring out. If a 250 is received the email is valid
func (d *Deliverabler) IsDeliverable(email string, retry int) error {
	if err := d.client.Rcpt(email); err != nil {
		// In the case of a timeout on the MX connection we need to re-establish and
		// retry the deliverability check
		if shouldReconnect(err) && retry > 0 {
			d.Close()
			time.Sleep(time.Second)                                               // Sleep for 1s as a backoff
			d2, err := NewDeliverabler(d.domain, d.host, d.sourceAddr, d.timeout) // Generate a new client
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
func (d *Deliverabler) HasCatchAll(domain string, retry int) bool {
	return d.IsDeliverable(randomEmail(domain), retry) == nil
}

// Close closes the Deliverablers smtp client connection
func (d *Deliverabler) Close() {
	d.client.Quit()
	d.client.Close()
}

// shouldReconnect determines whether or not we should retry connecting to the
// smtp server based on the response received
func shouldReconnect(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	if strings.Contains(errStr, "i/o timeout") ||
		strings.Contains(errStr, "broken pipe") ||
		strings.Contains(errStr, "use of closed network connection") ||
		strings.Contains(errStr, "connection reset by peer") ||
		strings.Contains(errStr, "multiple regions") ||
		strings.Contains(errStr, "server busy") ||
		strings.Contains(errStr, "EOF") ||
		err == ErrTooManyRCPT || err == ErrTryAgainLater {
		return true
	}
	return false
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
