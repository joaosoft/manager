package manager

import (
	"net/http"
	"sync"

	"github.com/joaosoft/logger"
)

// SimpleWebHttp ...
type SimpleWebHttp struct {
	server  *http.Server
	handler *HandlerFunc
	host    string
	logger  logger.ILogger
	started bool
}

// NewSimpleWebHttp...
func (manager *Manager) NewSimpleWebHttp(host string) IWeb {
	return &SimpleWebHttp{
		server: &http.Server{Addr: host},
		host:   host,
		logger: manager.logger,
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

// AddNamespace ...
func (w *SimpleWebHttp) AddNamespace(path string, middleware []MiddlewareFunc, routes ...*Route) error {
	return nil
}

func (w *SimpleWebHttp) AddFilter(pattern string, position string, middleware MiddlewareFunc, method string, methods ...string) {
	// TODO: implement
}

// Start ...
func (w *SimpleWebHttp) Start(waitGroup ...*sync.WaitGroup) error {
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

	go http.ListenAndServe(w.host, nil)

	w.started = true

	return nil
}

// Stop ...
func (w *SimpleWebHttp) Stop(waitGroup ...*sync.WaitGroup) error {
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
func (w *SimpleWebHttp) Started() bool {
	return w.started
}

// GetClient ...
func (w *SimpleWebHttp) GetClient() interface{} {
	return w.server
}
