package server

import (
	"net/http"

	"github.com/rajmohanram/prom-http-sd-endpoint/internal/logger"
)

// Server struct to hold server-related data
type Server struct {
	mux            *http.ServeMux
	listenAndServe func(addr string, handler http.Handler) error
}

// NewServer creates a new Server instance
func NewServer(discoveryHandler http.Handler, healthHandler http.Handler) *Server {
	mux := http.NewServeMux()

	// Register the health check endpoint
	mux.Handle("/healthz", healthHandler)

	// Register the discovery handler for all other paths
	mux.Handle("/", discoveryHandler)

	return &Server{
		mux:            mux,
		listenAndServe: http.ListenAndServe,
	}
}

// Start the HTTP server
func (s *Server) Start(address string) error {
	logger.Logger.Infof("Starting server on %s...", address)
	return s.listenAndServe(address, s.mux)
}
