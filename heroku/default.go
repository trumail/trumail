package heroku

import (
	"os"
)

// DefaultClient is the default client that will be used for
// all Heroku requests
var DefaultClient = NewClient(os.Getenv("HEROKU_APP_ID"), os.Getenv("HEROKU_TOKEN"))

// RestartDyno retarts the Dyno defined in the DefaultClient
func RestartDyno() error {
	return DefaultClient.RestartDyno()
}
