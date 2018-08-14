package manager

import "sync"

// IRabbitmqProducer ...
type IRabbitmqProducer interface {
	Start(wg *sync.WaitGroup) error
	Stop(wg *sync.WaitGroup) error
	Publish(routingKey string, body []byte, reliable bool) error
	Started() bool
}

// AddRabbitmqProducer ...
func (manager *Manager) AddRabbitmqProducer(key string, nsqProducer IRabbitmqProducer) error {
	manager.rabbitmqProducers[key] = nsqProducer
	log.Infof("nsq producer %s added", key)

	return nil
}

// RemoveRabbitmqProducer ...
func (manager *Manager) RemoveRabbitmqProducer(key string) (IRabbitmqProducer, error) {
	process := manager.rabbitmqProducers[key]

	delete(manager.processes, key)
	log.Infof("nsq producer %s removed", key)

	return process, nil
}

// GetRabbitmqProducer ...
func (manager *Manager) GetRabbitmqProducer(key string) IRabbitmqProducer {
	if process, exists := manager.rabbitmqProducers[key]; exists {
		return process
	}
	log.Infof("nsq producer %s doesn't exist", key)
	return nil
}
