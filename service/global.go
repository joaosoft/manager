package gomanager

import (
	"os"

	logger "github.com/joaosoft/go-log/service"
	writer "github.com/joaosoft/go-writer/service"
)

var global = make(map[string]interface{})
var log = logger.NewLog(
	logger.WithLevel(logger.InfoLevel),
	logger.WithFormatHandler(writer.JsonFormatHandler),
	logger.WithWriter(os.Stdout)).WithPrefixes(map[string]interface{}{
	"level":   logger.LEVEL,
	"time":    logger.TIME,
	"service": "go-manager"})
