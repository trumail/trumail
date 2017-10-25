package verifier

import (
	"gopkg.in/check.v1"
)

type addressSuite struct{}

var _ = check.Suite(&addressSuite{})

func (s *gravatarSuite) TestParseAddress(c *check.C) {
	email := "email_username@domain.com"
	address, err := ParseAddress(email)

	c.Assert(err, check.IsNil)
	c.Assert(address.Username, check.Equals, "email_username")
	c.Assert(address.Domain, check.Equals, "domain.com")
	c.Assert(address.Address, check.Equals, "email_username@domain.com")
}

func (s *gravatarSuite) TestParseAddressForUpperCaseEmails(c *check.C) {
	email := "EMAIL_USERNAME@DOMAIN.COM"
	address, err := ParseAddress(email)

	c.Assert(err, check.IsNil)
	c.Assert(address.Username, check.Equals, "email_username")
	c.Assert(address.Domain, check.Equals, "domain.com")
	c.Assert(address.Address, check.Equals, "email_username@domain.com")
}

func (s *gravatarSuite) TestParseAddressInvalidEmail(c *check.C) {
	email := "email_username@"
	address, err := ParseAddress(email)

	c.Assert(err, check.Not(check.IsNil))
	c.Assert(address, check.IsNil)
}

func (s *gravatarSuite) TestAddressMD5Method(c *check.C) {
	address := Address{
		Username: "user",
		Domain:   "email.com",
		Address:  "user@email.com",
	}

	md5 := address.MD5()
	c.Assert(md5, check.Equals, "b58c6f14d292556214bd64909bcdb118")
}
