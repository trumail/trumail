package verifier

import check "gopkg.in/check.v1"

type addressSuite struct{}

var _ = check.Suite(&addressSuite{})

func (s *addressSuite) TestParseAddress(c *check.C) {
	email := "email_username@domain.com"
	address, err := ParseAddress(email)

	c.Assert(err, check.IsNil)
	c.Assert(address.Username, check.Equals, "email_username")
	c.Assert(address.Domain, check.Equals, "domain.com")
	c.Assert(address.Address, check.Equals, "email_username@domain.com")
	c.Assert(address.MD5Hash, check.Equals, "629b2a45027be2158761fecb17eb79d6")
}

func (s *addressSuite) TestParseAddress2(c *check.C) {
	email := "email_username@DoMAIn.CoM"
	address, err := ParseAddress(email)

	c.Assert(err, check.IsNil)
	c.Assert(address.Username, check.Equals, "email_username")
	c.Assert(address.Domain, check.Equals, "domain.com")
	c.Assert(address.Address, check.Equals, "email_username@domain.com")
	c.Assert(address.MD5Hash, check.Equals, "629b2a45027be2158761fecb17eb79d6")
}

func (s *addressSuite) TestParseAddressForUpperCaseEmails(c *check.C) {
	email := "EMAIL_USERNAME@DOMAIN.COM"
	address, err := ParseAddress(email)

	c.Assert(err, check.IsNil)
	c.Assert(address.Username, check.Equals, "EMAIL_USERNAME")
	c.Assert(address.Domain, check.Equals, "domain.com")
	c.Assert(address.Address, check.Equals, "EMAIL_USERNAME@domain.com")
	c.Assert(address.MD5Hash, check.Equals, "94d8a553082c902d086c47bd40ccf3c1")
}

func (s *addressSuite) TestParseAddressInvalidEmail(c *check.C) {
	email := "email_username@"
	address, err := ParseAddress(email)

	c.Assert(err, check.Not(check.IsNil))
	c.Assert(address, check.IsNil)
}
