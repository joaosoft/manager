package manager

import "fmt"

// AppConfig ...
type AppConfig struct {
	Manager *ManagerConfig `json:"manager"`
}

// ManagerConfig ...
type ManagerConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
}

// NewConfig ...
func NewConfig() (*AppConfig, IConfig, error) {
	appConfig := &AppConfig{}
	simpleConfig, err := NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig)

	return appConfig, simpleConfig, err
}
