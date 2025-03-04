package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// Logger is a global variable for the logger
var Logger = logrus.New()

// InitLogger initializes the logger with basic configurations
func InitLogger(output io.Writer) {
	// Log to the specified output (stdout by default)
	if output == nil {
		output = os.Stdout
	}
	Logger.SetOutput(output)

	// Set log level, for example to Info
	Logger.SetLevel(logrus.InfoLevel)

	// Set log format to JSON for structured logs
	Logger.SetFormatter(&logrus.JSONFormatter{})
}
