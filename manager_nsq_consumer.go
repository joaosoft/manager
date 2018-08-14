package manager

import (
	"github.com/nsqio/go-nsq"
	"sync"
)

type INSQHandler interface {
	HandleMessage(message *nsq.Message) error
}

// INSQConsumer ...
type INSQConsumer interface {
	Start(wg *sync.WaitGroup) error
	Stop(wg *sync.WaitGroup) error
	HandleMessage(message *nsq.Message) error
	Started() bool
}

// AddNSQConsumer ...
func (manager *Manager) AddNSQConsumer(key string, nsqConsumer INSQConsumer) error {
	manager.nsqConsumers[key] = nsqConsumer
	log.Infof("consumer %s added", key)

	return nil
}

// RemoveNSQConsumer ...
func (manager *Manager) RemoveNSQConsumer(key string) (INSQConsumer, error) {
	nsqConsumer := manager.nsqConsumers[key]

	delete(manager.processes, key)
	log.Infof("consumer %s removed", key)

	return nsqConsumer, nil
}

// GetNSQConsumer ...
func (manager *Manager) GetNSQConsumer(key string) INSQConsumer {
	if nsqConsumer, exists := manager.nsqConsumers[key]; exists {
		return nsqConsumer
	}
	log.Infof("consumer %s doesn't exist", key)
	return nil
}
