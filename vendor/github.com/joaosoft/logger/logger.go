package logger

import (
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	writer "github.com/joaosoft/writers"
)

var logger = NewLoggerEmpty(InfoLevel)

// NewLogger ...
func NewLogger(options ...LoggerOption) ILogger {
	logger := &Logger{
		writer:        []io.Writer{os.Stdout},
		formatHandler: writer.JsonFormatHandler,
		level:         InfoLevel,
		prefixes:      make(map[string]interface{}),
		tags:          make(map[string]interface{}),
		fields:        make(map[string]interface{}),
	}
	logger.Reconfigure(options...)

	return logger
}

// NewLogDefault
func NewLogDefault(service string, level Level) ILogger {
	return NewLogger(
		WithLevel(level),
		WithFormatHandler(writer.JsonFormatHandler),
		WithWriter(os.Stdout)).
		With(
			map[string]interface{}{"level": LEVEL, "timestamp": TIMESTAMP},
			map[string]interface{}{"service": service},
			map[string]interface{}{},
			map[string]interface{}{"ip": IP, "package": PACKAGE, "function": FUNCTION, "stack": STACK, "trace": TRACE})
}

// NewLoggerEmpty
func NewLoggerEmpty(level Level) ILogger {
	return NewLogger(
		WithLevel(level),
		WithFormatHandler(writer.JsonFormatHandler),
		WithWriter(os.Stdout)).
		WithPrefixes(map[string]interface{}{"level": LEVEL, "timestamp": TIMESTAMP}).
		WithSufixes(map[string]interface{}{"ip": IP, "package": PACKAGE, "function": FUNCTION, "stack": STACK, "trace": TRACE})
}

func (logger *Logger) SetLevel(level Level) {
	logger.level = level
}

func (logger *Logger) With(prefixes, tags, fields, sufixes map[string]interface{}) ILogger {
	newLog := logger.clone().WithPrefixes(prefixes).WithTags(tags).WithFields(fields).WithSufixes(sufixes)
	return newLog
}

func (logger *Logger) WithPrefixes(prefixes map[string]interface{}) ILogger {
	newLog := logger.clone()
	newLog.prefixes = prefixes
	return newLog
}

func (logger *Logger) WithTags(tags map[string]interface{}) ILogger {
	newLog := logger.clone()
	newLog.tags = tags
	return newLog
}

func (logger *Logger) WithFields(fields map[string]interface{}) ILogger {
	newLog := logger.clone()
	newLog.fields = fields
	return newLog
}

func (logger *Logger) WithSufixes(sufixes map[string]interface{}) ILogger {
	newLog := logger.clone()
	newLog.sufixes = sufixes
	return newLog
}

func (logger *Logger) WithPrefix(key string, value interface{}) ILogger {
	newLog := logger.clone()
	newLog.prefixes[key] = fmt.Sprintf("%s", value)
	return newLog
}

func (logger *Logger) WithTag(key string, value interface{}) ILogger {
	newLog := logger.clone()
	newLog.tags[key] = fmt.Sprintf("%s", value)
	return newLog
}

func (logger *Logger) WithField(key string, value interface{}) ILogger {
	newLog := logger.clone()
	newLog.fields[key] = fmt.Sprintf("%s", value)
	return newLog
}

func (logger *Logger) WithSufix(key string, value interface{}) ILogger {
	newLog := logger.clone()
	newLog.sufixes[key] = fmt.Sprintf("%s", value)
	return newLog
}

// CopyDependency ...
func (logger *Logger) clone() *Logger {
	return &Logger{
		level:         logger.level,
		writer:        logger.writer,
		formatHandler: logger.formatHandler,
		specialWriter: logger.specialWriter,
		tags:          logger.tags,
		prefixes:      logger.prefixes,
		fields:        logger.fields,
		sufixes:       logger.sufixes,
	}
}

func (logger *Logger) Print(message interface{}) IAddition {
	msg := fmt.Sprint(message)
	logger.writeLog(PrintLevel, message)

	return NewAddition(msg)
}

func (logger *Logger) Debug(message interface{}) IAddition {
	msg := fmt.Sprint(message)
	logger.writeLog(DebugLevel, message)

	return NewAddition(msg)
}

func (logger *Logger) Info(message interface{}) IAddition {
	msg := fmt.Sprint(message)
	logger.writeLog(InfoLevel, msg)

	return NewAddition(msg)
}

func (logger *Logger) Warn(message interface{}) IAddition {
	msg := fmt.Sprint(message)
	logger.writeLog(WarnLevel, msg)

	return NewAddition(msg)
}

func (logger *Logger) Error(message interface{}) IAddition {
	msg := fmt.Sprint(message)
	logger.writeLog(ErrorLevel, msg)

	return NewAddition(msg)
}

func (logger *Logger) Panic(message interface{}) IAddition {
	msg := fmt.Sprint(message)
	logger.writeLog(PanicLevel, msg)

	return NewAddition(msg)
}

func (logger *Logger) Fatal(message interface{}) IAddition {
	msg := fmt.Sprint(message)
	logger.writeLog(FatalLevel, msg)

	return NewAddition(msg)
}

func (logger *Logger) Printf(format string, arguments ...interface{}) IAddition {
	msg := fmt.Sprintf(format, arguments...)
	logger.writeLog(PrintLevel, msg)

	return NewAddition(msg)
}

