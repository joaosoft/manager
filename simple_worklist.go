package manager

import "sync"

// SimpleWorkList ...
type SimpleWorkList struct {
	name    string
	config  *WorkListConfig
	handler WorkHandler
	list    IList
	workers []*Worker
	started bool
}

// NewSimpleWorkList ...
func NewSimpleWorkList(config *WorkListConfig, handler WorkHandler) IWorkList {
	return &SimpleWorkList{
		name:    config.Name,
		list:    NewQueue(WithMode(config.Mode)),
		config:  config,
		handler: handler,
	}
}

// Start ...
func (worklist *SimpleWorkList) Start(wg *sync.WaitGroup) error {
	defer wg.Done()

	var workers []*Worker
	for i := 1; i <= worklist.config.MaxWorkers; i++ {
		log.Infof("starting worker [ %d ]", i)
		worker := NewWorker(i, worklist.config, worklist.handler, worklist.list)
		worker.Start()
		workers = append(workers, worker)
	}
	worklist.workers = workers
	worklist.started = true

	return nil
}

// Stop ...
func (worklist *SimpleWorkList) Stop(wg *sync.WaitGroup) error {
	defer wg.Done()
	for _, worker := range worklist.workers {
		log.Infof("stopping worker [ %d: %s ]", worker.id, worker.name)
		worker.Stop()
	}
	return nil
}

// Started ...
func (worklist *SimpleWorkList) Started() bool {
	return worklist.started
}

// AddWork ...
func (worklist *SimpleWorkList) AddWork(id string, data interface{}) {
	log.Infof("adding work to the list [ name: %s ]", worklist.name)
	work := NewWork(id, data)
	worklist.list.Add(id, work)
}
