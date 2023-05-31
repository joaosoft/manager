# logger
[![Build Status](https://travis-ci.org/joaosoft/logger.svg?branch=master)](https://travis-ci.org/joaosoft/logger) | [![codecov](https://codecov.io/gh/joaosoft/logger/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/logger) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/logger)](https://goreportcard.com/report/github.com/joaosoft/logger) | [![GoDoc](https://godoc.org/github.com/joaosoft/logger?status.svg)](https://godoc.org/github.com/joaosoft/logger)

A simplified logger that allows you to add complexity depending of your requirements.
The easy way to use the logger:
``` Go
import log github.com/joaosoft/logger

log.Info("hello")
```
you also can config it, as i prefer, please see below
After a read of the project https://gitlab.com/vredens/loggerger extracted some concepts like allowing to add tags and fields to logger infrastructure. 

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for
* formatted messages
* prefixes (special prefixes: **DATE, TIME, TIMESTAMP, LEVEL, IP, PACKAGE, FUNCTION, FILE, TRACE, STACK**)
* tags
* fields
* writers at [[writer]](https://github.com/joaosoft/writers/tree/master/bin/examples)
  * to file (with queue processing)[1] 
  * to stdout (with queue processing)[1] [[here]](https://github.com/joaosoft/writers/tree/master/examples)
* addition commands (ToError())
  
  **[1]** this writer allows you to continue the processing and dispatch the logging

## Levels
* **DefaultLevel**, default level
* **PanicLevel**, when there is no recover
* **FatalLevel**, when the error is fatal to the application
* **ErrorLevel**, when there is a controlled error
* **WarnLevel**, when there is a warning
* **InfoLevel**, when it is a informational message
* **DebugLevel**, when it is a debugging message
* **PrintLevel**, when it is a system message
* **NoneLevel**, when the logging is disabled

## Special Prefix's
* **LEVEL**, add the level value to the prefix
* **TIMESTAMP**, add the timestamp value to the prefix
* **DATE**, add the date value to the prefix
* **TIME**, add the time value to the prefix
* **IP**, add the client ip address
* **TRACE**, add the error trace
* **PACKAGE**, add the package name
* **FILE**, add the file
* **FUNCTION**, add the function name
* **STACK**, add the debug stack


## Dependecy Management 
>### Dep

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Install dependencies: `dep ensure`
* Update dependencies: `dep ensure -update`


>### Go
```
go get github.com/joaosoft/logger/service
```

## Interface 
```go
type Logger interface {
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

type IAddition interface {
	ToError() error
}

type ISpecialWriter interface {
	SWrite(prefixes map[string]interface{}, tags map[string]interface{}, message interface{}, fields map[string]interface{}, sufixes map[string]interface{}) (n int, err error)
}
```

## Usage 
This examples are available in the project at [logger/examples](https://github.com/joaosoft/logger/tree/master/examples)

```go
//
// log to text
fmt.Println(":: LOG TEXT")
log := logger.NewLogger(
    logger.WithLevel(logger.InfoLevel),
    logger.WithFormatHandler(writer.TextFormatHandler),
    logger.WithWriter(os.Stdout)).
    With(
        map[string]interface{}{"level": logger.LEVEL, "timestamp": logger.TIMESTAMP, "date": logger.DATE, "time": logger.TIME},
        map[string]interface{}{"service": "log"},
        map[string]interface{}{"name": "joão"},
        map[string]interface{}{"ip": logger.IP, "function": logger.FUNCTION, "file": logger.FILE})

// logging...
log.Error("isto é uma mensagem de error")
log.Info("isto é uma mensagem de info")
log.Debug("isto é uma mensagem de debug")
log.Error("")

fmt.Println("--------------")
<-time.After(time.Second)

//
// log to json
fmt.Println(":: LOG JSON")
log = logger.NewLogger(
    logger.WithLevel(logger.InfoLevel),
    logger.WithFormatHandler(writer.JsonFormatHandler),
    logger.WithWriter(os.Stdout)).
    With(
        map[string]interface{}{"level": logger.LEVEL, "timestamp": logger.TIMESTAMP, "date": logger.DATE, "time": logger.TIME},
        map[string]interface{}{"service": "log"},
        map[string]interface{}{"name": "joão"},
        map[string]interface{}{"ip": logger.IP, "function": logger.FUNCTION, "file": logger.FILE})

// logging...
log.Errorf("isto é uma mensagem de error %s", "hello")
log.Infof("isto é uma  mensagem de info %s ", "hi")
log.Debugf("isto é uma mensagem de debug %s", "ehh")
```

###### Output 

```javascript
default...
:: LOG TEXT
{prefixes:map[level:error timestamp:2018-08-16 20:27:13:18 date:2018-08-16 time:20:27:13:18] tags:map[service:log] message:isto é uma mensagem de error fields:map[name:joão] sufixes:map[ip:192.168.1.4 function:Example.ExampleDefaultLogger file:/Users/joaoribeiro/workspace/go/personal/src/logger/examples/main.go]}
{prefixes:map[level:info timestamp:2018-08-16 20:27:13:18 date:2018-08-16 time:20:27:13:18] tags:map[service:log] message:isto é uma mensagem de info fields:map[name:joão] sufixes:map[ip:192.168.1.4]}
{prefixes:map[level:error timestamp:2018-08-16 20:27:13:18 date:2018-08-16 time:20:27:13:18] tags:map[service:log] message: fields:map[name:joão] sufixes:map[ip:192.168.1.4 function:Example.ExampleDefaultLogger file:/Users/joaoribeiro/workspace/go/personal/src/logger/examples/main.go]}
--------------
:: LOG JSON
{"prefixes":{"date":"2018-08-16","level":"error","time":"20:27:14:18","timestamp":"2018-08-16 20:27:14:18"},"tags":{"service":"log"},"message":"isto é uma mensagem de error hello","fields":{"name":"joão"},"sufixes":{"file":"/Users/joaoribeiro/workspace/go/personal/src/logger/examples/main.go","function":"Example.ExampleDefaultLogger","ip":"192.168.1.4"}}
{"prefixes":{"date":"2018-08-16","level":"info","time":"20:27:14:18","timestamp":"2018-08-16 20:27:14:18"},"tags":{"service":"log"},"message":"isto é uma  mensagem de info hi ","fields":{"name":"joão"},"sufixes":{"ip":"192.168.1.4"}}
```

## Known issues
* all the maps do not guarantee order of the items! 


## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
