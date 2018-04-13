package heroku

import (
	"os"
)

// DefaultClient is the default client that will be used for
// all Heroku requests
var DefaultClient = NewClient(os.Getenv("HEROKU_TOKEN"))

// RestartApp retarts the Dyno defined in the DefaultClient
func RestartApp() error {
	return DefaultClient.RestartApp(os.Getenv("HEROKU_APP_ID"))
}

// RestartDyno retarts the Dyno defined in the DefaultClient
func RestartDyno() error {
	return DefaultClient.RestartDyno(os.Getenv("HEROKU_APP_ID"), os.Getenv("DYNO"))
}
