package manager

import (
	"sync"

	"web/common"

	"github.com/joaosoft/web/server"
)

// SimpleWebServer ...
type SimpleWebServer struct {
	server  *server.Server
	host    string
	started bool
}

// NewSimpleWebServer...
func NewSimpleWebServer(host string) IWeb {
	server, _ := server.NewServer(server.WithAddress(host))
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
func (web *SimpleWebServer) AddRoute(method string, path string, handler HandlerFunc, middleware ...MiddlewareFunc) error {
	middlewares := make([]server.MiddlewareFunc, 0)
	for _, m := range middleware {
		middlewares = append(middlewares, m.(server.MiddlewareFunc))
	}

	return web.server.AddRoute(common.Method(method), path, handler.(func(*server.Context) error), middlewares...)
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
	go web.server.Start()

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
