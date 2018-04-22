package main

import (
	"net/http"
	"strings"
	"time"

	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sdwolfe32/trumail/api"
	"github.com/sdwolfe32/trumail/config"
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
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://trumail.io"},
		AllowMethods: []string{http.MethodGet},
	}))
	s := api.NewService(logger, config.SourceAddr, config.HTTPClientTimeout)

	// Bind endpoints to router
	l.Info("Binding API endpoints to the router")
	if config.RateLimitHours != 0 && config.RateLimitMax != 0 {
		r := api.NewRateLimiter(config.Token, config.RateLimitMax,
			time.Hour*time.Duration(config.RateLimitHours))
		e.GET("/v1/:format/:email", s.Lookup, r.RateLimit)
		e.GET("/v1/limit-status", r.LimitStatus)
	} else {
		e.GET("/v1/:format/:email", s.Lookup)
	}
	e.GET("/v1/health", s.Health)
	e.GET("/v1/debug", api.Debug)

	// Listen and Serve
	l.WithField("port", config.Port).Info("Listening and Serving")
	l.Fatal(e.Start(":" + config.Port))
}
