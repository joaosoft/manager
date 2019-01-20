package manager

import (
	logger "github.com/joaosoft/logger"
)

var global = make(map[string]interface{})
var log = logger.NewLogDefault("Manager", logger.InfoLevel)

func init() {
	global[path_key] = defaultPath
}
