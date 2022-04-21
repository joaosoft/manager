package web

import (
	"fmt"
)

var (
	ErrorNotFound     = NewError(StatusNotFound, "route not found")
	ErrorInvalidChunk = NewError(StatusNotFound, "invalid chunk length")
)

type Error struct {
	Status  Status `json:"status"`
	Message string `json:"message"`
}

func NewError(status Status, messages ...string) *Error {
	err := &Error{
		Status: status,
	}

	if len(messages) > 0 {
		err.Message = messages[0]
	} else {
		err.Message = StatusText(status)
	}

	return err
}

func (e *Error) Error() string {
	return fmt.Sprintf("status=%d, message=%v", e.Status, e.Message)
}

func (w *Server) DefaultErrorHandler(ctx *Context, err error) error {
	w.logger.Infof("handling error: %s", err)

	if e, ok := err.(*Error); ok {
		return ctx.Response.JSON(e.Status, e)
	}

	return ctx.Response.JSON(StatusInternalServerError, NewError(StatusInternalServerError, err.Error()))
}
