package mgr

import (
	"fmt"

	"github.com/joaosoft/go-manager/services/config"
	"github.com/labstack/gommon/log"
)

// -------------- CONFIGURATION CLIENTS --------------
// NewJSONFile ... creates a new nsq producer
func (instance *Manager) NewSimpleConfig(path string, file string, extension string) (config.IConfig, error) {
	return config.NewSimpleConfig(path, file, extension)
}

// -------------- METHODS --------------
// GetConfig ... get a config with key
func (instance *Manager) GetConfig(key string) config.IConfig {
	return instance.ConfigController[key]
}

// Unmarshal ... unmarshal configuration
func (instance *Manager) Unmarshal(key string, obj interface{}) (interface{}, error) {
	if err := instance.ConfigController[key].Unmarshal(obj); err != nil {
		return nil, err
	}
	return obj, nil
}

// AddProcess ... add a config with key
func (instance *Manager) AddConfig(key string, cfg config.IConfig) error {
	if instance.Started {
		panic("Manager, can not add config after start")
	}

	instance.ConfigController[key] = &config.ConfigController{
		Path:   "",
		Config: cfg}

	log.Infof(fmt.Sprintf("Manager, config '%s' added", key))

	return nil
}

// RemConfig ... remove the config by bey
func (instance *Manager) RemConfig(key string) (config.IConfig, error) {
	// get config
	controller := instance.ConfigController[key]

	// delete config
	delete(instance.ConfigController, key)
	log.Infof(fmt.Sprintf("Manager, config '%s' removed", key))

	return controller, nil
}
