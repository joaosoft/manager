package gomanager

import (
	"github.com/joaosoft/go-log/service"
)

var global = make(map[string]interface{})
var log = golog.NewLogDefault("go-manager", golog.InfoLevel)
