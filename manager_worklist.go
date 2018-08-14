package manager

import (
	"sync"
	"time"
)

type IWorkList interface {
	Start(wg *sync.WaitGroup) error
	Stop(wg *sync.WaitGroup) error
	Started() bool
	AddWork(id string, work interface{})
}

// WorkListConfig ...
type WorkListConfig struct {
	Name       string        `json:"name"`
	MaxWorkers int           `json:"max_workers"`
	MaxRetries int           `json:"max_retries"`
	SleepTime  time.Duration `json:"sleep_time"`
	Mode       Mode          `json:"mode"`
}

// NewWorkListConfig...
func NewWorkListConfig(name string, maxWorkers, maxRetries int, sleepTime time.Duration, mode Mode) *WorkListConfig {
	return &WorkListConfig{
		Name:       name,
		MaxWorkers: maxWorkers,
		MaxRetries: maxRetries,
		SleepTime:  sleepTime,
		Mode:       mode,
	}
}

// AddWorkList ...
func (manager *Manager) AddWorkList(key string, worklist IWorkList) error {
	manager.worklist[key] = worklist
	log.Infof("work list %s added", key)

	return nil
}

// RemoveWorkList ...
func (manager *Manager) RemoveWorkList(key string) (IWorkList, error) {
	list := manager.worklist[key]

	delete(manager.worklist, key)
	log.Infof("work list %s removed", key)

	return list, nil
}

// GetWorkList ...
func (manager *Manager) GetWorkList(key string) IWorkList {
	if list, exists := manager.worklist[key]; exists {
		return list
	}
	log.Infof("work list %s doesn't exist", key)
	return nil
}
