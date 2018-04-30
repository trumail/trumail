package verifier

import (
	"encoding/xml"
	"math/rand"
	"time"
)

// Verifier contains all dependencies needed to perform educated email
// verification lookups
type Verifier struct{ hostname, sourceAddr string }

// NewVerifier generates a new Verifier using the passed hostname and
// source email address
func NewVerifier(hostname, sourceAddr string) *Verifier {
	// Seed the random number generator and return a new Verifier
	rand.Seed(time.Now().UTC().UnixNano())
	return &Verifier{hostname, sourceAddr}
}

// Lookup contains all output data for an email verification Lookup
type Lookup struct {
	XMLName xml.Name `json:"-" xml:"lookup"`
	Address
	ValidFormat bool `json:"validFormat" xml:"validFormat"`
	Deliverable bool `json:"deliverable" xml:"deliverable"`
	FullInbox   bool `json:"fullInbox" xml:"fullInbox"`
	HostExists  bool `json:"hostExists" xml:"hostExists"`
	CatchAll    bool `json:"catchAll" xml:"catchAll"`
}

// VerifyTimeout performs an email verification, failing with an ErrTimeout
// if a valid Lookup isn't produced within the timeout passed
func (v *Verifier) VerifyTimeout(email string, timeout time.Duration) (*Lookup, error) {
	ch := make(chan interface{}, 1)

	// Create a goroutine that will attempt to connect to the SMTP server
	go func() {
		d, err := v.Verify(email)
		if err != nil {
			ch <- err
		} else {
			ch <- d
		}
	}()

	// Block until a response is produced or timeout
	select {
	case res := <-ch:
		switch r := res.(type) {
		case *Lookup:
			return r, nil
		case error:
			return nil, r
		default:
			return nil, newLookupError(ErrUnexpectedResponse, ErrUnexpectedResponse)
		}
	case <-time.After(timeout):
		return nil, newLookupError(ErrTimeout, ErrTimeout)
	}
}

// Verify performs an email verification on the passed email address
func (v *Verifier) Verify(email string) (*Lookup, error) {
	// Initialize the lookup
	var l Lookup
	l.Address.Address = email

	// First parse the email address passed
	address, err := ParseAddress(email)
	if err != nil {
		l.ValidFormat = false
		return &l, nil
	}
	l.ValidFormat = true
	l.Address = *address

	// Attempt to form an SMTP Connection
	del, err := NewDeliverabler(address.Domain, v.hostname, v.sourceAddr)
	if err != nil {
		le := parseSMTPError(err)
		if le != nil {
			if le.Message == ErrNoSuchHost {
				l.HostExists = false
				return &l, nil
			}
			return nil, le
		}
		return nil, parseBasicErr(err)
	}
	l.HostExists = true
	defer del.Close() // Defer close the SMTP connection

	// Retrieve the catchall status or check deliverability
	if del.HasCatchAll(3) {
		l.CatchAll = true
		l.Deliverable = true
	} else {
		if err := del.IsDeliverable(address.Address, 3); err != nil {
			le := parseSMTPError(err)
			if le != nil {
				if le.Message == ErrFullInbox {
					l.FullInbox = true // and FullInbox and move on
					return &l, nil
				}
				return nil, le // Return if it's a legit error
			}
		} else {
			l.Deliverable = true
		}
	}
	return &l, nil
}
