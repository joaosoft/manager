package manager

import "sync"

type HandlerFunc interface{}
type MiddlewareFunc interface{}
type Route struct {
	Method      string
	Path        string
	Handler     HandlerFunc
	Middlewares []MiddlewareFunc
}

func NewRoute(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route {
	return &Route{Method: method, Path: path, Handler: handler, Middlewares: middleware}
}

// IConfig ...
type IWeb interface {
	AddRoute(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) error
	AddRoutes(routes ...*Route) error
	AddNamespace(path string, middleware []MiddlewareFunc, routes ...*Route) error
	AddFilter(pattern string, position string, middleware MiddlewareFunc, method string, methods ...string)
	Start(waitGroup ...*sync.WaitGroup) error
	Stop(waitGroup ...*sync.WaitGroup) error
	Started() bool
	GetClient() interface{}
}

// AddWeb ...
func (manager *Manager) AddWeb(key string, web IWeb) error {
	manager.webs[key] = web
	manager.logger.Infof("web %s added", key)

	return nil
}

// RemoveWeb ...
func (manager *Manager) RemoveWeb(key string) (IWeb, error) {
	web := manager.webs[key]

	delete(manager.webs, key)
	manager.logger.Infof("web %s removed", key)

	return web, nil
}

// GetWeb ...
func (manager *Manager) GetWeb(key string) IWeb {
	if web, ok := manager.webs[key]; ok {
		return web
	}
	manager.logger.Infof("web %s doesn't exist", key)
	return nil
}
