package gomanager

type IQueue interface {
	Start() error
	Stop() error
	Started() bool
}

// AddQueue ...
func (manager *GoManager) AddQueue(key string, queue IQueue) error {
	manager.queues[key] = queue
	log.Infof("queue %s added", key)

	return nil
}

// RemoveQueue ...
func (manager *GoManager) RemoveQueue(key string) (IQueue, error) {
	queue := manager.queues[key]

	delete(manager.queues, key)
	log.Infof("queue %s removed", key)

	return queue, nil
}

// GetQueue ...
func (manager *GoManager) GetQueue(key string) IQueue {
	if queue, exists := manager.queues[key]; exists {
		return queue
	}
	log.Infof("queue %s doesn't exist", key)
	return nil
}
