package gomanager

// SimpleWorkQueue ...
type SimpleWorkQueue struct {
	name    string
	config  *WorkQueueConfig
	handler WorkHandler
	queue   *Queue
	workers []*Worker
	started bool
}

// NewSimpleWorkQueue ...
func NewSimpleWorkQueue(config *WorkQueueConfig, handler WorkHandler) IWorkQueue {
	return &SimpleWorkQueue{
		name:    config.Name,
		queue:   NewQueue(WithMode(FIFO)),
		config:  config,
		handler: handler,
	}
}

// Start ...
func (workqueue *SimpleWorkQueue) Start() error {
	var workers []*Worker
	for i := 1; i <= workqueue.config.MaxWorkers; i++ {
		log.Infof("starting worker [ %d ]", i)
		worker := NewWorker(i, workqueue.config, workqueue.handler, workqueue.queue)
		worker.Start()
		workers = append(workers, worker)
	}
	workqueue.workers = workers
	workqueue.started = true

	return nil
}

// Stop ...
func (workqueue *SimpleWorkQueue) Stop() error {
	for _, worker := range workqueue.workers {
		log.Infof("stopping worker [ %d: %s ]", worker.id, worker.name)
		worker.Stop()
	}
	return nil
}

// Started ...
func (workqueue *SimpleWorkQueue) Started() bool {
	return workqueue.started
}

// AddWork ...
func (workqueue *SimpleWorkQueue) AddWork(id string, data interface{}) {
	log.Infof("adding work to the list [ name: %s ]", workqueue.name)
	work := NewWork(id, data)
	workqueue.queue.Add(id, work)
}