func (logger *Logger) Debugf(format string, arguments ...interface{}) IAddition {
	msg := fmt.Sprintf(format, arguments...)
	logger.writeLog(DebugLevel, msg)

	return NewAddition(msg)
}

func (logger *Logger) Infof(format string, arguments ...interface{}) IAddition {
	msg := fmt.Sprintf(format, arguments...)
	logger.writeLog(InfoLevel, msg)

	return NewAddition(msg)
}

func (logger *Logger) Warnf(format string, arguments ...interface{}) IAddition {
	msg := fmt.Sprintf(format, arguments...)
	logger.writeLog(WarnLevel, msg)

	return NewAddition(msg)
}

func (logger *Logger) Errorf(format string, arguments ...interface{}) IAddition {
	msg := fmt.Sprintf(format, arguments...)
	logger.writeLog(ErrorLevel, msg)

	return NewAddition(msg)
}

func (logger *Logger) Fatalf(format string, arguments ...interface{}) IAddition {
	msg := fmt.Sprintf(format, arguments...)
	logger.writeLog(FatalLevel, msg)

	return NewAddition(msg)
}

func (logger *Logger) Panicf(format string, arguments ...interface{}) IAddition {
	msg := fmt.Sprintf(format, arguments...)
	logger.writeLog(PanicLevel, msg)

	return NewAddition(msg)
}

func (logger *Logger) IsDebugEnabled() bool {
	return logger.level == DebugLevel
}

func (logger *Logger) IsInfoEnabled() bool {
	return logger.level == InfoLevel
}

func (logger *Logger) IsWarnEnabled() bool {
	return logger.level == WarnLevel
}

func (logger *Logger) IsErrorEnabled() bool {
	return logger.level == ErrorLevel
}

func (logger *Logger) IsPanicEnabled() bool {
	return logger.level == PanicLevel
}

func (logger *Logger) IsFatalEnabled() bool {
	return logger.level == FatalLevel
}

func (logger *Logger) IsPrintEnabled() bool {
	return logger.level == PrintLevel
}

func (logger *Logger) writeLog(level Level, message interface{}) {
	if level > logger.level {
		return
	}

	prefixes := handleSpecialTags(level, logger.prefixes)
	sufixes := handleSpecialTags(level, logger.sufixes)
	if logger.specialWriter == nil {
		if bytes, err := logger.formatHandler(prefixes, logger.tags, message, logger.fields, sufixes); err != nil {
			return
		} else {
			for _, w := range logger.writer {
				w.Write(bytes)
			}
		}
	} else {
		for _, w := range logger.specialWriter {
			w.SWrite(prefixes, logger.tags, message, logger.fields, sufixes)
		}
	}
}

func handleSpecialTags(level Level, prefixes map[string]interface{}) map[string]interface{} {
	newPrefixes := make(map[string]interface{})
	for key, value := range prefixes {
		switch value {
		case LEVEL:
			value = level.String()

		case TIMESTAMP:
			value = time.Now().Format("2006-01-02 15:04:05:06")

		case DATE:
			value = time.Now().Format("2006-01-02")

		case TIME:
			value = time.Now().Format("15:04:05:06")

		case IP:
			addresses, _ := net.InterfaceAddrs()
			for _, a := range addresses {
				if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						value = ipNet.IP.String()
					}
				}
			}

		case TRACE:
			if level <= ErrorLevel {
				pc := make([]uintptr, 1)
				runtime.Callers(4, pc)
				function := runtime.FuncForPC(pc[0])
				file, line := function.FileLine(pc[0])
				info := strings.SplitN(function.Name(), ".", 2)
				stack := string(debug.Stack())
				stack = stack[strings.Index(stack, function.Name()):]

				value = struct {
					File     string `json:"file"`
					Line     int    `json:"line"`
					Package  string `json:"package"`
					Function string `json:"function"`
					Stack    string `json:"stack"`
				}{
					File:     file,
					Line:     line,
					Package:  info[0],
					Function: info[1],
					Stack:    stack,
				}
			} else {
				continue
			}

		case FILE:
			if level <= ErrorLevel {
				pc := make([]uintptr, 1)
				runtime.Callers(4, pc)
				function := runtime.FuncForPC(pc[0])
				value, _ = function.FileLine(pc[0])
			} else {
				continue
			}

		case PACKAGE:
			if level <= ErrorLevel {
				pc := make([]uintptr, 1)
				runtime.Callers(4, pc)
				function := runtime.FuncForPC(pc[0])
				value = strings.SplitN(function.Name(), ".", 2)[0]
			} else {
				continue
			}

		case FUNCTION:
			if level <= ErrorLevel {
				pc := make([]uintptr, 1)
				runtime.Callers(4, pc)
				function := runtime.FuncForPC(pc[0])
				value = strings.SplitN(function.Name(), ".", 2)[1]
			} else {
				continue
			}

		case STACK:
			if level <= ErrorLevel {
				pc := make([]uintptr, 1)
				runtime.Callers(4, pc)
				function := runtime.FuncForPC(pc[0])
				stack := string(debug.Stack())
				value = stack[strings.Index(stack, function.Name()):]
			} else {
				continue
			}
		}

		newPrefixes[key] = value
	}
	return newPrefixes
}
