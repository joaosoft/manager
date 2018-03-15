package services

import (
	"os"

	logger "github.com/sirupsen/logrus"
)

var global = make(map[string]interface{})
var log = logger.WithFields(logger.Fields{
	"application": "go-manager",
})

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logger.SetFormatter(&logger.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logger.SetLevel(logger.DebugLevel)
}
