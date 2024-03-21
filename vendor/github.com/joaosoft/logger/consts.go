package logger

type Level int
type Prefix string

const (
	LevelDefault = LevelInfo // LevelDefault Level

	LevelNone  Level = iota // LevelNone, when the logging is disabled
	LevelPrint              // LevelPrint, when it is a system message
	LevelPanic              // LevelPanic, when there is no recover
	LevelFatal              // LevelFatal, when the error is fatal to the application
	LevelError              // LevelError, when there is a controlled error
	LevelWarn               // LevelWarn, when there is a warning
	LevelInfo               // LevelInfo, when it is a informational message
	LevelDebug              // LevelDebug, when it is a debugging message

	// Special Prefixes
	LEVEL     Prefix = "{{LEVEL}}"     // Add the level value to the prefix
	TIMESTAMP Prefix = "{{TIMESTAMP}}" // Add the timestamp value to the prefix
	DATE      Prefix = "{{DATE}}"      // Add the date value to the prefix
	TIME      Prefix = "{{TIME}}"      // Add the time value to the prefix
	IP        Prefix = "{{IP}}"        // Add the client ip address
	TRACE     Prefix = "{{TRACE}}"     // Add the error trace
	PACKAGE   Prefix = "{{PACKAGE}}"   // Add the package name
	FILE      Prefix = "{{FILE}}"      // Add the file
	FUNCTION  Prefix = "{{FUNCTION}}"  // Add the function name
	STACK     Prefix = "{{STACK}}"     // Add the debug stack
)
