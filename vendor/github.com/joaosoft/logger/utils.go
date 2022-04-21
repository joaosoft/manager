package logger

import (
	"fmt"
	"strings"
)

func ParseLevel(level string) (Level, error) {
	switch strings.ToLower(level) {
	case "print":
		return PrintLevel, nil
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	case "none":
		return NoneLevel, nil
	default:
		return DefaultLevel, fmt.Errorf("invalid level: %s, set default level: %s", level, DefaultLevel)
	}
}

func (level Level) String() string {
	switch level {
	case PrintLevel:
		return "print"
	case PanicLevel:
		return "panic"
	case FatalLevel:
		return "fatal"
	case ErrorLevel:
		return "error"
	case WarnLevel:
		return "warn"
	case InfoLevel:
		return "info"
	case DebugLevel:
		return "debug"
	case NoneLevel:
		return "none"
	default:
		return "info"
	}
}
