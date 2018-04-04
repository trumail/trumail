package api

import (
	"errors"
	"fmt"
	"net/http"
)

const herokuBaseURL = "https://api.heroku.com"

// restartDyno takes a Heroku app ID and an auth token in order to
// restar a Heroku Dyno
func restartDyno(appID, token string) error {
	if appID == "" || token == "" {
		return errors.New("No credentials found to restart heroku dynos")
	}
	// Create the restart request
	path := fmt.Sprintf("/apps/%s/dynos", appID)
	req, err := http.NewRequest(http.MethodDelete, herokuBaseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/vnd.heroku+json; version=3")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// Execute the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// Check the status code before returning with no errors
	if res.StatusCode != http.StatusOK {
		return errors.New("Non 200 Status-code received")
	}
	return nil
}
