package manager

import (
	golog "github.com/joaosoft/go-log/app"
)

var global = make(map[string]interface{})
var logger = golog.NewLogDefault("manager", golog.InfoLevel)

func init() {
	global[path_key] = defaultPath
}
