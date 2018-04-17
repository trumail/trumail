package api

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo"
)

var (
	// ErrRateLimitExceeded is thrown when an IP exceeds the
	// specified rate-limit
	ErrRateLimitExceeded = echo.NewHTTPError(http.StatusTooManyRequests,
		"Rate limit exceeded - If you'd like a higher request volume please contact steven@swolfe.me")
)

// RateLimiter is a middleware for limiting request
// speed to a maximum over a set interval
type RateLimiter struct {
	max        int64         // The maximum number of requests allowed in the interval
	interval   time.Duration // The duration to assert the max
	ipMap      *sync.Map     // IP-Address -> ReqData
	cidrcustom string        // list of custom cidr in config.go
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
func NewRateLimiter(max int64, interval time.Duration, cidrcustom string) *RateLimiter {
	return &RateLimiter{max: max, interval: interval, ipMap: &sync.Map{}, cidrcustom: cidrcustom}
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
		ip := c.RealIP()    // The requestors IP
		rd := r.reqData(ip) // The requestors ReqData
		isExcluded := false
		interval := r.interval
		max := r.max
		//log.Println("REAL IP " + ip)
		// if we have custom subnets set
		if r.cidrcustom != "" {
			//log.Println("CIDR CUSTOM " + r.cidrcustom)
			// split the cidr list string to array
			s := strings.Split(r.cidrcustom, ",")
			for _, i := range s {
				//if not previously excluded
				if !isExcluded {
					//log.Println("Found " + i)
					// split into pieces
					_s := strings.Split(i, "|")
					cidr, cmax, cint := _s[0], _s[1], _s[2]
					//log.Println("cidr " + cidr)
					_, subnet, _ := net.ParseCIDR(cidr)
					_ip := net.ParseIP(ip)
					// if current ip fits in the subnet, apply values
					if subnet.Contains(_ip) {
						if cmax == "0" && cint == "0" {
							isExcluded = true
							log.Println(ip + " is in range " + cidr + " and has no limit")
						} else {
							timd, err := strconv.ParseInt(cint, 10, 64)
							if err != nil {
								fmt.Println(err.Error())
								timd = 10
							}
							interval = time.Hour * time.Duration(timd)
							max, err = strconv.ParseInt(cmax, 10, 64)
							if err != nil {
								fmt.Println(err.Error())
								max = 1
							}
							log.Println(ip + " is in range " + cidr + " and has max " + string(max) + " in " + string(timd) + " hour(s)")
						}
					}
				}
			}
		}
		// Whether the ReqData is expired
		valid := rd.start.After(time.Now().Add(-1 * interval))

		// If the valid count for this timeframe exceeds the max
		if !isExcluded && valid && rd.count >= max {

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
	ip := c.RealIP()    // The requestors IP
	rd := r.reqData(ip) // The requestors ReqData

	// Return the current rate limit standing
	return c.JSON(http.StatusOK, &LimitStatus{
		Max:      r.max,
		Interval: r.interval,
		Current:  rd.count,
	})
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
