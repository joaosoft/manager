package gomanager

type IWork interface {
	Run() error
}

type IWorkQueue interface {
	Start() error
	Stop() error
	Started() bool
	AddWork(work IWork)
}

// WorkQueueConfig ...
type WorkQueueConfig struct {
	Name        string `json:"name"`
	MaxWorkers  int    `json:"max_workers"`
	MaxLenQueue int    `json:"max_len_queue"`
	Mode        Mode   `json:"mode"`
}

// NewWorkQueueConfig...
func NewWorkQueueConfig(name string, maxWorkers, maxLenQueue int) *WorkQueueConfig {
	return &WorkQueueConfig{
		Name:        name,
		MaxWorkers:  maxWorkers,
		MaxLenQueue: maxLenQueue,
	}
}

// AddQueue ...
func (manager *GoManager) AddWorkQueue(key string, workqueue IWorkQueue) error {
	manager.workqueue[key] = workqueue
	log.Infof("work queue %s added", key)

	return nil
}

// RemoveQueue ...
func (manager *GoManager) RemoveWorkQueue(key string) (IWorkQueue, error) {
	queue := manager.workqueue[key]

	delete(manager.workqueue, key)
	log.Infof("work queue %s removed", key)

	return queue, nil
}

// GetQueue ...
func (manager *GoManager) GetWorkQueue(key string) IWorkQueue {
	if queue, exists := manager.workqueue[key]; exists {
		return queue
	}
	log.Infof("work queue %s doesn't exist", key)
	return nil
}
