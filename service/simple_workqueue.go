package gomanager

// SimpleWorkQueue ...
type SimpleWorkQueue struct {
	name    string
	config  *WorkQueueConfig
	handler WorkHandler
	list    IList
	workers []*Worker
	started bool
}

// NewSimpleWorkQueue ...
func NewSimpleWorkQueue(config *WorkQueueConfig, handler WorkHandler) IWorkQueue {
	return &SimpleWorkQueue{
		name:    config.Name,
		list:    NewQueue(WithMode(config.Mode)),
		config:  config,
		handler: handler,
	}
}

// Start ...
func (workqueue *SimpleWorkQueue) Start() error {
	var workers []*Worker
	for i := 1; i <= workqueue.config.MaxWorkers; i++ {
		log.Infof("starting worker [ %d ]", i)
		worker := NewWorker(i, workqueue.config, workqueue.handler, workqueue.list)
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
	workqueue.list.Add(id, work)
}
