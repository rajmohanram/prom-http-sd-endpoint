package discovery

import (
	"encoding/json"
	"net/http"

	"github.com/rajmohanram/prom-http-sd-endpoint/internal/config"
	"github.com/rajmohanram/prom-http-sd-endpoint/internal/logger"
)

// Handler struct to hold the loaded configuration
type Handler struct {
	config *config.Config
}

// NewHandler creates a new Handler for discovery endpoints
func NewHandler(config *config.Config) *Handler {
	return &Handler{config: config}
}

// ServeHTTP is the main entry point for the HTTP request
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jobName := r.URL.Path[1:] // Remove the leading slash

	// Search for the job in the configuration
	for _, job := range h.config.Jobs {
		if job.Name == jobName {
			// Respond with Prometheus service discovery JSON format
			sdResponse := []map[string]interface{}{
				{
					"targets": job.Targets,
					"labels":  job.Labels,
				},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(sdResponse); err != nil {
				logger.Logger.WithField("error", err).Error("Error encoding JSON response")
				http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			} else {
				logger.Logger.WithField("jobName", jobName).Info("Successfully served discovery endpoint")
			}
			return
		}
	}

	// If Target not found, return 404
	logger.Logger.WithField("jobName", jobName).Warn("Target not found")
	http.Error(w, "Target not found", http.StatusNotFound)
}
