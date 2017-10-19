package config

import (
	"fmt"
	"github.com/labstack/gommon/log"
	viper "github.com/spf13/viper"
)

type simpleConfig struct {
	path        string
	file        string
	extension   string
	object      interface{}
	viperConfig *viper.Viper
}

// Get ... get a configuration by key
func (instance *simpleConfig) Get(key string) interface{} {
	return instance.viperConfig.Get(key)
}

func NewSimpleConfig(path string, file string, extension string) (*simpleConfig, error) {
	config, err := LoadConfig(path, file, extension)
	return &simpleConfig{
		path:        path,
		extension:   extension,
		viperConfig: config,
	}, err
}

// Unmarshal ... unmarshal configuration
func (instance *simpleConfig) Unmarshal(obj interface{}) error {
	if err := instance.viperConfig.Unmarshal(obj); err != nil {
		return err
	}
	instance.object = obj

	return nil
}

// Reload ... reload the configuration file
func (instance *simpleConfig) Reload(key string) error {
	var cfg *viper.Viper
	var err error

	if cfg, err = LoadConfig(instance.path, instance.file, instance.extension); err != nil {
		return err
	}
	instance.viperConfig = cfg
	instance.object = nil

	return nil
}

// LoadConfig ... loads the configuration file
func LoadConfig(path string, file string, extension string) (*viper.Viper, error) {
	viperConfig := viper.New()

	viperConfig.SetConfigName(file)
	viperConfig.SetConfigType(extension)
	viperConfig.AddConfigPath(path)

	if err := viperConfig.ReadInConfig(); err != nil {
		log.Error(err)
		return nil, err
	}
	fmt.Println(viperConfig)

	return viperConfig, nil
}
