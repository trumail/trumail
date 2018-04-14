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

// init seeds the random number generator on
func init() { rand.Seed(time.Now().UTC().UnixNano()) }

// Deliverabler contains the context and smtp.Client needed to check email
// address deliverability
type Deliverabler struct {
	client                       *smtp.Client
	domain, hostname, sourceAddr string
}

// NewDeliverabler generates a new Deliverabler reference, failing if it's
// unable to produce an output within the specified timeout on the client
func (v *Verifier) NewDeliverabler(domain string) (*Deliverabler, error) {
	ch := make(chan interface{}, 1)

	// Create a goroutine that will attempt to connect to the SMTP server
	go func() {
		d, err := newDeliverabler(domain, v.hostname, v.sourceAddr)
		if err != nil {
			ch <- err
		} else {
			ch <- d
		}
	}()

	// Block until a response is produced or timeout
	select {
	case r := <-ch:
		// Return the successful response
		if del, ok := r.(*Deliverabler); ok {
			return del, nil
		}
		// Return the error
		if err, ok := r.(error); ok {
			return nil, err
		}
		return nil, newLookupError(ErrUnexpectedResponse, ErrUnexpectedResponse, false)
	case <-time.After(v.client.Client.Timeout):
		return nil, newLookupError(ErrTimeout, ErrTimeout, false)
	}
}

// newDeliverabler generates a new Deliverabler reference
func newDeliverabler(domain, hostname, sourceAddr string) (*Deliverabler, error) {
	// Convert any internationalized domain names to ascii
	asciiDomain, err := idna.ToASCII(domain)
	if err != nil {
		asciiDomain = domain
	}

	// Dial the SMTP server
	client, err := dialSMTP(asciiDomain)
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

// IsDeliverable takes an email address and performs the operation of adding
// the email to the envelope. It also receives a number of retries to reconnect
// to the MX server before erring out. If a 250 is received the email is valid
func (d *Deliverabler) IsDeliverable(email string, retry int) error {
	if err := d.client.Rcpt(email); err != nil {
		// If we determine a retry should take place
		if shouldRetry(err) && retry > 0 {
			d.Close()                                                    // Close the previous Deliverabler
			d, err = newDeliverabler(d.domain, d.hostname, d.sourceAddr) // Generate a new Deliverabler
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

// dialSMTP receives a domain and attempts to dial the mail server having
// retrieved one or more MX records
func dialSMTP(domain string) (*smtp.Client, error) {
	// Retrieve all MX records
	records, err := net.LookupMX(domain)
	if err != nil {
		return nil, err
	}

	// Verify that at least 1 MX record is found
	if len(records) == 0 {
		return nil, errors.New("No MX records found")
	}

	// Create a channel for receiving the first successful
	// connection on
	client := make(chan *smtp.Client, 1)

	// Attempt to connect to all SMTP servers concurrently
	for _, record := range records {
		addr := record.Host + ":25"
		go func() {
			// Dial the server with a timeout
			conn, err := net.DialTimeout("tcp", addr, time.Minute)
			if err != nil {
				return
			}

			// Generate an smtp client form the connection
			host, _, _ := net.SplitHostPort(addr)
			sc, err := smtp.NewClient(conn, host)
			if err != nil {
				conn.Close()
				return
			}

			// Place the connection on the channel or close it
			select {
			case client <- sc:
			default:
				sc.Close()
			}
		}()
	}
	return <-client, nil
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
		"use of closed network connection",
		"connection reset by peer",
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
