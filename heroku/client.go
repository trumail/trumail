package heroku

import (
	"errors"
	"fmt"
	"net/http"
)

// baseURL is the base URL of the Heroku API
const baseURL = "https://api.heroku.com"

// Client contains credentials needed to communicate with the
// Heroku API
type Client struct{ appID, dyno, token string }

// NewClient takes a Heroku App ID and token and returns a newly
// generated Heroku Client
func NewClient(appID, dyno, token string) *Client {
	return &Client{appID: appID, dyno: dyno, token: token}
}

// RestartDyno takes a Heroku app ID and an auth token in order to
// restar a Heroku Dyno
func (c *Client) RestartDyno() error {
	if c.appID == "" || c.dyno == "" || c.token == "" {
		return errors.New("Credentials missing to restart heroku dyno")
	}

	// Execute the request on the built path
	return c.do(http.MethodDelete, fmt.Sprintf("/apps/%s/dynos/%s", c.appID, c.dyno))
}

// do creates and executes a request using the passed method
// and path
func (c *Client) do(method, path string) error {
	// Create the request using the passed method and path
	req, err := http.NewRequest(method, baseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/vnd.heroku+json; version=3")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	// Execute the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check the status code before returning with no errors
	if res.StatusCode != http.StatusAccepted {
		return errors.New("Non 202 Status-code received")
	}
	return nil
}
