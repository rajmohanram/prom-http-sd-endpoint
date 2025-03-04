package discovery

import (
    "testing"

    "github.com/rajmohanram/prom-http-sd-endpoint/internal/config"
    "github.com/stretchr/testify/assert"
)

func TestGenerateDiscoveryResponse_NoTargets(t *testing.T) {
    // Create a job with no targets
    job := config.Job{
        Name:    "test-job",
        Targets: []string{},
        Labels:  map[string]string{"env": "test"},
    }

    // Call the GenerateDiscoveryResponse function
    response, err := GenerateDiscoveryResponse(job)

    // Check that an error is returned
    assert.Error(t, err)
    assert.Nil(t, response)
}

func TestGenerateDiscoveryResponse_WithTargets(t *testing.T) {
    // Create a job with targets
    job := config.Job{
        Name:    "test-job",
        Targets: []string{"127.0.0.1:9090"},
        Labels:  map[string]string{"env": "test"},
    }

    // Call the GenerateDiscoveryResponse function
    response, err := GenerateDiscoveryResponse(job)

    // Check that no error is returned
    assert.NoError(t, err)

    // Check the response
    expectedResponse := []map[string]interface{}{
        {
            "targets": []string{"127.0.0.1:9090"},
            "labels":  map[string]string{"env": "test"},
        },
    }
    assert.Equal(t, expectedResponse, response)
}
