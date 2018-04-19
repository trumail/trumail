package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo"
)

// AuthToken is the key to be used on the Trumail Token header
const AuthToken = "X-Auth-Token"

var (
	// ErrRateLimitExceeded is thrown when an IP exceeds the
	// specified rate-limit
	ErrRateLimitExceeded = echo.NewHTTPError(http.StatusTooManyRequests,
		"Rate limit exceeded - If you'd like a higher request volume please contact steven@swolfe.me")
)

// RateLimiter is a middleware for limiting request
// speed to a maximum over a set interval
type RateLimiter struct {
	token    string        // A token that can be used to bypass the rate limit
	max      int64         // The maximum number of requests allowed in the interval
	interval time.Duration // The duration to assert the max
	ipMap    *sync.Map     // IP-Address -> ReqData
}

// ReqData contains recent request data
type ReqData struct {
	start time.Time
	count int64
}

// LimitStatus is returned when a request is made for an
// IPs current rate limit standing
type LimitStatus struct {
	Max      int64         `json:"max"`
	Interval time.Duration `json:"interval"`
	Current  int64         `json:"current"`
}

// NewRateLimiter generates a new RateLimiter reference
func NewRateLimiter(token string, max int64, interval time.Duration) *RateLimiter {
	return &RateLimiter{token, max, interval, &sync.Map{}}
}

// NewReqData generates a new ReqData reference with the
// start time
func NewReqData() *ReqData { return &ReqData{start: time.Now()} }

// Count increments the count on a ReqData
func (f *ReqData) Count() { f.count++ }

// RateLimit returns an error if the ip passed has performed too
// many requests in the defined period of time.
func (r *RateLimiter) RateLimit(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// First check the header for the existence of the Trumail Token
		if r.token != "" {
			if c.Request().Header.Get(AuthToken) == r.token {
				return next(c)
			}
		}

		// Allocate the users IP and ReqData
		ip := realIP(c)
		rd := r.reqData(ip)

		// Whether the ReqData is expired
		valid := rd.start.After(time.Now().Add(-1 * r.interval))

		// If the valid count for this timeframe exceeds the max
		if valid && rd.count >= r.max {
			return ErrRateLimitExceeded
		}

		// If the IPMeta is invalid (expired), store a new one
		if !valid {
			r.ipMap.Store(ip, NewReqData())
		}

		// Count a new request and return
		rd.Count()
		return next(c)
	}
}

// LimitStatus retrieves and returns general Trumail statistics
func (r *RateLimiter) LimitStatus(c echo.Context) error {
	// Allocate the users ReqData
	rd := r.reqData(realIP(c))

	// Return the current rate limit standing
	return c.JSON(http.StatusOK, &LimitStatus{
		Max:      r.max,
		Interval: r.interval,
		Current:  rd.count,
	})
}

// realIP checks for a Cloudflare connecting IP on the request
// and returns it if found, otherwise it returns the IP echo
// determines from X-Forwarded-For
func realIP(c echo.Context) string {
	// Allocate the Cloudflare connecting IP
	cfip := c.Request().Header.Get("CF-Connecting-IP")

	// Return it or return the default IP
	if cfip != "" {
		return cfip
	} else {
		return c.RealIP()
	}
}

// reqData returns ReqData found in the syncmap keyed
// by the requestors IP address
func (r *RateLimiter) reqData(ip string) *ReqData {
	// Load an existing or new ReqData interface
	if rdIface, ok := r.ipMap.Load(ip); ok {
		return rdIface.(*ReqData)
	}

	// Create a new ReqData and return it
	newRD := NewReqData()
	r.ipMap.Store(ip, newRD)
	return newRD
}
