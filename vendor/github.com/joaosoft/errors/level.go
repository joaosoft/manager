package errors

type Level int

const (
	PanicLevel Level = iota // PanicLevel, when there is no recover
	FatalLevel              // FatalLevel, when the error is fatal to the application
	ErrorLevel              // ErrorLevel, when there is a controlled error
	WarnLevel               // WarnLevel, when there is a warning
	InfoLevel               // InfoLevel, when it is a informational message
	DebugLevel              // DebugLevel, when it is a debugging message
)
