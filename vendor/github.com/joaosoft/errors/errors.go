package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

func (e *Err) Add(newErr *Err) *Err {
	prevErr := &Err{
		Previous: e.Previous,
		Level:    e.Level,
		Code:     e.Code,
		Message:  e.Message,
		Stack:    e.Stack,
	}

	var stack string
	if newErr.Level <= ErrorLevel {
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
		Previous: prevErr,
		Level:    newErr.Level,
		Code:     newErr.Code,
		Message:  newErr.Message,
		Stack:    stack,
	}
}

func (e *Err) Error() string {
	return e.Message
}

func (e *Err) Cause() string {
	str := fmt.Sprintf("'%s'", e.Message)

	prevErr := e.Previous
	for prevErr != nil {
		str += fmt.Sprintf(", caused by '%s'", prevErr.Message)
		prevErr = prevErr.Previous
	}
	return str
}

func (e *Err) Errors() []*Err {
	errors := make([]*Err, 0)
	errors = append(errors, e)

	nextErr := e.Previous
	for nextErr != nil {
		errors = append(errors, e.Previous)
		nextErr = nextErr.Previous
	}

	return errors
}

func (e *Err) Format(values ...interface{}) *Err {
	e.Message = fmt.Sprintf(e.Error(), values...)
	return e
}

func (e *Err) String() string {
	b, _ := json.Marshal(e)
	return string(b)
}
