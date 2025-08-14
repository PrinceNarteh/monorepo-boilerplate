package routers

import (
	"net/http"

	"github.com/rs/zerolog"
)

// Router represents the HTTP router
type Router struct {
	mux    *http.ServeMux
	logger *zerolog.Logger
}

// New creates a new router instance
func New(logger *zerolog.Logger) *Router {
	return &Router{
		mux:    http.NewServeMux(),
		logger: logger,
	}
}

// SetupRoutes sets up all the routes for the application
func (r *Router) SetupRoutes() {
	// Health check endpoint
	r.mux.HandleFunc("GET /health", r.healthCheckHandler)
	
	// API routes can be added here
	r.mux.HandleFunc("GET /api/v1/status", r.statusHandler)
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// healthCheckHandler handles health check requests
func (r *Router) healthCheckHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","service":"go-boilerplate"}`))
}

// statusHandler handles status requests
func (r *Router) statusHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"running","version":"1.0.0"}`))
}
