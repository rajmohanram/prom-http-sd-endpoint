package main

import (
	"flag"
	"os"

	"github.com/rajmohanram/prom-http-sd-endpoint/internal/config"
	"github.com/rajmohanram/prom-http-sd-endpoint/internal/discovery"
	"github.com/rajmohanram/prom-http-sd-endpoint/internal/health"
	"github.com/rajmohanram/prom-http-sd-endpoint/internal/logger"
	"github.com/rajmohanram/prom-http-sd-endpoint/pkg/server"
)

func main() {
	// Initialize logger
	logger.InitLogger(os.Stdout)

	// Log that the server is starting
	logger.Logger.Info("Starting the API server")

	// Define a command-line flag for the configuration file path
	configFile := flag.String("config", "targets.yaml", "Path to the configuration file")
	flag.Parse()

	// Start the server
	if err := startServer(*configFile); err != nil {
		logger.Logger.Fatalf("Error starting server: %v", err)
	}
}

// startServer loads the configuration and starts the server
func startServer(configFile string) error {
	logger.Logger.Info("Starting server process...")

	// Check if the configuration file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		logger.Logger.Errorf("Configuration file %s does not exist\n", configFile)
		return err
	}
	logger.Logger.Info("Configuration file exists, loading...")

	// Load and validate configuration
	configData, err := config.LoadConfig(configFile)
	if err != nil {
		logger.Logger.WithField("configFile", configFile).Errorf("Error loading config: %v", err)
		return err
	}
	logger.Logger.Info("Configuration loaded successfully")

	// Create the HTTP handler for dynamic discovery endpoints
	logger.Logger.Info("Creating discovery handler...")
	discoveryHandler := discovery.NewHandler(configData)
	logger.Logger.Info("Created discovery handler")

	// Create the health check handler
	logger.Logger.Info("Creating health check handler...")
	healthHandler := health.NewHandler()
	logger.Logger.Info("Created health check handler")

	// Initialize and start the server
	logger.Logger.Info("Creating server...")
	srv := server.NewServer(discoveryHandler, healthHandler)
	logger.Logger.Info("About to start server on port 8080...")
	if err := srv.Start(":8080"); err != nil {
		logger.Logger.WithField("error", err).Errorf("Error starting server")
		return err
	}

	return nil
}
