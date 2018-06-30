package manager

// IGateway ...
type IGateway interface {
	Request(method, host, endpoint string, headers map[string][]string, body []byte) (int, []byte, error)
}

// AddGateway ...
func (manager *Manager) AddGateway(key string, gateway IGateway) error {
	manager.gateways[key] = gateway
	log.Infof("gateway %s added", key)

	return nil
}

// RemoveGateway ...
func (manager *Manager) RemoveGateway(key string) (IGateway, error) {
	gateway := manager.gateways[key]

	delete(manager.configs, key)
	log.Infof("gateway %s removed", key)

	return gateway, nil
}

// GetGateway ...
func (manager *Manager) GetGateway(key string) IGateway {
	if gateway, exists := manager.gateways[key]; exists {
		return gateway
	}
	log.Infof("gateway %s doesn't exist", key)
	return nil
}
