package gomanager

// SimpleWorkQueue ...
type SimpleWorkQueue struct {
	name        string
	queue       Queue
	maxWorkers  int
	maxLenQueue int
	started     bool
}

// NewSimpleWorkQueue ...
func NewSimpleWorkQueue(config *WorkQueueConfig) IWorkQueue {
	return &SimpleWorkQueue{
		name:        config.Name,
		maxWorkers:  config.MaxWorkers,
		maxLenQueue: config.MaxLenQueue,
		queue: Queue{
			mode: config.Mode,
		},
	}
}

// Start ...
func (workqueue *SimpleWorkQueue) Start() error {
	if !workqueue.started {
		workqueue.started = true
	}
	return nil
}

// Stop ...
func (workqueue *SimpleWorkQueue) Stop() error {
	if workqueue.started {
		workqueue.started = false
	}
	return nil
}

// Started ...
func (workqueue *SimpleWorkQueue) Started() bool {
	return workqueue.started
}

// AddWork ...
func (workqueue *SimpleWorkQueue) AddWork(work IWork) {
	log.Infof("adding work to queue [ name: %s ]", workqueue.name)
}
