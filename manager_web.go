package manager

type HandlerFunc interface{}
type MiddlewareFunc interface{}

// IConfig ...
type IWeb interface {
	AddRoute(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) error
	Start() error
	Stop() error
	Started() bool
}

// AddWeb ...
func (manager *Manager) AddWeb(key string, web IWeb) error {
	manager.webs[key] = web
	logger.Infof("web %s added", key)

	return nil
}

// RemoveWeb ...
func (manager *Manager) RemoveWeb(key string) (IWeb, error) {
	web := manager.webs[key]

	delete(manager.webs, key)
	logger.Infof("web %s removed", key)

	return web, nil
}

// GetWeb ...
func (manager *Manager) GetWeb(key string) IWeb {
	if web, ok := manager.webs[key]; ok {
		return web
	}
	logger.Infof("web %s doesn't exist", key)
	return nil
}
