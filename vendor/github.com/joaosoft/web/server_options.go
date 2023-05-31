package web

import (
	"github.com/joaosoft/logger"
)

// ServerOption ...
type ServerOption func(builder *Server)

// Reconfigure ...
func (w *Server) Reconfigure(options ...ServerOption) {
	for _, option := range options {
		option(w)
	}
}

// WithServerName ...
func WithServerName(name string) ServerOption {
	return func(webserver *Server) {
		webserver.name = name
	}
}

// WithServerConfiguration ...
func WithServerConfiguration(config *ServerConfig) ServerOption {
	return func(webserver *Server) {
		webserver.config = config
	}
}

// WithServerLogger ...
func WithServerLogger(logger logger.ILogger) ServerOption {
	return func(webserver *Server) {
		webserver.logger = logger
		webserver.isLogExternal = true
	}
}

// WithServerLogLevel ...
func WithServerLogLevel(level logger.Level) ServerOption {
	return func(webserver *Server) {
		webserver.logger.SetLevel(level)
	}
}

// WithServerAddress ...
func WithServerAddress(address string) ServerOption {
	return func(webserver *Server) {
		webserver.config.Address = address
	}
}

// WithServerMultiAttachmentMode ...
func WithServerMultiAttachmentMode(mode MultiAttachmentMode) ServerOption {
	return func(webserver *Server) {
		webserver.multiAttachmentMode = mode
	}
}
