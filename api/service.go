package api

import (
	"net/http"
	"time"

	"github.com/sdwolfe32/trumail/verifier"
	"github.com/sirupsen/logrus"
)

// maxWorkerCount specifies a maximum number of goroutines allowed
// when processing bulk email lists (not a public endpoint yet)
const maxWorkerCount = 20

// TrumailAPI contains all dependencies for the Trumail API
type TrumailAPI struct {
	log         *logrus.Entry
	herokuAppID string
	herokuToken string
	hostname    string
	verify      *verifier.Verifier
}

// NewTrumailAPI generates a new Trumail reference
func NewTrumailAPI(log *logrus.Logger, herokuAppID, herokuToken, hostname, sourceAddr string, timeoutSecs int) *TrumailAPI {
	return &TrumailAPI{
		log:         log.WithField("service", "lookup"),
		herokuAppID: herokuAppID,
		herokuToken: herokuToken,
		hostname:    hostname,
		verify: verifier.NewVerifier(&http.Client{Timeout: time.Duration(timeoutSecs) * time.Second},
			maxWorkerCount, hostname, sourceAddr),
	}
}
