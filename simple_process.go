package manager

import "sync"

// SimpleProcess ...
type SimpleProcess struct {
	function func() error
	started  bool
}

// NewSimpleProcess...
func NewSimpleProcess(function func() error) IProcess {
	return &SimpleProcess{
		function: function,
	}
}

// Start ...
func (process *SimpleProcess) Start(wg *sync.WaitGroup) error {
	process.started = true
	if wg != nil {
		defer wg.Done()
	}

	if !process.started {
		if err := process.function(); err != nil {
			return err
		}
	}

	return nil
}

// Stop ...
func (process *SimpleProcess) Stop(wg *sync.WaitGroup) error {
	process.started = false
	if wg != nil {
		defer wg.Done()
	}

	return nil
}

// Started ...
func (process *SimpleProcess) Started() bool {
	return process.started
}
