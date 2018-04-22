package heroku

import (
	"errors"
	"fmt"
	"time"

	"github.com/sdwolfe32/httpclient"
)

// baseURL is the base URL of the Heroku API
const baseURL = "https://api.heroku.com"

// Client contains credentials needed to communicate with the
// Heroku API
type Client struct{ client *httpclient.Client }

// NewClient takes a Heroku App ID and token and returns a newly
// generated Heroku Client
func NewClient(token string) *Client {
	return &Client{
		client: httpclient.New(time.Second*10,
			map[string]string{
				"Content-Type":  "application/json",
				"Accept":        "application/vnd.heroku+json; version=3",
				"Authorization": "Bearer " + token,
			},
		),
	}
}

// RestartDyno takes a Heroku app ID and an auth token in order to
// restart a Heroku Dyno
func (c *Client) RestartDyno(appID, dyno string) error {
	if appID == "" || dyno == "" {
		return errors.New("Credentials missing to restart heroku dyno")
	}

	// Execute the request on the built path
	return c.client.Delete(fmt.Sprintf("%s/apps/%s/dynos/%s", baseURL, appID, dyno))
}
