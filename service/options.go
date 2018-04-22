package gomanager

import (
	"github.com/joaosoft/go-log/service"
)

// ManagerOption ...
type ManagerOption func(manager *Manager)

// Reconfigure ...
func (manager *Manager) Reconfigure(options ...ManagerOption) {
	for _, option := range options {
		option(manager)
	}
}

// WithRunInBackground ...
func WithRunInBackground(runInBackground bool) ManagerOption {
	return func(manager *Manager) {
		manager.runInBackground = runInBackground
	}
}

// WithLogger ...
func WithLogger(logger golog.ILog) ManagerOption {
	return func(manager *Manager) {
		log = logger
		manager.logIsExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level golog.Level) ManagerOption {
	return func(manager *Manager) {
		log.SetLevel(level)
	}
}
