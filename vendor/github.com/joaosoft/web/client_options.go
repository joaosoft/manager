package web

import (
	"github.com/joaosoft/logger"
)

// ClientOption ...
type ClientOption func(builder *Client)

// Reconfigure ...
func (c *Client) Reconfigure(options ...ClientOption) {
	for _, option := range options {
		option(c)
	}
}

// WithClientConfiguration ...
func WithClientConfiguration(config *ClientConfig) ClientOption {
	return func(webclient *Client) {
		webclient.config = config
	}
}

// WithClientLogger ...
func WithClientLogger(logger logger.ILogger) ClientOption {
	return func(webclient *Client) {
		webclient.logger = logger
		webclient.isLogExternal = true
	}
}

// WithClientLogLevel ...
func WithClientLogLevel(level logger.Level) ClientOption {
	return func(webclient *Client) {
		webclient.logger.SetLevel(level)
	}
}

// WithClientMultiAttachmentMode ...
func WithClientMultiAttachmentMode(mode MultiAttachmentMode) ClientOption {
	return func(webclient *Client) {
		webclient.multiAttachmentMode = mode
	}
}
