package web

import (
	"fmt"
)

type ClientConfig struct {
	Log Log `json:"log"`
}

func NewClientConfig() (*AppClientConfig, error) {
	appConfig := &AppClientConfig{}
	err := NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig)

	return appConfig, err
}
