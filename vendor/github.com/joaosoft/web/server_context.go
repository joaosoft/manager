package web

import (
	"time"
)

func NewContext(startTime time.Time, request *Request, response *Response) *Context {
	return &Context{
		StartTime: startTime,
		Request:   request,
		Response:  response,
	}
}
