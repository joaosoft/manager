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
func (w *SimpleWebHttp) AddRoutes(routes ...*Route) error {
	for _, route := range routes {
		err := w.AddRoute(route.Method, route.Path, route.Handler, route.Middlewares...)

		if err != nil {
			return err
		}
	}

	return nil
}

// AddRoute ...
func (w *SimpleWebHttp) AddRoute(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) error {
	http.HandleFunc(path, handler.(func(http.ResponseWriter, *http.Request)))
	return nil
}

// Start ...
func (w *SimpleWebHttp) Start(wg *sync.WaitGroup) error {
	if wg == nil {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	}

	defer wg.Done()

	if w.started {
		return nil
	}

	go http.ListenAndServe(w.host, nil)

	w.started = true

	return nil
}

// Stop ...
func (w *SimpleWebHttp) Stop(wg *sync.WaitGroup) error {
	if wg == nil {
		wg = &sync.WaitGroup{}
		wg.Add(1)
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
func (w *SimpleWebHttp) Started() bool {
	return w.started
}

// GetClient ...
func (w *SimpleWebHttp) GetClient() interface{} {
	return w.server
}
