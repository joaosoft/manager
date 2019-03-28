package manager

import (
	"sync"

	"github.com/joaosoft/logger"

	"github.com/joaosoft/web"
)

// SimpleWebServer ...
type SimpleWebServer struct {
	server  *web.Server
	host    string
	logger  logger.ILogger
	started bool
}

// NewSimpleWebServer...
func (manager *Manager) NewSimpleWebServer(host string) IWeb {
	server, _ := web.NewServer(web.WithServerAddress(host))
	return &SimpleWebServer{
		server: server,
		logger: manager.logger,
	}
}

// AddRoutes ...
func (w *SimpleWebServer) AddRoutes(routes ...*Route) error {
	for _, route := range routes {
		err := w.AddRoute(route.Method, route.Path, route.Handler, route.Middlewares...)

		if err != nil {
			return err
		}
	}

	return nil
}

// AddRoute ...
func (w *SimpleWebServer) AddRoute(method string, path string, handler HandlerFunc, middleware ...MiddlewareFunc) error {
	middlewares := make([]web.MiddlewareFunc, 0)
	for _, m := range middleware {
		middlewares = append(middlewares, m.(web.MiddlewareFunc))
	}

	return w.server.AddRoute(web.Method(method), path, handler.(func(*web.Context) error), middlewares...)
}

// AddNamespace ...
func (w *SimpleWebServer) AddNamespace(path string, middleware []MiddlewareFunc, routes ...*Route) error {

	middlewares := make([]web.MiddlewareFunc, 0)
	for _, m := range middleware {
		middlewares = append(middlewares, m.(web.MiddlewareFunc))
	}

	namespace := w.server.AddNamespace(path, middlewares...)

	for _, route := range routes {
		err := w.AddRoute(route.Method, route.Path, route.Handler, route.Middlewares...)

		if err != nil {
			return err
		}

		middlewares := make([]web.MiddlewareFunc, 0)
		for _, m := range route.Middlewares {
			middlewares = append(middlewares, m.(web.MiddlewareFunc))
		}

		if err = namespace.AddRoute(web.Method(route.Method), route.Path, route.Handler.(func(*web.Context) error), middlewares...); err != nil {
			return err
		}
	}

	return nil
}

func (w *SimpleWebServer) AddFilter(pattern string, position string, middleware MiddlewareFunc, method string, methods ...string) {
	webMethods := make([]web.Method, 0)
	for _, m := range methods {
		webMethods = append(webMethods, web.Method(m))
	}

	w.server.AddFilter(pattern, web.Position(position), middleware.(web.MiddlewareFunc), web.Method(method), webMethods...)
}

// Start ...
func (w *SimpleWebServer) Start(waitGroup ...*sync.WaitGroup) error {
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

	go w.server.Start()
	w.started = true

	return nil
}

// Stop ...
func (w *SimpleWebServer) Stop(waitGroup ...*sync.WaitGroup) error {
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

	if err := w.server.Stop(); err != nil {
		return err
	}

	w.started = false

	return nil
}

// Started ...
func (w *SimpleWebServer) Started() bool {
	return w.started
}

// GetClient ...
func (w *SimpleWebServer) GetClient() interface{} {
	return w.server
}
