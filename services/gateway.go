package mgr

import (
	"fmt"
	"io"

	"github.com/joaosoft/go-manager/services/gateway"
	"github.com/labstack/gommon/log"
)

// -------------- GATEWAY --------------
// NewGateway ... creates a new web server
func (instance *Manager) NewGateway(config *gateway.Config) (*gateway.Gateway, error) {
	log.Infof(fmt.Sprintf("gateway, creating gateway"))
	return gateway.NewGateway(config), nil
}

// -------------- METHODS --------------
// GetGateway ... get a gateway by key
func (instance *Manager) GetGateway(key string) (*gateway.Gateway, error) {
	return instance.GatewayController[key], nil
}

// AddGateway, add a new gateway
func (instance *Manager) AddGateway(key string, gateway *gateway.Gateway) error {
	log.Infof(fmt.Sprintf("gateway, add a new gateway '%s'", key))
	instance.GatewayController[key] = gateway
	return nil
}

// RemGateway, remove a gateway by key
func (instance *Manager) RemGateway(key string) (*gateway.Gateway, error) {
	log.Infof(fmt.Sprintf("gateway, remove the gateway '%s'", key))

	// get gateway
	controller := instance.GatewayController[key]

	// delete gateway
	delete(instance.GatewayController, key)
	log.Infof(fmt.Sprintf("gateway, gateway '%s' removed", key))

	return controller, nil
}

// RequestGateway ... make a http request
func (instance *Manager) RequestGateway(key string, method string, endpoint string, headers map[string]string, body io.Reader) (int, []byte, error) {
	log.Infof(fmt.Sprintf("gateway, request the gateway '%s' method:'%s', endpoint: '%s', headers:'%s', body:'%s'", key, method, endpoint, headers, body))
	return instance.GatewayController[key].Request(method, endpoint, headers, body)
}
