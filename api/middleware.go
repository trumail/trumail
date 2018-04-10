package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo"
)

// DefaultRateLimitMiddleware is a default middleware using
// a 2 second rate-limit
var DefaultRateLimitMiddleware = NewRateLimitMiddleware()

// RateLimit uses the DefaultRateLimitMiddleware to rate
// limit requests
func RateLimit(next echo.HandlerFunc) echo.HandlerFunc {
	return DefaultRateLimitMiddleware.Do(next)
}

// ErrRateLimitExceeded is thrown when an IP exceeds the
// specified rate-limit
var ErrRateLimitExceeded = echo.NewHTTPError(http.StatusTooManyRequests,
	"Rate limit exceeded - If you'd like a higher request volume please contact steven@swolfe.me")

// RateLimitMiddleware is a middleware for limiting request
// speed
type RateLimitMiddleware struct {
	sync.Mutex
	ipMap map[string]time.Time
}

// NewRateLimitMiddleware generates a new RateLimitMiddleware
// reference
func NewRateLimitMiddleware() *RateLimitMiddleware {
	return &RateLimitMiddleware{ipMap: make(map[string]time.Time)}
}

// Do returns an error if the ip passed has performed too
// many requests in the defined period of time.
func (m *RateLimitMiddleware) Do(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract the IP from the request
		ip := c.RealIP()

		// Lock the map
		m.Lock()
		defer m.Unlock()

		// If this IP is in the map and it's last request
		// was within the specified ratelimit timeframe
		if last, ok := m.ipMap[ip]; ok &&
			last.After(time.Now().Add(-1*time.Second)) {
			return ErrRateLimitExceeded
		}

		// Set a new last request time and allow the request
		m.ipMap[ip] = time.Now()
		return next(c)
	}
}
