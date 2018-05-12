package verifier

import (
	"encoding/xml"
)

// Verifier contains all dependencies needed to perform educated email
// verification lookups
type Verifier struct{ hostname, sourceAddr string }

// NewVerifier generates a new Verifier using the passed hostname and
// source email address
func NewVerifier(hostname, sourceAddr string) *Verifier {
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
