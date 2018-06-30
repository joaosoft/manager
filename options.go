package manager

import "github.com/joaosoft/logger"

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
func WithLogger(logger logger.ILogger) ManagerOption {
	return func(manager *Manager) {
		logger = logger
		manager.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) ManagerOption {
	return func(manager *Manager) {
		log.SetLevel(level)
	}
}
