package verifier

import (
	"encoding/xml"
	"time"

	"github.com/sdwolfe32/trumail/httpclient"
)

// Verifier contains all dependencies needed to perform educated email
// verification lookups
type Verifier struct {
	client               *httpclient.Client
	hostname, sourceAddr string
	disp                 *Disposabler
}

// NewVerifier generates a new AddressVerifier reference
func NewVerifier(timeout time.Duration, hostname, sourceAddr string) *Verifier {
	client := httpclient.New(timeout, nil)
	return &Verifier{client, hostname, sourceAddr, NewDisposabler(client)}
}

// Lookup contains all output data for an email verification Lookup
type Lookup struct {
	XMLName xml.Name `json:"-" xml:"lookup"`
	Address
	Deliverable bool `json:"deliverable" xml:"deliverable"`
	FullInbox   bool `json:"fullInbox" xml:"fullInbox"`
	CatchAll    bool `json:"catchAll" xml:"catchAll"`
	Disposable  bool `json:"disposable" xml:"disposable"`
	Gravatar    bool `json:"gravatar" xml:"gravatar"`
}

// Verify parses the passed email and verifies it's deliverability,
// returning any errors that are encountered
func (v *Verifier) Verify(email string) (*Lookup, error) {
	// First parse the email string passed
	a, err := ParseAddress(email)
	if err != nil {
		return nil, newLookupError("", ErrEmailParseFailure)
	}

	// Attempt to form an SMTP Connection
	del, err := v.NewDeliverabler(a.Domain)
	if err != nil {
		return nil, parseSTDErr(err)
	}
	defer del.Close() // Defer close the SMTP connection

	// Declare the lookup to be populated and returned
	var l Lookup
	l.Address = *a

	// Retrieve the catchall status
	if del.HasCatchAll(3) {
		l.CatchAll = true
		l.Deliverable = true
	}
	l.Disposable = v.disp.IsDisposable(a.Domain)

	// Perform the main address verification if not a catchall server
	if !l.CatchAll {
		if err := del.IsDeliverable(a.Address, 3); err != nil {
			le := parseRCPTErr(err)
			if le != nil {
				if le.Err == ErrFullInbox {
					l.FullInbox = true // Set FullInbox and move on
				} else {
					return nil, le // Return if it's a legit error
				}
			}
		} else {
			l.Deliverable = true
		}
	}

	// Check if the email has a Gravatar associated with it
	l.Gravatar = v.HasGravatar(a)
	return &l, nil
}
