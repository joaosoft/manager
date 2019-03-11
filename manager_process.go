package manager

import "sync"

// IProcess ...
type IProcess interface {
	Start(waitGroup ...*sync.WaitGroup) error
	Stop(waitGroup ...*sync.WaitGroup) error
	Started() bool
}

// AddProcess ...
func (manager *Manager) AddProcess(key string, process IProcess) error {
	manager.processes[key] = process
	manager.logger.Infof("process %s added", key)

	return nil
}

// RemoveProcess ...
func (manager *Manager) RemoveProcess(key string) (IProcess, error) {
	process := manager.processes[key]

	delete(manager.processes, key)
	manager.logger.Infof("process %s removed", key)

	return process, nil
}

// GetProcess ...
func (manager *Manager) GetProcess(key string) IProcess {
	if process, exists := manager.processes[key]; exists {
		return process
	}
	manager.logger.Infof("process %s doesn't exist", key)
	return nil
}
