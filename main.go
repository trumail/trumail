package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/sdwolfe32/trumail/api"
)

func main() {
	log := logrus.New() // New Logger

	// Retrieve environment variables and initialize logger
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	sourceAddr := os.Getenv("SOURCE_ADDR")
	if sourceAddr == "" {
		sourceAddr = "admin@gmail.com"
	}
	rateLimit, _ := strconv.ParseBool(os.Getenv("RATE_LIMIT"))
	serveWeb, _ := strconv.ParseBool(os.Getenv("SERVE_WEB"))
	env := os.Getenv("ENVIRONMENT")
	if strings.Contains(env, "prod") {
		log.Formatter = new(logrus.JSONFormatter)
	}
	l := log.WithField("port", port)

	// Initialize service Builder and Lookuper
	h := retrievePTR()
	r := api.NewRouter(rateLimit)
	s := api.NewLookuper(log, h, sourceAddr)

	// Bind endpoints to the router
	l.Info("Binding all Trumail endpoints to the router")
	r.HandleEndpoint("/{format:(?:xml|json)}/{email}", s.Lookup).Methods(http.MethodGet)
	r.HandleEndpoint("/healthcheck", s.Healthcheck).Methods(http.MethodGet)

	if serveWeb {
		// Set all remaining paths to point to static files (must come after)
		r.HandleStatic("./web")
	}

	// Listen and Serve
	l.Info("Listening and Serving")
	r.ListenAndServe(port)
}

// retrievePTR attempts to retrieve the PTR record for the IP address retrieved
// via an API call on api.ipify.org
func retrievePTR() string {
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

	// Retrieves the PTR record for our IP and returns without the trailing dot
	names, err := net.LookupAddr(string(data))
	if err != nil {
		return string(data)
	}
	return strings.TrimSuffix(names[0], ".")
}
