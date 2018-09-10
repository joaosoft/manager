package manager

import (
	"sync"

	"github.com/joaosoft/webserver"
)

// SimpleWebServer ...
type SimpleWebServer struct {
	server  *webserver.WebServer
	host    string
	started bool
}

// NewSimpleWebServer...
func NewSimpleWebServer(host string) IWeb {
	server, _ := webserver.NewWebServer(webserver.WithAddress(host))
	return &SimpleWebServer{
		server: server,
	}
}

// AddRoutes ...
func (web *SimpleWebServer) AddRoutes(routes ...*Route) error {
	for _, route := range routes {
		err := web.AddRoute(route.Method, route.Path, route.Handler, route.Middlewares...)

		if err != nil {
			return err
		}
	}

	return nil
}

// AddRoute ...
func (web *SimpleWebServer) AddRoute(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) error {
	middlewares := make([]webserver.MiddlewareFunc, 0)
	for _, m := range middleware {
		middlewares = append(middlewares, m.(webserver.MiddlewareFunc))
	}

	return web.server.AddRoute(webserver.Method(method), path, handler.(func(*webserver.Context) error), middlewares...)
}

// Start ...
func (web *SimpleWebServer) Start(wg *sync.WaitGroup) error {
	if wg == nil {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	}

	defer wg.Done()

	if web.started {
		return nil
	}

	web.started = true
	if err := web.server.Start(); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

// Stop ...
func (web *SimpleWebServer) Stop(wg *sync.WaitGroup) error {
	if wg == nil {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	}

	defer wg.Done()

	if !web.started {
		return nil
	}

	web.started = false
	if err := web.server.Stop(); err != nil {
		return err
	}
	return nil
}

// Started ...
func (web *SimpleWebServer) Started() bool {
	return web.started
}

// GetClient ...
func (web *SimpleWebServer) GetClient() interface{} {
	return web.server
}
