package gomanager

import "time"

type IWorkQueue interface {
	Start() error
	Stop() error
	Started() bool
	AddWork(id string, work interface{})
}

// WorkQueueConfig ...
type WorkQueueConfig struct {
	Name       string        `json:"name"`
	MaxWorkers int           `json:"max_workers"`
	MaxRetries int           `json:"max_retries"`
	SleepTime  time.Duration `json:"sleep_time"`
	Mode       Mode          `json:"mode"`
}

// NewWorkQueueConfig...
func NewWorkQueueConfig(name string, maxWorkers, maxRetries int, sleepTime time.Duration, mode Mode) *WorkQueueConfig {
	return &WorkQueueConfig{
		Name:       name,
		MaxWorkers: maxWorkers,
		MaxRetries: maxRetries,
		SleepTime:  sleepTime,
		Mode:       mode,
	}
}

// AddQueue ...
func (manager *GoManager) AddWorkQueue(key string, workqueue IWorkQueue) error {
	manager.workqueue[key] = workqueue
	log.Infof("work list %s added", key)

	return nil
}

// RemoveQueue ...
func (manager *GoManager) RemoveWorkQueue(key string) (IWorkQueue, error) {
	queue := manager.workqueue[key]

	delete(manager.workqueue, key)
	log.Infof("work list %s removed", key)

	return queue, nil
}

// GetQueue ...
func (manager *GoManager) GetWorkQueue(key string) IWorkQueue {
	if queue, exists := manager.workqueue[key]; exists {
		return queue
	}
	log.Infof("work list %s doesn't exist", key)
	return nil
}
