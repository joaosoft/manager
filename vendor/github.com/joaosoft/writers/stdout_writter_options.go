package writers

import (
	"time"
)

// StdoutWriterOption ...
type StdoutWriterOption func(fileWriter *StdoutWriter)

// Reconfigure ...
func (stdoutWriter *StdoutWriter) Reconfigure(options ...StdoutWriterOption) {
	for _, option := range options {
		option(stdoutWriter)
	}
}

// WithStdoutFlushTime ...
func WithStdoutFlushTime(flushTime time.Duration) StdoutWriterOption {
	return func(stdoutWriter *StdoutWriter) {
		stdoutWriter.config.flushTime = flushTime
	}
}

// WithStdoutQuitChannel ...
func WithStdoutQuitChannel(quit chan bool) StdoutWriterOption {
	return func(stdoutWriter *StdoutWriter) {
		stdoutWriter.quit = quit
	}
}

// WithStdoutFormatHandler ...
func WithStdoutFormatHandler(formatHandler FormatHandler) StdoutWriterOption {
	return func(stdoutWriter *StdoutWriter) {
		stdoutWriter.formatHandler = formatHandler
	}
}
