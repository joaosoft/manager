package logger

var Instance = NewLoggerEmpty(InfoLevel)

func SetLevel(level Level) {
	Instance.SetLevel(level)
}

func With(prefixes, tags, fields, sufixes map[string]interface{}) ILogger {
	return Instance.With(prefixes, tags, fields, sufixes)
}

func WithPrefixes(prefixes map[string]interface{}) ILogger {
	return Instance.WithPrefixes(prefixes)
}

func WithTags(tags map[string]interface{}) ILogger {
	return Instance.WithTags(tags)
}

func WithFields(fields map[string]interface{}) ILogger {
	return Instance.WithFields(fields)
}

func WithSufixes(sufixes map[string]interface{}) ILogger {
	return Instance.WithSufixes(sufixes)
}

func WithPrefix(key string, value interface{}) ILogger {
	return Instance.WithPrefix(key, value)
}

func WithTag(key string, value interface{}) ILogger {
	return Instance.WithTag(key, value)
}

func WithField(key string, value interface{}) ILogger {
	return Instance.WithField(key, value)
}

func WithSufix(key string, value interface{}) ILogger {
	return Instance.WithSufix(key, value)
}

func Print(message interface{}) IAddition {
	return Instance.Print(message)
}

func Debug(message interface{}) IAddition {
	return Instance.Debug(message)
}

func Info(message interface{}) IAddition {
	return Instance.Info(message)
}

func Warn(message interface{}) IAddition {
	return Instance.Warn(message)
}

func Error(message interface{}) IAddition {
	return Instance.Error(message)
}

func Panic(message interface{}) IAddition {
	return Instance.Panic(message)
}

func Fatal(message interface{}) IAddition {
	return Instance.Fatal(message)
}

func Printf(format string, arguments ...interface{}) IAddition {
	return Instance.Printf(format, arguments)
}

func Debugf(format string, arguments ...interface{}) IAddition {
	return Instance.Debugf(format, arguments)
}

func Infof(format string, arguments ...interface{}) IAddition {
	return Instance.Infof(format, arguments)
}

func Warnf(format string, arguments ...interface{}) IAddition {
	return Instance.Warnf(format, arguments)
}

func Errorf(format string, arguments ...interface{}) IAddition {
	return Instance.Errorf(format, arguments)
}

func Panicf(format string, arguments ...interface{}) IAddition {
	return Instance.Panicf(format, arguments)
}

func Fatalf(format string, arguments ...interface{}) IAddition {
	return Instance.Fatalf(format, arguments)
}

func  IsDebugEnabled() bool {
return Instance.IsDebugEnabled()
}

func  IsInfoEnabled() bool {
	return Instance.IsInfoEnabled()
}

func  IsWarnEnabled() bool {
	return Instance.IsWarnEnabled()
}

func  IsErrorEnabled() bool {
	return Instance.IsErrorEnabled()
}

func  IsPanicEnabled() bool {
	return Instance.IsPanicEnabled()
}

func  IsFatalEnabled() bool {
	return Instance.IsFatalEnabled()
}

func  IsPrintEnabled() bool {
	return Instance.IsPrintEnabled()
}

func Reconfigure(options ...LoggerOption) {
	Instance.Reconfigure(options...)
}
