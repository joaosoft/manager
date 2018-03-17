package gomanager

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

// HandlerFunc ...
func (web *SimpleWebHttp) AddRoute(method, path string, handler HandlerFunc) error {
	http.HandleFunc(path, handler.(func(http.ResponseWriter, *http.Request)))
	return nil
}

func (web *SimpleWebHttp) Start() error {
	web.started = true
	if err := http.ListenAndServe(web.host, nil); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (web *SimpleWebHttp) Started() bool {
	return web.started
}

func (web *SimpleWebHttp) Handle(writer http.ResponseWriter, request *http.Request) (http.ResponseWriter, *http.Request) {
	return writer, request
}
