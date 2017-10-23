package verifier

import (
	"errors"
	"net/http"
	"testing"

	"gopkg.in/check.v1"
	"gopkg.in/h2non/gock.v1"
)

type gravatarSuite struct{}

var _ = check.Suite(&gravatarSuite{})

func Test(t *testing.T) { check.TestingT(t) } // Just to make discoverable

func configureRequestMock(addressMD5 string, statusCode int) *gock.Response {
	url := "https://en.gravatar.com/" + addressMD5 + ".json"
	return gock.New(url).
		Head("").
		Reply(statusCode)
}

func (s *gravatarSuite) TestHasGravatarStatusOk(c *check.C) {
	address := &Address{
		Username: "username",
		Domain:   "domain.com",
		Address:  "username@domain.com",
	}

	configureRequestMock(address.MD5(), http.StatusOK)
	defer gock.Off()

	c.Assert(HasGravatar(address), check.Equals, true)
}

func (s *gravatarSuite) TestHasGravatarRequestError(c *check.C) {
	address := &Address{
		Username: "username",
		Domain:   "domain.com",
		Address:  "username@domain.com",
	}

	gockResponse := configureRequestMock(address.MD5(), 200)
	gockResponse.SetError(errors.New("Some error while requesting"))
	defer gock.Off()

	c.Assert(HasGravatar(address), check.Equals, false)
}

func (s *gravatarSuite) TestHasGravatarStatusNotOk(c *check.C) {
	address := &Address{
		Username: "username",
		Domain:   "domain.com",
		Address:  "username@domain.com",
	}

	configureRequestMock(address.MD5(), http.StatusBadRequest)
	defer gock.Off()

	c.Assert(HasGravatar(address), check.Equals, false)
}
