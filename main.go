package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/sdwolfe32/trumail/api"
	"github.com/sdwolfe32/trumail/config"
)

func main() {
	http.DefaultClient = &http.Client{
		Timeout: time.Duration(config.HTTPClientTimeout) * time.Second,
	}

	logger := logrus.New() // New Logger

	if strings.Contains(config.Env, "prod") {
		logger.Formatter = new(logrus.JSONFormatter)
	}
	l := logger.WithField("port", config.Port)

	r, s := api.Initialize(logger)

	l.Info("Binding all Trumail endpoints to the router")
	api.RegisterEndpoints(r, s)

	if config.ServeWeb {
		// Set all remaining paths to point to static files (must come after)
		r.HandleStatic("./web")
	}

	// Listen and Serve
	l.Info("Listening and Serving")
	r.ListenAndServe(config.Port)
}
