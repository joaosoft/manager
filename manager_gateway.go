package manager

// IGateway ...
type IGateway interface {
	Request(method, host, endpoint string, contentType string, headers map[string][]string, body []byte) (int, []byte, error)
}

// AddGateway ...
func (manager *Manager) AddGateway(key string, gateway IGateway) error {
	manager.gateways[key] = gateway
	manager.logger.Infof("gateway %s added", key)

	return nil
}

// RemoveGateway ...
func (manager *Manager) RemoveGateway(key string) (IGateway, error) {
	gateway := manager.gateways[key]

	delete(manager.configs, key)
	manager.logger.Infof("gateway %s removed", key)

	return gateway, nil
}

// GetGateway ...
func (manager *Manager) GetGateway(key string) IGateway {
	if gateway, exists := manager.gateways[key]; exists {
		return gateway
	}
	manager.logger.Infof("gateway %s doesn't exist", key)
	return nil
}
