package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_FileNotExist(t *testing.T) {
	_, err := LoadConfig("nonexistent.yaml")
	assert.Error(t, err)
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	// Create a temporary invalid config file
	tmpFile, err := os.CreateTemp("", "invalid_config.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("invalid_yaml_content")
	assert.NoError(t, err)
	tmpFile.Close()

	_, err = LoadConfig(tmpFile.Name())
	assert.Error(t, err)
}

func TestLoadConfig_NoJobsDefined(t *testing.T) {
	// Create a temporary config file with no jobs
	tmpFile, err := os.CreateTemp("", "no_jobs_config.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	noJobsConfig := `
jobs: []
`
	_, err = tmpFile.WriteString(noJobsConfig)
	assert.NoError(t, err)
	tmpFile.Close()

	_, err = LoadConfig(tmpFile.Name())
	assert.Error(t, err)
}

func TestLoadConfig_ValidConfig(t *testing.T) {
	// Create a temporary valid config file
	tmpFile, err := os.CreateTemp("", "valid_config.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

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

	config, err := LoadConfig(tmpFile.Name())
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Len(t, config.Jobs, 1)
	assert.Equal(t, "test-job", config.Jobs[0].Name)
	assert.Equal(t, []string{"127.0.0.1:9090"}, config.Jobs[0].Targets)
	assert.Equal(t, map[string]string{"env": "test"}, config.Jobs[0].Labels)
}

func TestValidateConfig_ValidConfig(t *testing.T) {
	config := &Config{
		Jobs: []Job{
			{
				Name:    "test-job",
				Targets: []string{"127.0.0.1:9090"},
				Labels:  map[string]string{"env": "test"},
			},
		},
	}

	err := ValidateConfig(config)
	assert.NoError(t, err)
}

func TestValidateConfig_InvalidConfig(t *testing.T) {
	config := &Config{
		Jobs: []Job{
			{
				Name:    "test-job",
				Targets: []string{"invalid_target"},
				Labels:  map[string]string{"env": "test"},
			},
		},
	}

	err := ValidateConfig(config)
	assert.Error(t, err)
}
