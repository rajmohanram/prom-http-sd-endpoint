package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rajmohanram/prom-http-sd-endpoint/internal/logger"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

type Job struct {
	Name    string            `yaml:"name" json:"name"`
	Targets []string          `yaml:"targets" json:"targets"`
	Labels  map[string]string `yaml:"labels" json:"labels"`
}

type Config struct {
	Jobs []Job `yaml:"jobs" json:"jobs"`
}

// LoadConfig reads the configuration file, unmarshals it into a Config struct, and validates it
func LoadConfig(filename string) (*Config, error) {
	logger.Logger.WithField("filename", filename).Info("Reading configuration file")
	file, err := os.ReadFile(filename)
	if err != nil {
		logger.Logger.WithField("filename", filename).Errorf("Error reading file: %v", err)
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		logger.Logger.WithField("filename", filename).Errorf("Error unmarshalling YAML: %v", err)
		return nil, fmt.Errorf("error unmarshalling YAML: %v", err)
	}

	if len(config.Jobs) == 0 {
		logger.Logger.WithField("filename", filename).Error("No jobs defined in the configuration")
		return nil, fmt.Errorf("no jobs defined in the configuration")
	}

	if err := ValidateConfig(&config); err != nil {
		logger.Logger.WithField("filename", filename).Errorf("Configuration validation failed: %v", err)
		return nil, err
	}

	logger.Logger.WithField("filename", filename).Info("Configuration loaded and validated successfully")
	return &config, nil
}

// ValidateConfig validates the configuration data with a JSON schema
func ValidateConfig(config *Config) error {
	logger.Logger.Info("Validating configuration")
	// Convert the config to JSON
	configJSON, err := json.Marshal(config)
	if err != nil {
		logger.Logger.Errorf("Error marshalling config to JSON: %v", err)
		return fmt.Errorf("error marshalling config to JSON: %v", err)
	}

	// Load the JSON schema
	schemaLoader := gojsonschema.NewStringLoader(jsonSchema)
	documentLoader := gojsonschema.NewBytesLoader(configJSON)

	// Validate the JSON against the schema
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		logger.Logger.Errorf("Error validating config: %v", err)
		return fmt.Errorf("error validating config: %v", err)
	}

	if !result.Valid() {
		for _, desc := range result.Errors() {
			logger.Logger.Errorf("Validation error: %s", desc)
		}
		return fmt.Errorf("configuration validation failed")
	}

	logger.Logger.Info("Configuration validation successful")
	return nil
}

// JSON schema for validation
const jsonSchema = `
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "jobs": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "name": {
            "type": "string"
          },
          "targets": {
            "type": "array",
            "items": {
              "type": "string",
              "pattern": "^([a-zA-Z0-9.-]+|\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}):\\d+$"
            }
          },
          "labels": {
            "type": "object",
            "additionalProperties": {
              "type": "string"
            }
          }
        },
        "required": ["name", "targets"]
      }
    }
  },
  "required": ["jobs"]
}
`
