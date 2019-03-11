package manager

import "sync"

// INSQProducer ...
type INSQProducer interface {
	Start(waitGroup ...*sync.WaitGroup) error
	Stop(waitGroup ...*sync.WaitGroup) error
	Publish(topic string, body []byte, maxRetries int) error
	Ping() error
	Started() bool
}

// AddNSQProducer ...
func (manager *Manager) AddNSQProducer(key string, nsqProducer INSQProducer) error {
	manager.nsqProducers[key] = nsqProducer
	manager.logger.Infof("nsq producer %s added", key)

	return nil
}

// RemoveNSQProducer ...
func (manager *Manager) RemoveNSQProducer(key string) (INSQProducer, error) {
	process := manager.nsqProducers[key]

	delete(manager.processes, key)
	manager.logger.Infof("nsq producer %s removed", key)

	return process, nil
}

// GetNSQProducer ...
func (manager *Manager) GetNSQProducer(key string) INSQProducer {
	if process, exists := manager.nsqProducers[key]; exists {
		return process
	}
	manager.logger.Infof("nsq producer %s doesn't exist", key)
	return nil
}
