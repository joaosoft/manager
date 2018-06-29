package manager

import (
	"net/http"
)

// SimpleWebHttp ...
type SimpleWebHttp struct {
	*http.Server
	Handler *HandlerFunc
	host    string
	started bool
}

// NewSimpleWebHttp...
func NewSimpleWebHttp(host string) IWeb {
	return &SimpleWebHttp{
		Server: &http.Server{Addr: host},
		host:   host,
	}
}

// AddRoute ...
func (web *SimpleWebHttp) AddRoute(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) error {
	http.HandleFunc(path, handler.(func(http.ResponseWriter, *http.Request)))
	return nil
}

// Start ...
func (web *SimpleWebHttp) Start() error {
	if !web.started {
		if err := http.ListenAndServe(web.host, nil); err != nil {
			logger.Error(err)
			return err
		}
	}
	web.started = true

	return nil
}

// Stop ...
func (web *SimpleWebHttp) Stop() error {
	if web.started {
		if err := web.Server.Close(); err != nil {
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
