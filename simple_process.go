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
	if wg != nil {
		defer wg.Done()
	}

	if process.started {
		return nil
	}

	process.started = true
	if err := process.function(); err != nil {
		return err
	}

	return nil
}

// Stop ...
func (process *SimpleProcess) Stop(wg *sync.WaitGroup) error {
	if wg != nil {
		defer wg.Done()
	}

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
