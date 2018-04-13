package verifier

import (
	"crypto/md5"
	"encoding/hex"
	"net/mail"
	"strings"
)

// Address stores all information about an email Address
type Address struct {
	Address  string `json:"address" xml:"address"`
	Username string `json:"username" xml:"username"`
	Domain   string `json:"domain" xml:"domain"`
}

// ParseAddress attempts to parse an email address and return it in the form
// of an Address struct pointer
func ParseAddress(email string) (*Address, error) {
	// Parses the address with the internal go mail address parser
	a, err := mail.ParseAddress(strings.ToLower(email))
	if err != nil {
		return nil, err
	}

	// Find the last occurrence of an @ sign
	index := strings.LastIndex(a.Address, "@")

	// Returns the address with the username and domain split out
	return &Address{
		Username: a.Address[:index],
		Domain:   a.Address[index+1:],
		Address:  a.Address,
	}, nil
}

// MD5 takes a calculation of the md5 Hash and returns a string representation
func (a *Address) MD5() string {
	hashBytes := md5.Sum([]byte(a.Address))
	return hex.EncodeToString(hashBytes[:])
}
