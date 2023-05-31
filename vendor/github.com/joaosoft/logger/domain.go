package logger

import (
	"io"

	"github.com/joaosoft/writers"
)

type IAddition interface {
	ToError() error
}

type ISpecialWriter interface {
	SWrite(prefixes map[string]interface{}, tags map[string]interface{}, message interface{}, fields map[string]interface{}, sufixes map[string]interface{}) (n int, err error)
}

type ILogger interface {
	SetLevel(level Level)

	With(prefixes, tags, fields, sufixes map[string]interface{}) ILogger
	WithPrefixes(prefixes map[string]interface{}) ILogger
	WithTags(tags map[string]interface{}) ILogger
	WithFields(fields map[string]interface{}) ILogger
	WithSufixes(sufixes map[string]interface{}) ILogger

	WithPrefix(key string, value interface{}) ILogger
	WithTag(key string, value interface{}) ILogger
	WithField(key string, value interface{}) ILogger
	WithSufix(key string, value interface{}) ILogger

	Debug(message interface{}) IAddition
	Info(message interface{}) IAddition
	Warn(message interface{}) IAddition
	Error(message interface{}) IAddition
	Panic(message interface{}) IAddition
	Fatal(message interface{}) IAddition
	Print(message interface{}) IAddition

	Debugf(format string, arguments ...interface{}) IAddition
	Infof(format string, arguments ...interface{}) IAddition
	Warnf(format string, arguments ...interface{}) IAddition
	Errorf(format string, arguments ...interface{}) IAddition
	Panicf(format string, arguments ...interface{}) IAddition
	Fatalf(format string, arguments ...interface{}) IAddition
	Printf(format string, arguments ...interface{}) IAddition

	IsDebugEnabled() bool
	IsInfoEnabled() bool
	IsWarnEnabled() bool
	IsErrorEnabled() bool
	IsPanicEnabled() bool
	IsFatalEnabled() bool
	IsPrintEnabled() bool

	Reconfigure(options ...LoggerOption)
}

// Logger ...
type Logger struct {
	level         Level
	writer        []io.Writer
	specialWriter []ISpecialWriter
	prefixes      map[string]interface{} `json:"prefixes"`
	tags          map[string]interface{} `json:"tags"`
	fields        map[string]interface{} `json:"fields"`
	sufixes       map[string]interface{} `json:"sufixes"`
	formatHandler writers.FormatHandler
}
