package web

import (
	"fmt"
)

type ServerConfig struct {
	Address string `json:"address"`
	Log     struct {
		Level string `json:"level"`
	} `json:"log"`
}

func NewServerConfig() (*AppServerConfig, error) {
	appConfig := &AppServerConfig{}
	err := NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig)

	return appConfig, err
}
