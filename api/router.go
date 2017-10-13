package api

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/sync/syncmap"
)

// Router defines all functionality for our api service router
type Router interface {
	HandleStatic(path string) *mux.Route
	HandleEndpoint(pattern string, endpoint Endpoint) *mux.Route
	ListenAndServe(port string) error
}

// Router contains an embedded router that we can use to bind
// Endpoints to
type router struct {
	router    *mux.Router
	rateLimit bool
	rateMap   *syncmap.Map
}

// NewRouter generates a new Router that will be used to bind
// handlers to the *mux.Router
func NewRouter(rateLimit bool) Router {
	return &router{
		router:    mux.NewRouter(),
		rateLimit: rateLimit,
		rateMap:   &syncmap.Map{}, // ip-address -> last request time
	}
}

// HandleStatic binds a new fileserver using the passed path to the router
func (r *router) HandleStatic(path string) *mux.Route {
	return r.router.PathPrefix("/").Handler(http.FileServer(http.Dir(path)))
}

// HandleEndpoint binds a new Endpoint handler to the router
func (r *router) HandleEndpoint(pattern string, endpoint Endpoint) *mux.Route {
	return r.router.HandleFunc(pattern, r.endpointWrapper(endpoint))
}

// ListenAndServe applies CORS headers and starts the server
// using the embedded router
func (r *router) ListenAndServe(port string) error {
	// Create the basic HTTP server with base parameters
	srv := &http.Server{
		Handler:      r.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Apply CORS headers
	srv.Handler = handlers.CORS(
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
		handlers.AllowedMethods([]string{"POST", "PUT", "GET", "OPTIONS", "HEAD"}),
		handlers.AllowedOrigins([]string{"*"}),
	)(r.router)

	// Set the port to run on and serve
	srv.Addr = ":" + port
	return srv.ListenAndServe()
}

// An Endpoint is a service endpoint that receives a request and returns either
// a successfully processed response body or an Error. In either case both
// responses are encoded and returned to the user with the appropriate status
// code
type Endpoint func(*http.Request) (interface{}, error)

// endpointWrapper transforms an endpoint into a standard http.Handlerfunc
func (r *router) endpointWrapper(e Endpoint) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		// Extract the request params
		ipAddr := req.Header.Get("X-Forwarded-For")
		format := strings.ToLower(mux.Vars(req)["format"])

		// Execute rate limit logic if the server is configured to do so
		if r.rateLimit {
			if last, ok := r.rateMap.Load(ipAddr); ok {
				if last.(time.Time).Add(time.Second).After(time.Now()) {
					respond(rw, http.StatusTooManyRequests, format,
						NewError("You have exceeded the rate-limit", http.StatusTooManyRequests, nil))
					return
				}
			}
			// Defer store fresh rate-limit time
			defer r.rateMap.Store(ipAddr, time.Now())
		}

		// Handle the request and respond appropriately
		res, err := e(req)
		if err != nil {
			if e, ok := err.(*Error); ok {
				respond(rw, e.StatusCode, format, e)
			} else {
				respond(rw, http.StatusInternalServerError, format,
					NewError("An error has occurred", http.StatusInternalServerError, err))
			}
			return
		}
		respond(rw, http.StatusOK, format, res)
	}
}

// respond is responsible for encoding the response and writing the desired
// format to the http.ResponseWriter
func respond(w http.ResponseWriter, status int, format string, res interface{}) {
	// Encode the response using the format passed
	switch format {
	case "xml":
		encodeXML(w, status, res)
	case "json":
		encodeJSON(w, status, res)
	default:
		encodeJSON(w, status, res)
	}
}

// encodeXML encodes the response to XML and writes it to the ResponseWriter
func encodeXML(w http.ResponseWriter, status int, res interface{}) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(status)
	xml.NewEncoder(w).Encode(res)
}

// encodeJSON encodes the response to JSON and writes it to the ResponseWriter
func encodeJSON(w http.ResponseWriter, status int, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}
