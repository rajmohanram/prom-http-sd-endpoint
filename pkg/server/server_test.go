package server

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/rajmohanram/prom-http-sd-endpoint/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHandler is a mock implementation of http.Handler
type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func TestNewServer(t *testing.T) {
	// Create mock handlers
	mockDiscoveryHandler := &MockHandler{}
	mockHealthHandler := &MockHandler{}

	// Create a new server with the mock handlers
	srv := NewServer(mockDiscoveryHandler, mockHealthHandler)

	// Assert that the server is not nil
	assert.NotNil(t, srv)
}

func TestStartServer(t *testing.T) {
	// Initialize logger with a buffer to capture log output
	var buf bytes.Buffer
	logger.InitLogger(&buf)

	// Create mock handlers
	mockDiscoveryHandler := &MockHandler{}
	mockHealthHandler := &MockHandler{}

	// Create a new server with the mock handlers
	srv := NewServer(mockDiscoveryHandler, mockHealthHandler)

	// Set up a mock ListenAndServe function
	called := false
	srv.listenAndServe = func(addr string, handler http.Handler) error {
		called = true
		assert.Equal(t, ":8080", addr)
		return nil
	}

	// Call the Start method
	err := srv.Start(":8080")

	// Assert that no error was returned
	assert.NoError(t, err)

	// Assert that our mock function was called
	assert.True(t, called)

	// Assert that the log message was output correctly
	assert.Contains(t, buf.String(), "Starting server on :8080")
}
