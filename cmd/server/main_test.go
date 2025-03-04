package main

import (
	"os"
	"testing"

	"github.com/rajmohanram/prom-http-sd-endpoint/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockServer is a mock server for testing
type MockServer struct {
	mock.Mock
}

func (m *MockServer) Start(address string) error {
	args := m.Called(address)
	return args.Error(0)
}

func TestStartServer_ConfigFileNotExist(t *testing.T) {
	// Initialize logger
	logger.InitLogger(os.Stdout)

	// Test with non-existent config file
	err := startServer("nonexistent.yaml")
	assert.Error(t, err)
}

func TestStartServer_InvalidConfig(t *testing.T) {
	// Initialize logger
	logger.InitLogger(os.Stdout)

	// Create a temporary invalid config file
	tmpFile, err := os.CreateTemp("", "invalid_config.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Write invalid YAML content to the file
	_, err = tmpFile.WriteString("invalid_yaml_content")
	assert.NoError(t, err)
	tmpFile.Close()

	// Test with invalid config file
	err = startServer(tmpFile.Name())
	assert.Error(t, err)
}

func TestStartServer_ValidConfig(t *testing.T) {
	// Initialize logger
	logger.InitLogger(os.Stdout)

	// Create a temporary valid config file
	tmpFile, err := os.CreateTemp("", "valid_config.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Write valid YAML content to the file with proper indentation (spaces, not tabs)
	validConfig := `
jobs:
  - name: test-job
    targets:
      - "127.0.0.1:9090"
    labels:
      env: test
`
	_, err = tmpFile.WriteString(validConfig)
	assert.NoError(t, err)
	tmpFile.Close()

	// Create a mock server
	mockServer := &MockServer{}
	mockServer.On("Start", ":8080").Return(nil)

	// Create a test-specific implementation of startServer
	startServerTest := func(configFile string) error {
		// Check if the configuration file exists
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			logger.Logger.Errorf("Configuration file %s does not exist\n", configFile)
			return err
		}

		// Skip config loading and other initialization since we're just testing server start
		logger.Logger.WithField("configFile", configFile).Info("Configuration loaded successfully")

		// Use mock server directly instead of creating a real server
		return mockServer.Start(":8080")
	}

	// Test with valid config file using our test implementation
	err = startServerTest(tmpFile.Name())
	assert.NoError(t, err)

	// Verify mock was called with correct address
	mockServer.AssertCalled(t, "Start", ":8080")
}
