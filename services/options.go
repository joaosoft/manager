package gomanager

import (
	"github.com/sirupsen/logrus"
)

// GoManagerOption ...
type GoManagerOption func(gomanager *GoManager)

// Reconfigure ...
func (gomanager *GoManager) Reconfigure(options ...GoManagerOption) {
	for _, option := range options {
		option(gomanager)
	}
}

// WithLogLevel ...
func WithLogLevel(level logrus.Level) GoManagerOption {
	return func(gomanager *GoManager) {
		logrus.SetLevel(level)
	}
}

// WithRunInBackground ...
func WithRunInBackground(runInBackground bool) GoManagerOption {
	return func(gomanager *GoManager) {
		gomanager.runInBackground = runInBackground
	}
}
