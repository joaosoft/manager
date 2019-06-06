package manager

import (
	"sync"
	"time"
)

type IWorkList interface {
	Start(waitGroup ...*sync.WaitGroup) error
	Stop(waitGroup ...*sync.WaitGroup) error
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

// BulkWorkListConfig ...
type BulkWorkListConfig struct {
	Name       string        `json:"name"`
	MaxWorks   int           `json:"max_works"`
	MaxWorkers int           `json:"max_workers"`
	MaxRetries int           `json:"max_retries"`
	SleepTime  time.Duration `json:"sleep_time"`
	Mode       Mode          `json:"mode"`
}

// NewBulkWorkListConfig...
func NewBulkWorkListConfig(name string, maxWorks, maxWorkers, maxRetries int, sleepTime time.Duration, mode Mode) *BulkWorkListConfig {
	return &BulkWorkListConfig{
		Name:       name,
		MaxWorks:   maxWorks,
		MaxWorkers: maxWorkers,
		MaxRetries: maxRetries,
		SleepTime:  sleepTime,
		Mode:       mode,
	}
}

// AddWorkList ...
func (manager *Manager) AddWorkList(key string, worklist IWorkList) error {
	manager.worklist[key] = worklist
	manager.logger.Infof("work list %s added", key)

	return nil
}

// RemoveWorkList ...
func (manager *Manager) RemoveWorkList(key string) (IWorkList, error) {
	list := manager.worklist[key]

	delete(manager.worklist, key)
	manager.logger.Infof("work list %s removed", key)

	return list, nil
}

// GetWorkList ...
func (manager *Manager) GetWorkList(key string) IWorkList {
	if list, exists := manager.worklist[key]; exists {
		return list
	}
	manager.logger.Infof("work list %s doesn't exist", key)
	return nil
}
