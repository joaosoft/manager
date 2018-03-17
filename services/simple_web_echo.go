package gomanager

import (
	"github.com/labstack/echo"
)

// SimpleWebEcho ...
type SimpleWebEcho struct {
	*echo.Echo
	host    string
	started bool
}

// NewSimpleWebEcho...
func NewSimpleWebEcho(host string) IWeb {
	e := echo.New()
	e.HideBanner = true

	return &SimpleWebEcho{
		Echo: e,
		host: host,
	}
}

// HandlerFunc ...
func (web *SimpleWebEcho) AddRoute(method, path string, handler HandlerFunc) error {
	web.Add(method, path, handler.(func(echo.Context) error))
	return nil
}

func (web *SimpleWebEcho) Start() error {
	web.started = true
	if err := web.Echo.Start(web.host); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (web *SimpleWebEcho) Started() bool {
	return web.started
}

type HandleFuncEcho func(ctx echo.Context)
