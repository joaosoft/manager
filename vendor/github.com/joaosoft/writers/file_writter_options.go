package writers

import "time"

// StdoutWriterOption ...
type FileWriterOption func(fileWriter *FileWriter)

// Reconfigure ...
func (fileWriter *FileWriter) Reconfigure(options ...FileWriterOption) {
	for _, option := range options {
		option(fileWriter)
	}
}

// WithFileDirectory ...
func WithFileDirectory(directory string) FileWriterOption {
	return func(fileWriter *FileWriter) {
		fileWriter.config.directory = directory
	}
}

// WithFileName ...
func WithFileName(fileName string) FileWriterOption {
	return func(fileWriter *FileWriter) {
		fileWriter.config.fileName = fileName
	}
}

// WithFileMaxMegaByteSize ...
func WithFileMaxMegaByteSize(fileMaxSize int64) FileWriterOption {
	return func(fileWriter *FileWriter) {
		fileWriter.config.fileMaxSize = fileMaxSize * MB_IN_BYTE
	}
}

// WithFileFlushTime ...
func WithFileFlushTime(flushTime time.Duration) FileWriterOption {
	return func(fileWriter *FileWriter) {
		fileWriter.config.flushTime = flushTime
	}
}

// WithFileQuitChannel ...
func WithFileQuitChannel(quit chan bool) FileWriterOption {
	return func(fileWriter *FileWriter) {
		fileWriter.quit = quit
	}
}

// WithFileFormatHandler ...
func WithFileFormatHandler(formatHandler FormatHandler) FileWriterOption {
	return func(fileWriter *FileWriter) {
		fileWriter.formatHandler = formatHandler
	}
}
