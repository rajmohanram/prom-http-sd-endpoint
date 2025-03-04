package health

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rajmohanram/prom-http-sd-endpoint/internal/logger"
)

// Status represents the application health status
type Status struct {
	Status string `json:"status"`
	Uptime string `json:"uptime"`
}

// Handler is the health check handler
type Handler struct {
	startTime time.Time
}

// NewHandler creates a new health check handler
func NewHandler() http.Handler {
	return &Handler{
		startTime: time.Now(),
	}
}

// ServeHTTP handles health check requests
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := Status{
		Status: "ok",
		Uptime: time.Since(h.startTime).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(status); err != nil {
		logger.Logger.WithField("error", err).Error("Error encoding health status response")
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	logger.Logger.Debug("Health check request processed successfully")
}
