package verifier

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/sdwolfe32/httpclient"
	check "gopkg.in/check.v1"
	gock "gopkg.in/h2non/gock.v1"
)

type gravatarSuite struct{}

var _ = check.Suite(&gravatarSuite{})

func Test(t *testing.T) { check.TestingT(t) } // Just to make discoverable

func configureRequestMock(addressMD5 string, statusCode int) *gock.Response {
	url := "https://en.gravatar.com/" + addressMD5 + ".json"
	return gock.New(url).Head("").Reply(statusCode)
}

func (s *gravatarSuite) TestHasGravatarStatusOk(c *check.C) {
	v := Verifier{client: httpclient.New(time.Second*5, nil)}
	configureRequestMock("asdf1234", http.StatusOK)
	defer gock.Off()

	c.Assert(v.HasGravatar("asdf1234"), check.Equals, true)
}

func (s *gravatarSuite) TestHasGravatarRequestError(c *check.C) {
	v := Verifier{client: httpclient.New(time.Second*5, nil)}
	gockResponse := configureRequestMock("asdf1234", 200)
	gockResponse.SetError(errors.New("Some error while requesting"))
	defer gock.Off()

	c.Assert(v.HasGravatar("asdf1234"), check.Equals, false)
}

func (s *gravatarSuite) TestHasGravatarStatusNotOk(c *check.C) {
	v := Verifier{client: httpclient.New(time.Second*5, nil)}
	configureRequestMock("asdf1234", http.StatusBadRequest)
	defer gock.Off()

	c.Assert(v.HasGravatar("asdf1234"), check.Equals, false)
}
