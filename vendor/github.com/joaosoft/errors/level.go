package errors

type Level int

const (
	LevelPanic Level = iota // LevelPanic, when there is no recover
	LevelFatal              // LevelFatal, when the error is fatal to the application
	LevelError              // LevelError, when there is a controlled error
	LevelWarn               // LevelWarn, when there is a warning
	LevelInfo               // LevelInfo, when it is a informational message
	LevelDebug              // LevelDebug, when it is a debugging message
)
