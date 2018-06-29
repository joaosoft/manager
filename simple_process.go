package manager

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
func (process *SimpleProcess) Start() error {
	if !process.started {
		if err := process.function(); err != nil {
			return err
		}
		process.started = true
	}

	return nil
}

// Stop ...
func (process *SimpleProcess) Stop() error {
	if process.started {
		process.started = false
	}
	return nil
}

// Started ...
func (process *SimpleProcess) Started() bool {
	return process.started
}
