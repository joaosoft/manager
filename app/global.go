package gomanager

import (
	golog "github.com/joaosoft/go-log/app"
)

var global = make(map[string]interface{})
var log = golog.NewLogDefault("go-manager", golog.InfoLevel)

func init() {
	global[path_key] = defaultPath
}
