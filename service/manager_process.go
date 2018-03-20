package gomanager

// IProcess ...
type IProcess interface {
	Start() error
	Stop() error
	Started() bool
}

// AddProcess ...
func (manager *GoManager) AddProcess(key string, process IProcess) error {
	manager.processes[key] = process
	log.Infof("process %s added", key)

	return nil
}

// RemoveProcess ...
func (manager *GoManager) RemoveProcess(key string) (IProcess, error) {
	process := manager.processes[key]

	delete(manager.processes, key)
	log.Infof("process %s removed", key)

	return process, nil
}

// GetProcess ...
func (manager *GoManager) GetProcess(key string) IProcess {
	if process, exists := manager.processes[key]; exists {
		return process
	}
	log.Infof("process %s doesn't exist", key)
	return nil
}
