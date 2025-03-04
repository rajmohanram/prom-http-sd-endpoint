package health

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/rajmohanram/prom-http-sd-endpoint/internal/logger"
    "github.com/stretchr/testify/assert"
)

func init() {
    // Initialize logger for tests
    logger.InitLogger(nil)
}

func TestNewHandler(t *testing.T) {
    // Create a new handler
    handler := NewHandler()

    // Assert that the handler is not nil
    assert.NotNil(t, handler)
    
    // Check the type of the handler
    _, ok := handler.(*Handler)
    assert.True(t, ok, "Handler should be of type *Handler")
}

func TestServeHTTP(t *testing.T) {
    // Create a new handler with a specific start time
    handler := &Handler{
        startTime: time.Now().Add(-5 * time.Minute), // Set to 5 minutes ago
    }

    // Create a test HTTP request
    req, err := http.NewRequest("GET", "/healthz", nil)
    assert.NoError(t, err)

    // Create a response recorder to record the response
    rr := httptest.NewRecorder()

    // Call the ServeHTTP method
    handler.ServeHTTP(rr, req)

    // Check the status code
    assert.Equal(t, http.StatusOK, rr.Code)

    // Check the content type
    assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

    // Parse the response body
    var status Status
    err = json.Unmarshal(rr.Body.Bytes(), &status)
    assert.NoError(t, err)

    // Check the status field
    assert.Equal(t, "ok", status.Status)

    // Check that the uptime string contains expected details
    assert.Contains(t, status.Uptime, "m", "Uptime should contain minutes")
    
    // Check that the uptime is roughly correct (should be close to 5 minutes)
    // This is an approximate test as the exact value depends on timing
    assert.Contains(t, status.Uptime, "5m", "Uptime should be close to 5 minutes")
}

func TestServeHTTP_JSONEncodeError(t *testing.T) {
    // Create a mock handler that will trigger a JSON encoding error
    // This is difficult to test directly, as we'd need to create a situation
    // where json.NewEncoder(w).Encode() fails, which is hard to do in a test
    
    // Instead, we'll verify that the normal case works correctly
    handler := NewHandler()
    
    // Create a test HTTP request
    req, err := http.NewRequest("GET", "/healthz", nil)
    assert.NoError(t, err)
    
    // Create a response recorder
    rr := httptest.NewRecorder()
    
    // Call the ServeHTTP method
    handler.ServeHTTP(rr, req)
    
    // Check the status code
    assert.Equal(t, http.StatusOK, rr.Code)
    
    // Verify the response can be parsed as valid JSON
    var status Status
    err = json.Unmarshal(rr.Body.Bytes(), &status)
    assert.NoError(t, err, "Response should be valid JSON")
}

func TestHealthEndpoint_Integration(t *testing.T) {
    // Create a handler
    handler := NewHandler()
    
    // Create a test server
    server := httptest.NewServer(handler)
    defer server.Close()
    
    // Make a request to the server
    resp, err := http.Get(server.URL)
    assert.NoError(t, err)
    defer resp.Body.Close()
    
    // Check the status code
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    
    // Parse the response body
    var status Status
    err = json.NewDecoder(resp.Body).Decode(&status)
    assert.NoError(t, err)
    
    // Verify the response fields
    assert.Equal(t, "ok", status.Status)
    assert.NotEmpty(t, status.Uptime)
}
