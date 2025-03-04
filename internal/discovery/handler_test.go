package discovery

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/rajmohanram/prom-http-sd-endpoint/internal/config"
    "github.com/stretchr/testify/assert"
)

func TestServeHTTP_JobFound(t *testing.T) {
    // Create a sample configuration with one job
    cfg := &config.Config{
        Jobs: []config.Job{
            {
                Name:    "test-job",
                Targets: []string{"127.0.0.1:9090"},
                Labels:  map[string]string{"env": "test"},
            },
        },
    }

    // Create a new handler with the sample configuration
    handler := NewHandler(cfg)

    // Create a new HTTP request
    req, err := http.NewRequest("GET", "/test-job", nil)
    assert.NoError(t, err)

    // Create a new HTTP response recorder
    rr := httptest.NewRecorder()

    // Call the ServeHTTP method
    handler.ServeHTTP(rr, req)

    // Check the status code
    assert.Equal(t, http.StatusOK, rr.Code)

    // Check the response body
    expectedResponse := `[{"targets":["127.0.0.1:9090"],"labels":{"env":"test"}}]`
    assert.JSONEq(t, expectedResponse, rr.Body.String())
}

func TestServeHTTP_JobNotFound(t *testing.T) {
    // Create a sample configuration with one job
    cfg := &config.Config{
        Jobs: []config.Job{
            {
                Name:    "test-job",
                Targets: []string{"127.0.0.1:9090"},
                Labels:  map[string]string{"env": "test"},
            },
        },
    }

    // Create a new handler with the sample configuration
    handler := NewHandler(cfg)

    // Create a new HTTP request
    req, err := http.NewRequest("GET", "/nonexistent-job", nil)
    assert.NoError(t, err)

    // Create a new HTTP response recorder
    rr := httptest.NewRecorder()

    // Call the ServeHTTP method
    handler.ServeHTTP(rr, req)

    // Check the status code
    assert.Equal(t, http.StatusNotFound, rr.Code)

    // Check the response body
    expectedResponse := "Target not found\n"
    assert.Equal(t, expectedResponse, rr.Body.String())
}
