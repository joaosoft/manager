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

	wg.Done()

	if !process.started {
		if err := process.function(); err != nil {
			return err
		}
		process.started = true
	}

	return nil
}

// Stop ...
func (process *SimpleProcess) Stop(wg *sync.WaitGroup) error {
	defer wg.Done()

	if process.started {
		process.started = false
	}
	return nil
}

// Started ...
func (process *SimpleProcess) Started() bool {
	return process.started
}
