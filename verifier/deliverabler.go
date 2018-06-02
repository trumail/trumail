package verifier

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/smtp"
	"time"

	"golang.org/x/net/idna"
)

// Deliverabler contains the context and smtp.Client needed to check
// email address deliverability
type Deliverabler struct {
	client                       *smtp.Client
	domain, hostname, sourceAddr string
}

// NewDeliverabler generates a new Deliverabler reference
func NewDeliverabler(domain, hostname, sourceAddr string) (*Deliverabler, error) {
	// Dial any SMTP server that will accept a connection
	client, err := mailDialTimeout(domain, time.Minute)
	if err != nil {
		return nil, err
	}

	// Sets the HELO/EHLO hostname
	if err := client.Hello(hostname); err != nil {
		return nil, err
	}

	// Sets a source address
	if err := client.Mail(sourceAddr); err != nil {
		return nil, err
	}

	// Return the deliverabler if successful
	return &Deliverabler{client, domain, hostname, sourceAddr}, nil
}

// dialSMTP receives a domain and attempts to dial the mail server having
// retrieved one or more MX records
func mailDialTimeout(domain string, timeout time.Duration) (*smtp.Client, error) {
	// Convert any internationalized domain names to ascii
	asciiDomain, err := idna.ToASCII(domain)
	if err != nil {
		asciiDomain = domain
	}

	// Retrieve all MX records
	records, err := net.LookupMX(asciiDomain)
	if err != nil {
		return nil, err
	}

	// Verify that at least 1 MX record is found
	if len(records) == 0 {
		return nil, errors.New("No MX records found")
	}

	// Create a channel for receiving responses from
	ch := make(chan interface{}, 1)

	// Done indicates if we're still waiting on dial responses
	var done bool

	// Attempt to connect to all SMTP servers concurrently
	for _, record := range records {
		addr := record.Host + ":25"
		go func() {
			c, err := smtpDialTimeout(addr, timeout)
			if err != nil {
				if !done {
					ch <- err
				}
				return
			}

			// Place the client on the channel or close it
			switch {
			case !done:
				done = true
				ch <- c
			default:
				c.Close()
			}
		}()
	}

	// Collect errors or return a client
	var errSlice []error
	for {
		res := <-ch
		switch r := res.(type) {
		case *smtp.Client:
			return r, nil
		case error:
			errSlice = append(errSlice, r)
			if len(errSlice) == len(records) {
				return nil, errSlice[0]
			}
		default:
			return nil, errors.New("Unexpected response dialing SMTP server")
		}
	}
}

// smtpDialTimeout is a timeout wrapper for smtp.Dial. It attempts to dial an
// SMTP server and fails with a timeout if the passed timeout is reached while
// attempting to establish a new connection
func smtpDialTimeout(addr string, timeout time.Duration) (*smtp.Client, error) {
	// Channel holding the new smtp.Client or error
	ch := make(chan interface{}, 1)

	// Dial the new smtp connection
	go func() {
		client, err := smtp.Dial(addr)
		if err != nil {
			ch <- err
			return
		}
		ch <- client
	}()

	// Retrieve the smtp client from our client channel or timeout
	select {
	case res := <-ch:
		switch r := res.(type) {
		case *smtp.Client:
			return r, nil
		case error:
			return nil, r
		default:
			return nil, errors.New("Unexpected response dialing SMTP server")
		}
	case <-time.After(timeout):
		return nil, errors.New("Timeout connecting to mail-exchanger")
	}
}

// IsDeliverable takes an email address and performs the operation of adding
// the email to the envelope. It also receives a number of retries to reconnect
// to the MX server before erring out. If a 250 is received the email is valid
func (d *Deliverabler) IsDeliverable(email string, retry int) error {
	if err := d.client.Rcpt(email); err != nil {
		// If we determine a retry should take place
		if shouldRetry(err) && retry > 0 {
			d.Close()                                                    // Close the previous Deliverabler
			d, err = NewDeliverabler(d.domain, d.hostname, d.sourceAddr) // Generate a new Deliverabler
			if err != nil {
				return err
			}
			return d.IsDeliverable(email, retry-1) // Retry deliverability check
		}
		return err
	}
	return nil
}

// HasCatchAll checks the deliverability of a randomly generated address in
// order to verify the existence of a catch-all
func (d *Deliverabler) HasCatchAll(retry int) bool {
	return d.IsDeliverable(randomEmail(d.domain), retry) == nil
}

// Close closes the Deliverablers SMTP client connection
func (d *Deliverabler) Close() {
	d.client.Quit()
	d.client.Close()
}

// shouldRetry determines whether or not we should retry connecting to the
// smtp server based on the response received
func shouldRetry(err error) bool {
	if err == nil {
		return false
	}
	return insContains(err.Error(),
		"i/o timeout",
		"broken pipe",
		"exceeded the maximum number of connections",
		"use of closed network connection",
		"connection reset by peer",
		"connection declined",
		"connection refused",
		"multiple regions",
		"server busy",
		"eof")
}

// randomEmail generates a random email address using the domain passed. Used
// primarily for checking the existence of a catch-all address
func randomEmail(domain string) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 20)
	for i := 0; i < 20; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return fmt.Sprintf("%s@%s", string(result), domain)
}
