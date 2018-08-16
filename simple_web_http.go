package manager

import (
	"net/http"
	"sync"
)

// SimpleWebHttp ...
type SimpleWebHttp struct {
	server  *http.Server
	handler *HandlerFunc
	host    string
	started bool
}

// NewSimpleWebHttp...
func NewSimpleWebHttp(host string) IWeb {
	return &SimpleWebHttp{
		server: &http.Server{Addr: host},
		host:   host,
	}
}

// AddRoutes ...
func (web *SimpleWebHttp) AddRoutes(routes ...*Route) error {
	for _, route := range routes {
		err := web.AddRoute(route.Method, route.Path, route.Handler, route.Middlewares...)

		if err != nil {
			return err
		}
	}

	return nil
}

// AddRoute ...
func (web *SimpleWebHttp) AddRoute(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) error {
	http.HandleFunc(path, handler.(func(http.ResponseWriter, *http.Request)))
	return nil
}

// Start ...
func (web *SimpleWebHttp) Start(wg *sync.WaitGroup) error {
	wg.Done()

	if !web.started {
		if err := http.ListenAndServe(web.host, nil); err != nil {
			log.Error(err)
			return err
		}
	}
	web.started = true

	return nil
}

// Stop ...
func (web *SimpleWebHttp) Stop(wg *sync.WaitGroup) error {
	defer wg.Done()

	if web.started {
		if err := web.server.Close(); err != nil {
			return err
		}
		web.started = false
	}
	return nil
}

// Started ...
func (web *SimpleWebHttp) Started() bool {
	return web.started
}

// GetClient ...
func (web *SimpleWebHttp) GetClient() interface{} {
	return web.server
}
