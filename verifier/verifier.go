package verifier

// Verifier contains all dependencies needed to perform educated email
// verification lookups
type Verifier struct{ hostname, sourceAddr string }

// Lookup contains all output data for an email verification Lookup
type Lookup struct {
	Address
	ValidFormat, Deliverable, FullInbox, HostExists, CatchAll bool
}

// NewVerifier generates a new Verifier using the passed hostname and
// source email address
func NewVerifier(hostname, sourceAddr string) *Verifier {
	return &Verifier{hostname, sourceAddr}
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
		le := ParseSMTPError(err)
		if le != nil {
			if le.Message == ErrNoSuchHost {
				l.HostExists = false
				return &l, nil
			}
			return &l, le
		}
		return &l, ParseBasicErr(err)
	}
	l.HostExists = true
	defer del.Close() // Defer close the SMTP connection

	// Retrieve the catchall status or check deliverability
	if del.HasCatchAll(3) {
		l.CatchAll = true
		l.Deliverable = true
	} else {
		if err := del.IsDeliverable(address.Address, 3); err != nil {
			le := ParseSMTPError(err)
			if le != nil {
				if le.Message == ErrFullInbox {
					l.FullInbox = true // and FullInbox and move on
					return &l, nil
				}
				return &l, le // Return if it's a legit error
			}
		} else {
			l.Deliverable = true
		}
	}
	return &l, nil
}
