package gomanager

import (
	logger "github.com/joaosoft/go-log/service"
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

// WithLogger ...
func WithLogger(logger logger.Log) GoManagerOption {
	return func(gomanager *GoManager) {
		log = logger
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) GoManagerOption {
	return func(gomanager *GoManager) {
		log.SetLevel(level)
	}
}
