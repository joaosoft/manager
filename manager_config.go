package manager

import "time"

// IConfig ...
type IConfig interface {
	Get(key string) interface{}
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetInt64(key string) int64
	GetFloat64(key string) float64
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration
	GetStringSlice(key string) []string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string

	GetObj() interface{}
	Set(config interface{})
	Save() error
	Reload() error
}

// AddConfig ...
func (manager *Manager) AddConfig(key string, config IConfig) error {
	manager.configs[key] = config
	manager.logger.Infof("config %s added", key)

	return nil
}

// RemoveConfig ...
func (manager *Manager) RemoveConfig(key string) (IConfig, error) {
	config := manager.configs[key]

	delete(manager.configs, key)
	manager.logger.Infof("config %s removed", key)

	return config, nil
}

// GetConfig ...
func (manager *Manager) GetConfig(key string) IConfig {
	if config, exists := manager.configs[key]; exists {
		return config
	}
	manager.logger.Infof("config %s doesn't exist", key)
	return nil
}
