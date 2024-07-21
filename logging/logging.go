package logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	// Configure the logger (format, output, level, etc.)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)       // Output to stdout, you can set it to a file or any io.Writer
	log.SetLevel(logrus.InfoLevel) // Default level is Info, change it as needed
}

// GetLogger provides a global logger instance
func GetLogger() *logrus.Logger {
	return log
}
