package logger

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Initialize the logger with the buffer as output
	InitLogger(&buf)

	// Check that the logger is set to output to the buffer
	assert.Equal(t, &buf, Logger.Out)

	// Check that the log level is set to Info
	assert.Equal(t, logrus.InfoLevel, Logger.Level)

	// Check that the log format is set to JSON
	_, ok := Logger.Formatter.(*logrus.JSONFormatter)
	assert.True(t, ok)

	// Log a test message
	Logger.Info("test message")

	// Check that the log message is written to the buffer
	assert.Contains(t, buf.String(), "test message")
}
