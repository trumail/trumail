package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/sdwolfe32/trumail/heroku"
	"github.com/sdwolfe32/trumail/spamhaus"
	"github.com/sdwolfe32/trumail/verifier"
	"github.com/sirupsen/logrus"
)

// maxWorkerCount specifies a maximum number of goroutines allowed
// when processing bulk email lists (not a public endpoint yet)
const maxWorkerCount = 20

// TrumailAPI contains all dependencies for the Trumail API
type TrumailAPI struct {
	log      *logrus.Entry
	hostname string
	verify   *verifier.Verifier
}

// NewTrumailAPI generates a new Trumail reference
func NewTrumailAPI(log *logrus.Logger, hostname, sourceAddr string, timeoutSecs int) *TrumailAPI {
	return &TrumailAPI{
		log:      log.WithField("service", "lookup"),
		hostname: hostname,
		verify: verifier.NewVerifier(&http.Client{Timeout: time.Duration(timeoutSecs) * time.Second},
			maxWorkerCount, hostname, sourceAddr),
	}
}

// RestartIfBlacklisted checks various providers for an SMTP
// block/blacklist on our IP and Triggers a Dyno restart, thus
// retrieving a new IP address if we are
func (t *TrumailAPI) RestartIfBlacklisted(errDetails string) error {
	l := t.log.WithField("method", "spamhaus_restart")

	// Perform Spamhaus blacklist check
	blocked, err := spamhaus.Blocked()
	if err != nil {
		l.WithError(err).Error("Failed to check Spamhaus blacklist status")
		return err
	}

	// Perform Proofpoint blacklist check (paid-service)
	if blocked == false {
		if strings.Contains(errDetails, "proofpoint") {
			blocked = true
		}
	}

	// Restart Dyno if blocked
	if blocked {
		return heroku.RestartDyno()
	}
	return nil
}
