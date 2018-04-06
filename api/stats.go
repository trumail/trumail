package api

import (
	"net/http"

	"github.com/labstack/echo"
	tinystat "github.com/sdwolfe32/tinystat/client"
	"golang.org/x/sync/errgroup"
)

// ErrStatsFailure is thrown when we fail to retrieve stats from Tinystat
var ErrStatsFailure = echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve statistics")

// Stats is a general summary of Trumail lookups
type Stats struct {
	Deliverable       int64 `json:"deliverable"`
	DeliverableRate   int64 `json:"deliverableRate"`
	Undeliverable     int64 `json:"undeliverable"`
	UndeliverableRate int64 `json:"undeliverableRate"`
	Error             int64 `json:"error"`
	ErrorRate         int64 `json:"errorRate"`
	SuccessRate       int64 `json:"successRate"`
	Total             int64 `json:"total"`
}

// Stats retrieves and returns general Trumail statistics
func (t *TrumailAPI) Stats(c echo.Context) error {
	l := t.log.WithField("handler", "stats")
	l.Debug("New Stats request received")

	// Retrieve all stats from Tinystat
	var g errgroup.Group
	var delCount, undelCount, errCount int64
	g.Go(func() (err error) {
		delCount, err = tinystat.ActionCount("deliverable", "730h")
		return
	})
	g.Go(func() (err error) {
		undelCount, err = tinystat.ActionCount("undeliverable", "730h")
		return
	})
	g.Go(func() (err error) {
		errCount, err = tinystat.ActionCount("error", "730h")
		return
	})
	if err := g.Wait(); err != nil {
		l.WithError(err).Error("Failed to retrieve Tinystat statistics")
		return ErrStatsFailure
	}

	// calculate the grand total and return a stats JSON response
	total := delCount + undelCount + errCount

	l.Debug("Returning Email Lookup")
	return c.JSON(http.StatusOK, &Stats{
		Deliverable:       delCount,
		DeliverableRate:   calcPercent(delCount, total),
		Undeliverable:     undelCount,
		UndeliverableRate: calcPercent(undelCount, total),
		Error:             errCount,
		ErrorRate:         calcPercent(errCount, total),
		SuccessRate:       100 - calcPercent(errCount, total),
		Total:             total,
	})
}

// calcPercent calculates a percentage given two int64
// values and returns the result as an int64
func calcPercent(sub, total int64) int64 {
	return int64((float64(sub) / float64(total)) * float64(100))
}
