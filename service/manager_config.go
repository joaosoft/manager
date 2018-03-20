package gomanager

// IConfig ...
type IConfig interface {
	Get() interface{}
	Set(config interface{})
	Save() error
	Reload() error
}

// AddConfig ...
func (manager *GoManager) AddConfig(key string, config IConfig) error {
	manager.configs[key] = config
	log.Infof("config %s added", key)

	return nil
}

// RemoveConfig ...
func (manager *GoManager) RemoveConfig(key string) (IConfig, error) {
	config := manager.configs[key]

	delete(manager.configs, key)
	log.Infof("config %s removed", key)

	return config, nil
}

// GetConfig ...
func (manager *GoManager) GetConfig(key string) IConfig {
	if config, exists := manager.configs[key]; exists {
		return config
	}
	log.Infof("config %s doesn't exist", key)
	return nil
}
