package web

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/joaosoft/logger"
)

func (w *Server) handlerFile(ctx *Context) error {
	logger.Infof("handling file %s", ctx.Request.Address.Full)

	dir, _ := os.Getwd()
	path := fmt.Sprintf("%s%s", dir, ctx.Request.Address.Full)

	if _, err := os.Stat(path); err == nil {
		if bytes, err := ioutil.ReadFile(path); err != nil {
			ctx.Response.Status = StatusNotFound
		} else {
			ctx.Response.Status = StatusOK
			ctx.Response.Body = bytes
		}
	} else {
		ctx.Response.Status = StatusNotFound
	}

	return nil
}
