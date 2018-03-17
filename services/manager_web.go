package gomanager

type HandlerFunc interface{}

// IConfig ...
type IWeb interface {
	AddRoute(method, path string, handler HandlerFunc) error
	Start() error
	Started() bool
}

// AddWeb ...
func (manager *GoManager) AddWeb(key string, web IWeb) error {
	manager.webs[key] = web
	log.Infof("web %s added", key)

	return nil
}

// RemoveWeb ...
func (manager *GoManager) RemoveWeb(key string) (IWeb, error) {
	web := manager.webs[key]

	delete(manager.webs, key)
	log.Infof("web %s removed", key)

	return web, nil
}

// GetWeb ...
func (manager *GoManager) GetWeb(key string) IWeb {
	if web, ok := manager.webs[key]; ok {
		return web
	}
	log.Infof("web %s doesn't exist", key)
	return nil
}
