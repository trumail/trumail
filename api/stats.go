package api

import (
	"net/http"

	"github.com/labstack/echo"
	tinystat "github.com/sdwolfe32/tinystat/client"
	"golang.org/x/sync/errgroup"
)

// ErrStatsFailure is thrown when we fail to retrieve stats from Tinystat
var ErrStatsFailure = echo.NewHTTPError(http.StatusInternalServerError,
	"Failed to retrieve statistics")

// StatsSummary is an overall stats summary
type StatsSummary struct {
	Daily   Stats `json:"daily"`
	Monthly Stats `json:"monthly"`
}

// Stats is a general summary of Trumail lookups
type Stats struct {
	Deliverable   int64 `json:"deliverable"`
	Undeliverable int64 `json:"undeliverable"`
	Errors        int64 `json:"errors"`
	SuccessRate   int64 `json:"successRate"`
}

// Stats retrieves and returns general Trumail statistics
func (t *TrumailAPI) Stats(c echo.Context) error {
	l := t.log.WithField("handler", "stats")
	l.Debug("New Stats request received")

	// Retrieve all stats from Tinystat
	var s StatsSummary
	var g errgroup.Group
	g.Go(func() error { return stats("24h", &s.Daily) })
	g.Go(func() error { return stats("730h", &s.Monthly) })
	if err := g.Wait(); err != nil {
		l.WithError(err).Error("Failed to retrieve Tinystat statistics")
		return ErrStatsFailure
	}

	l.Debug("Returning Stats Response")
	return c.JSON(http.StatusOK, s)
}

// stats populates the passed Stats reference with stats
// for the previous passed duration
func stats(duration string, s *Stats) error {
	var g errgroup.Group
	g.Go(func() (err error) {
		s.Deliverable, err = tinystat.ActionCount("deliverable", duration)
		return
	})
	g.Go(func() (err error) {
		s.Undeliverable, err = tinystat.ActionCount("undeliverable", duration)
		return
	})
	g.Go(func() (err error) {
		s.Errors, err = tinystat.ActionCount("error", duration)
		return
	})
	if err := g.Wait(); err != nil {
		return err
	}
	total := s.Deliverable + s.Undeliverable + s.Errors
	s.SuccessRate = 100 - calcPercent(s.Errors, total)
	return nil
}

// calcPercent calculates a percentage given two int64
// values and returns the result as an int64
func calcPercent(sub, total int64) int64 {
	return int64((float64(sub) / float64(total)) * float64(100))
}
