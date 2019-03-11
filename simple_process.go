package manager

import (
	"github.com/joaosoft/logger"
	"sync"
)

// SimpleProcess ...
type SimpleProcess struct {
	function func() error
	logger logger.ILogger
	started  bool
}

// NewSimpleProcess...
func (manager *Manager) NewSimpleProcess(function func() error) IProcess {
	return &SimpleProcess{
		function: function,
		logger: manager.logger,
	}
}

// Start ...
func (process *SimpleProcess) Start(waitGroup ...*sync.WaitGroup) error {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
	}

	defer wg.Done()

	if process.started {
		return nil
	}

	if err := process.function(); err != nil {
		return err
	}

	process.started = true

	return nil
}

// Stop ...
func (process *SimpleProcess) Stop(waitGroup ...*sync.WaitGroup) error {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
	}

	defer wg.Done()

	if !process.started {
		return nil
	}

	process.started = false

	return nil
}

// Started ...
func (process *SimpleProcess) Started() bool {
	return process.started
}
