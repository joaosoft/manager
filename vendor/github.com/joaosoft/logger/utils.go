package logger

import (
	"fmt"
	"strings"
)

func ParseLevel(level string) (Level, error) {
	switch strings.ToLower(level) {
	case "print":
		return LevelPrint, nil
	case "panic":
		return LevelPanic, nil
	case "fatal":
		return LevelFatal, nil
	case "error":
		return LevelError, nil
	case "warn":
		return LevelWarn, nil
	case "info":
		return LevelInfo, nil
	case "debug":
		return LevelDebug, nil
	case "none":
		return LevelNone, nil
	default:
		return LevelDefault, fmt.Errorf("invalid level: %s, set default level: %s", level, LevelDefault)
	}
}

func (level Level) String() string {
	switch level {
	case LevelPrint:
		return "print"
	case LevelPanic:
		return "panic"
	case LevelFatal:
		return "fatal"
	case LevelError:
		return "error"
	case LevelWarn:
		return "warn"
	case LevelInfo:
		return "info"
	case LevelDebug:
		return "debug"
	case LevelNone:
		return "none"
	default:
		return "info"
	}
}
