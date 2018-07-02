package manager

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

// AddRoutes ...
func (web *SimpleWebEcho) AddRoutes(routes ...Route) error {
	for _, route := range routes {
		err := web.AddRoute(route.Method, route.Path, route.Handler, route.Middlewares)

		if err != nil {
			return err
		}
	}

	return nil
}

// AddRoute ...
func (web *SimpleWebEcho) AddRoute(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) error {
	web.Add(method, path, handler.(func(echo.Context) error))
	for _, item := range middleware {
		web.Group(path, item.(echo.MiddlewareFunc))
	}
	return nil
}

// Start ...
func (web *SimpleWebEcho) Start() error {
	if !web.started {
		if err := web.Echo.Start(web.host); err != nil {
			log.Error(err)
			return err
		}
		web.started = true
	}

	return nil
}

// Stop ...
func (web *SimpleWebEcho) Stop() error {
	if web.started {
		if err := web.Echo.Close(); err != nil {
			return err
		}
		web.started = false
	}
	return nil
}

// Started ...
func (web *SimpleWebEcho) Started() bool {
	return web.started
}
