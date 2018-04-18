package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/labstack/echo"
	"github.com/sdwolfe32/trumail/api"
	"github.com/sdwolfe32/trumail/config"
	"github.com/sdwolfe32/trumail/heroku"
	"github.com/sdwolfe32/trumail/verifier"
	"github.com/sirupsen/logrus"
)

func main() {
	// Generate a new logrus logger
	logger := logrus.New()
	logger.Level = logrus.DebugLevel

	// Configure the logger based on the environment
	if strings.Contains(config.Env, "prod") {
		logger.Formatter = new(logrus.JSONFormatter)
		logger.Level = logrus.InfoLevel
	}
	l := logger.WithField("port", config.Port)

	// Define all required dependencies
	l.Info("Defining all service dependencies")
	hostname := retrievePTR()
	e := echo.New()
	v := verifier.NewVerifier(hostname, config.SourceAddr)
	// Restart Dyno if officially confirmed blacklisted
	if v.Blacklisted() {
		l.Info("Confirmed Blacklisted! - Restarting Dyno")
		go log.Println(heroku.RestartDyno())
	}
	s := api.NewTrumailAPI(logger,
		time.Duration(config.HTTPClientTimeout)*time.Second, v)

	// Bind endpoints to router
	l.Info("Binding API endpoints to the router")
	if config.RateLimitHours != 0 && config.RateLimitMax != 0 {
		r := api.NewRateLimiter(config.Token, config.RateLimitMax,
			time.Hour*time.Duration(config.RateLimitHours))
		e.GET("/:format/:email", s.Lookup, r.RateLimit)
		e.GET("/limit-status", r.LimitStatus)
	} else {
		e.GET("/:format/:email", s.Lookup)
	}
	e.GET("/stats", s.Stats)
	e.GET("/health", s.Health)

	// Host static demo pages if configured to do so
	if config.ServeWeb {
		l.Info("Serving web UI on index")
		e.Static("/", "web")
		e.Static("/assets", "web/assets")
	}

	// Listen and Serve
	l.WithField("port", config.Port).Info("Listening and Serving")
	l.Fatal(e.Start(":" + config.Port))
}

// retrievePTR attempts to retrieve the PTR record for the IP
// address retrieved via an API call on api.ipify.org
func retrievePTR() string {
	// Request the IP from ipify
	resp, err := http.Get("https://api.ipify.org/")
	if err != nil {
		log.Fatal("Failed to retrieve IP from api.ipify.org")
	}
	defer resp.Body.Close()

	// Decodes the IP response body
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read IP response body")
	}

	// Retrieve the PTR record for our IP and return without a trailing dot
	names, err := net.LookupAddr(string(data))
	if err != nil {
		return string(data)
	}
	return strings.TrimSuffix(names[0], ".")
}
