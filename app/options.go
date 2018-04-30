package gomanager

import (
	golog "github.com/joaosoft/go-log/app"
)

// managerOption ...
type managerOption func(manager *Manager)

// Reconfigure ...
func (manager *Manager) Reconfigure(options ...managerOption) {
	for _, option := range options {
		option(manager)
	}
}

// WithRunInBackground ...
func WithRunInBackground(runInBackground bool) managerOption {
	return func(manager *Manager) {
		manager.runInBackground = runInBackground
	}
}

// WithLogger ...
func WithLogger(logger golog.ILog) managerOption {
	return func(manager *Manager) {
		log = logger
		manager.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level golog.Level) managerOption {
	return func(manager *Manager) {
		log.SetLevel(level)
	}
}
