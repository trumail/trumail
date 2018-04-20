package verifier

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/sdwolfe32/trumail/httpclient"
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
	v := Verifier{client: httpclient.New(time.Second*5, nil)}
	address := &Address{
		Username: "username",
		Domain:   "domain.com",
		Address:  "username@domain.com",
	}

	configureRequestMock(address.MD5Hash, http.StatusOK)
	defer gock.Off()

	c.Assert(v.HasGravatar(address), check.Equals, true)
}

func (s *gravatarSuite) TestHasGravatarRequestError(c *check.C) {
	v := Verifier{client: httpclient.New(time.Second*5, nil)}
	address := &Address{
		Username: "username",
		Domain:   "domain.com",
		Address:  "username@domain.com",
	}

	gockResponse := configureRequestMock(address.MD5Hash, 200)
	gockResponse.SetError(errors.New("Some error while requesting"))
	defer gock.Off()

	c.Assert(v.HasGravatar(address), check.Equals, false)
}

func (s *gravatarSuite) TestHasGravatarStatusNotOk(c *check.C) {
	v := Verifier{client: httpclient.New(time.Second*5, nil)}
	address := &Address{
		Username: "username",
		Domain:   "domain.com",
		Address:  "username@domain.com",
	}

	configureRequestMock(address.MD5Hash, http.StatusBadRequest)
	defer gock.Off()

	c.Assert(v.HasGravatar(address), check.Equals, false)
}
