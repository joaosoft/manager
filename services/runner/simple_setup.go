package runner

import (
	"fmt"

	"github.com/labstack/gommon/log"
	viper "github.com/spf13/viper"
)

type SimpleSetup struct {
	path      string
	file      string
	extension string
	object    interface{}
	config    *viper.Viper
}

// Get...
func (setup *SimpleSetup) Get(key string) interface{} {
	return setup.config.Get(key)
}

// NewSimpleSetup...
func NewSimpleSetup(path string, file string, extension string) (*SimpleSetup, error) {
	config, err := loadConfig(path, file, extension)
	return &SimpleSetup{
		path:      path,
		extension: extension,
		config:    config,
	}, err
}

// unmarshal...
func (setup *SimpleSetup) unmarshal(obj interface{}) error {
	if err := setup.config.Unmarshal(obj); err != nil {
		return err
	}
	setup.object = obj

	return nil
}

// Reload...
func (setup *SimpleSetup) Reload() error {
	var cfg *viper.Viper
	var err error

	if cfg, err = loadConfig(setup.path, setup.file, setup.extension); err != nil {
		return err
	}
	setup.config = cfg
	setup.object = nil

	return nil
}

// loadConfig...
func loadConfig(path string, file string, extension string) (*viper.Viper, error) {
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
