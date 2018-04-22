package gomanager

import (
	"github.com/nsqio/go-nsq"
)

type INSQHandler interface {
	HandleMessage(message *nsq.Message) error
}

// INSQConsumer ...
type INSQConsumer interface {
	Start() error
	Stop() error
	HandleMessage(message *nsq.Message) error
	Started() bool
}

// AddProcess ...
func (manager *Manager) AddNSQConsumer(key string, nsqConsumer INSQConsumer) error {
	manager.nsqConsumers[key] = nsqConsumer
	log.Infof("consumer %s added", key)

	return nil
}

// RemoveProcess ...
func (manager *Manager) RemoveNSQConsumer(key string) (INSQConsumer, error) {
	nsqConsumers := manager.nsqConsumers[key]

	delete(manager.processes, key)
	log.Infof("consumer %s removed", key)

	return nsqConsumers, nil
}

// GetProcess ...
func (manager *Manager) GetNSQConsumer(key string) INSQConsumer {
	if nsqConsumers, exists := manager.nsqConsumers[key]; exists {
		return nsqConsumers
	}
	log.Infof("consumer %s doesn't exist", key)
	return nil
}
