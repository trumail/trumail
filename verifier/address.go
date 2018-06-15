package verifier

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"strings"
)

// ErrNoAtChar is thrown when no '@' character is found on an email
// address
var ErrNoAtChar = errors.New("No '@' character found on email address")

// Address stores all information about an email Address
type Address struct{ Address, Username, Domain, MD5Hash string }

// ParseAddress attempts to parse an email address and return it in the form
// of an Address struct pointer - domain case insensitive
func ParseAddress(email string) (*Address, error) {
	// Parses the address with the internal go mail address parser
	a, err := mail.ParseAddress(unescape(email))
	if err != nil {
		return nil, err
	}

	// Find the last occurrence of an @ sign
	index := strings.LastIndex(a.Address, "@")
	if index == -1 {
		return nil, ErrNoAtChar
	}

	// Parse the username, domain and case unique address
	username := a.Address[:index]
	domain := strings.ToLower(a.Address[index+1:])
	address := fmt.Sprintf("%s@%s", username, domain)

	// Hash the address
	hashBytes := md5.Sum([]byte(address))
	md5Hash := hex.EncodeToString(hashBytes[:])

	// Returns the Address with the username and domain split out
	return &Address{address, username, domain, md5Hash}, nil
}

// unescape attempts to return a query un-escaped version of the
// passed string, returning the original string of an error occurs
func unescape(str string) string {
	esc, err := url.QueryUnescape(str)
	if err != nil {
		return str
	}
	return esc
}
