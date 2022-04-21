package errors

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

func New(level Level, code int, err interface{}, params ...interface{}) *Err {

	var message string
	switch v := err.(type) {
	case error:
		message = v.Error()

	case string:
		message = v

	default:
		message = fmt.Sprint(v)
	}

	if len(params) > 0 {
		message = fmt.Sprintf(fmt.Sprint(message), params...)
	}

	var stack string
	if level <= ErrorLevel {
		pc := make([]uintptr, 1)
		runtime.Callers(2, pc)
		function := runtime.FuncForPC(pc[0])
		stack = string(debug.Stack())
		index := strings.Index(stack, function.Name())
		if index < 0 {
			index = 0
		}
		stack = stack[index:]
	}

	return &Err{
		Level:   level,
		Code:    code,
		Message: message,
		Stack:   stack,
	}
}

func Add(err *Err) *Err {

	var stack string
	if err.Level <= ErrorLevel {
		pc := make([]uintptr, 1)
		runtime.Callers(2, pc)
		function := runtime.FuncForPC(pc[0])
		stack = string(debug.Stack())
		index := strings.Index(stack, function.Name())
		if index < 0 {
			index = 0
		}
		stack = stack[index:]
	}

	return &Err{
		Level:   err.Level,
		Code:    err.Code,
		Message: err.Message,
		Stack:   stack,
	}
}
