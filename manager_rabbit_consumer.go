package manager

import (
	"github.com/streadway/amqp"
	"sync"
)

type RabbitmqHandler func(message amqp.Delivery) error

// IRabbitmqConsumer ...
type IRabbitmqConsumer interface {
	Start(wg *sync.WaitGroup) error
	Stop(wg *sync.WaitGroup) error
	Started() bool
}

// AddRabbitmqConsumer ...
func (manager *Manager) AddRabbitmqConsumer(key string, rabbitmqConsumer IRabbitmqConsumer) error {
	manager.rabbitmqConsumers[key] = rabbitmqConsumer
	log.Infof("consumer %s added", key)

	return nil
}

// RemoveRabbitmqConsumer ...
func (manager *Manager) RemoveRabbitmqConsumer(key string) (IRabbitmqConsumer, error) {
	rabbitmqConsumer := manager.rabbitmqConsumers[key]

	delete(manager.processes, key)
	log.Infof("consumer %s removed", key)

	return rabbitmqConsumer, nil
}

// GetRabbitmqConsumer ...
func (manager *Manager) GetRabbitmqConsumer(key string) IRabbitmqConsumer {
	if rabbitmqConsumer, exists := manager.rabbitmqConsumers[key]; exists {
		return rabbitmqConsumer
	}
	log.Infof("consumer %s doesn't exist", key)
	return nil
}
