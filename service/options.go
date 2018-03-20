package gomanager

import (
	"github.com/joaosoft/go-log/service"
)

// GoManagerOption ...
type GoManagerOption func(gomanager *GoManager)

// Reconfigure ...
func (gomanager *GoManager) Reconfigure(options ...GoManagerOption) {
	for _, option := range options {
		option(gomanager)
	}
}

// WithRunInBackground ...
func WithRunInBackground(runInBackground bool) GoManagerOption {
	return func(gomanager *GoManager) {
		gomanager.runInBackground = runInBackground
	}
}

// WithRunInBackground ...
func WithLogger(logger golog.Log) GoManagerOption {
	return func(gomanager *GoManager) {
		log = logger
	}
}

// WithLogLevel ...
func WithLogLevel(level golog.Level) GoManagerOption {
	return func(gomanager *GoManager) {
		log.SetLevel(level)
	}
}
