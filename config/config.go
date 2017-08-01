package config

// Config ... config interface
type IConfig interface {
	Get(key string) interface{}
	Unmarshal(obj interface{}) error
}

// ConfigController ... config structure
type ConfigController struct {
	Path   string
	Config IConfig
}

// NewConfig ... create a new ConfigController
func NewConfig(path string, config IConfig) IConfig {

	return &ConfigController{
		Path:   path,
		Config: config,
	}
}

// Get ... get a configuration by key
func (instance *ConfigController) Get(key string) interface{} {
	return instance.Config.Get(key)
}

// Reload ... reload the configuration file
func (instance *ConfigController) Reload(key string) error {
	return nil
}

// Unmarshal ... unmarshal configuration
func (instance *ConfigController) Unmarshal(obj interface{}) error {
	if err := instance.Config.Unmarshal(obj); err != nil {
		return err
	}
	return nil
}
