package manager

import (
	"sync"

	"github.com/joaosoft/web"
)

// SimpleWebServer ...
type SimpleWebServer struct {
	server  *web.Server
	host    string
	started bool
}

// NewSimpleWebServer...
func NewSimpleWebServer(host string) IWeb {
	server, _ := web.NewServer(web.WithServerAddress(host))
	return &SimpleWebServer{
		server: server,
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

// Start ...
func (w *SimpleWebServer) Start(wg *sync.WaitGroup) error {
	if wg == nil {
		wg = &sync.WaitGroup{}
		wg.Add(1)
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
func (w *SimpleWebServer) Stop(wg *sync.WaitGroup) error {
	if wg == nil {
		wg = &sync.WaitGroup{}
		wg.Add(1)
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
