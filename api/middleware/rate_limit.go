package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo"
)

var (
	// DefaultRateLimiter is a default rate-limiting middleware
	// that allows up to 1000 requests every 24 hours
	DefaultRateLimiter = NewRateLimiter(100, time.Hour*24)
	// ErrRateLimitExceeded is thrown when an IP exceeds the
	// specified rate-limit
	ErrRateLimitExceeded = echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded - If you'd like a higher request volume please contact steven@swolfe.me")
)

// RateLimit uses the DefaultRateLimiter to rate limit requests
func RateLimit(next echo.HandlerFunc) echo.HandlerFunc {
	return DefaultRateLimiter.Do(next)
}

// RateLimiter is a middleware for limiting request
// speed to a maximum over a set interval
type RateLimiter struct {
	max      int64         // The maximum number of requests allowed in the interval
	interval time.Duration // The duration to assert the max
	ipMap    *sync.Map     // IP-Address -> ReqData
}

// NewRateLimiter generates a new RateLimiter reference
func NewRateLimiter(max int64, interval time.Duration) *RateLimiter {
	return &RateLimiter{max: max, interval: interval, ipMap: &sync.Map{}}
}

// ReqData contains recent request data
type ReqData struct {
	start time.Time
	count int64
}

// NewReqData generates a new ReqData reference with the
// start time
func NewReqData() *ReqData { return &ReqData{start: time.Now()} }

// Count increments the count on a ReqData
func (f *ReqData) Count() { f.count++ }

// Do returns an error if the ip passed has performed too
// many requests in the defined period of time.
func (r *RateLimiter) Do(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract the IP from the request
		ip := c.RealIP()

		// Load an existing or new ReqData interface
		rfIface, ok := r.ipMap.Load(ip)
		if !ok {
			newRF := NewReqData()
			r.ipMap.Store(ip, newRF)
			rfIface = newRF
		}

		// Type assert the ReqData
		rf := rfIface.(*ReqData)

		// If the IPMeta is valid for this time inteval AND
		// If The count for this timeframe exceeds the max
		if rf.start.After(time.Now().Add(-1 * r.interval)) {
			if rf.count >= r.max {
				return ErrRateLimitExceeded
			}
		} else {
			// Otherwise refresh the freq data for this interval
			r.ipMap.Store(ip, NewReqData())
		}

		// Count a new request and return
		rf.Count()
		return next(c)
	}
}
