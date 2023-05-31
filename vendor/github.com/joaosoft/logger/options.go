package logger

import (
	"fmt"
	"io"

	writer "github.com/joaosoft/writers"
)

// LoggerOption ...
type LoggerOption func(log *Logger)

// Reconfigure ...
func (logger *Logger) Reconfigure(options ...LoggerOption) {
	for _, option := range options {
		option(logger)
	}
}

// WithWriter ...
func WithWriter(writer ...io.Writer) LoggerOption {
	return func(logger *Logger) {
		logger.writer = append(logger.writer, writer...)
	}
}

// WithSpecialWriter ...
func WithSpecialWriter(writer ...ISpecialWriter) LoggerOption {
	return func(logger *Logger) {
		logger.specialWriter = append(logger.specialWriter, writer...)
	}
}

// WithLevel ...
func WithLevel(level Level) LoggerOption {
	return func(logger *Logger) {
		logger.level = level
	}
}

// WithFormatHandler ...
func WithFormatHandler(formatHandler writer.FormatHandler) LoggerOption {
	return func(logger *Logger) {
		logger.formatHandler = formatHandler
	}
}

func WithOptions(prefixes, tags, fields map[string]interface{}) LoggerOption {
	return func(logger *Logger) {
		logger.prefixes = prefixes
		logger.tags = tags
		logger.fields = fields
	}
}

func WithOptPrefixes(prefixes map[string]interface{}) LoggerOption {
	return func(logger *Logger) {
		logger.prefixes = prefixes
	}
}

func WithOptTags(tags map[string]interface{}) LoggerOption {
	return func(logger *Logger) {
		logger.tags = tags
	}
}

func WithOptFields(fields map[string]interface{}) LoggerOption {
	return func(logger *Logger) {
		logger.fields = fields
	}
}

func WithOptPrefix(key string, value interface{}) LoggerOption {
	return func(logger *Logger) {
		logger.prefixes[key] = fmt.Sprintf("%s", value)
	}
}

func WithOptTag(key string, value interface{}) LoggerOption {
	return func(logger *Logger) {
		logger.tags[key] = fmt.Sprintf("%s", value)
	}
}

func WithOptField(key string, value interface{}) LoggerOption {
	return func(logger *Logger) {
		logger.fields[key] = fmt.Sprintf("%s", value)
	}
}
