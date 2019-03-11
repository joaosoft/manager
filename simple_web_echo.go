package manager

import (
	"sync"

	"github.com/joaosoft/logger"

	"github.com/labstack/echo"
)

// SimpleWebEcho ...
type SimpleWebEcho struct {
	server  *echo.Echo
	host    string
	logger  logger.ILogger
	started bool
}

// NewSimpleWebEcho...
func (manager *Manager) NewSimpleWebEcho(host string) IWeb {
	e := echo.New()
	e.HideBanner = true

	return &SimpleWebEcho{
		server: e,
		host:   host,
		logger: manager.logger,
	}
}

// AddRoutes ...
func (w *SimpleWebEcho) AddRoutes(routes ...*Route) error {
	for _, route := range routes {
		err := w.AddRoute(route.Method, route.Path, route.Handler, route.Middlewares...)

		if err != nil {
			return err
		}
	}

	return nil
}

// AddRoute ...
func (w *SimpleWebEcho) AddRoute(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) error {
	w.server.Add(method, path, handler.(func(echo.Context) error))
	for _, item := range middleware {
		w.server.Group(path, item.(echo.MiddlewareFunc))
	}
	return nil
}

// Start ...
func (w *SimpleWebEcho) Start(waitGroup ...*sync.WaitGroup) error {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
	}

	defer wg.Done()

	if w.started {
		return nil
	}

	go w.server.Start(w.host)
	w.started = true

	return nil
}

// Stop ...
func (w *SimpleWebEcho) Stop(waitGroup ...*sync.WaitGroup) error {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
	}

	defer wg.Done()

	if !w.started {
		return nil
	}

	if err := w.server.Close(); err != nil {
		return err
	}

	w.started = false

	return nil
}

// Started ...
func (w *SimpleWebEcho) Started() bool {
	return w.started
}

// GetClient ...
func (w *SimpleWebEcho) GetClient() interface{} {
	return w.server
}
