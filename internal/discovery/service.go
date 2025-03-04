package discovery

import (
	"fmt"

	"github.com/rajmohanram/prom-http-sd-endpoint/internal/config"
	"github.com/rajmohanram/prom-http-sd-endpoint/internal/logger"
)

// Generates a discovery response for the given job
func GenerateDiscoveryResponse(job config.Job) ([]map[string]interface{}, error) {
	if len(job.Targets) == 0 {
		err := fmt.Errorf("no targets found for job: %s", job.Name)
		logger.Logger.WithField("jobName", job.Name).Error(err)
		return nil, err
	}

	logger.Logger.WithField("jobName", job.Name).Info("Generating discovery response")
	return []map[string]interface{}{
		{
			"targets": job.Targets,
			"labels":  job.Labels,
		},
	}, nil
}
