package errors

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

func New(level Level, code interface{}, err interface{}, params ...interface{}) *Error {

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
	if level <= LevelError {
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

	return &Error{
		Level:   level,
		Code:    code,
		Message: message,
		Stack:   stack,
	}
}

func Add(err *Error) *Error {

	var stack string
	if err.Level <= LevelError {
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

	return &Error{
		Level:   err.Level,
		Code:    err.Code,
		Message: err.Message,
		Stack:   stack,
	}
}

func AddList(errs ...*Error) *ErrorList {
	errorList := &ErrorList{}

	for _, err := range errs {
		errorList.Add(err)
	}

	return errorList
}
