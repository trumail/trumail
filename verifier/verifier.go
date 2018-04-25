package verifier

import (
	"encoding/xml"
	"time"

	"github.com/sdwolfe32/httpclient"
)

// Verifier contains all dependencies needed to perform educated email
// verification lookups
type Verifier struct {
	client               *httpclient.Client
	hostname, sourceAddr string
	disp                 *Disposabler
}

// NewVerifier generates a new httpclient.Client using the passed timeout
// and then returns a new Verifier reference that will be used to Verify
// email addresses
func NewVerifier(hostname, sourceAddr string) *Verifier {
	client := httpclient.NewBaseClient(time.Second * 30)
	return &Verifier{client, hostname, sourceAddr, NewDisposabler(client)}
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
	Disposable  bool `json:"disposable" xml:"disposable"`
	Gravatar    bool `json:"gravatar" xml:"gravatar"`
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
	l := new(Lookup)
	l.Address.Address = email

	// First parse the email address passed
	address, err := ParseAddress(email)
	if err != nil {
		l.ValidFormat = false
		return l, nil
	}
	l.ValidFormat = true
	l.Address = *address

	// Set all parse dependent but SMTP independent values
	l.Disposable = v.disp.IsDisposable(address.Domain)
	l.Gravatar = v.HasGravatar(address.MD5Hash)

	// Attempt to form an SMTP Connection
	del, err := NewDeliverabler(address.Domain, v.hostname, v.sourceAddr)
	if err != nil {
		le := parseSMTPError(err)
		if le != nil && le.Message == ErrNoSuchHost {
			l.HostExists = false
			return l, nil
		}
		return nil, err
	}
	l.HostExists = true
	defer del.Close() // Defer close the SMTP connection

	// Retrieve the catchall status
	if del.HasCatchAll(3) {
		l.CatchAll = true
		l.Deliverable = true
	}

	// Perform the main address verification if not a catchall server
	if !l.CatchAll {
		if err := del.IsDeliverable(address.Address, 3); err != nil {
			le := parseSMTPError(err)
			if le != nil {
				if le.Message == ErrFullInbox {
					l.FullInbox = true // Set FullInbox and move on
					return l, nil
				}
				return nil, le // Return if it's a legit error
			}
		} else {
			l.Deliverable = true
		}
	}
	return l, nil
}
