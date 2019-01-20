package manager

import (
	"github.com/joaosoft/logger"
	"sync"
)

// SimpleWorkList ...
type SimpleWorkList struct {
	name    string
	config  *WorkListConfig
	handler WorkHandler
	list    IList
	workers []*Worker
	logger logger.ILogger
	started bool
}

// NewSimpleWorkList ...
func (manager *Manager) NewSimpleWorkList(config *WorkListConfig, handler WorkHandler) IWorkList {
	return &SimpleWorkList{
		name:    config.Name,
		list:    NewQueue(WithMode(config.Mode)),
		config:  config,
		handler: handler,
		logger: manager.logger,
	}
}

// Start ...
func (worklist *SimpleWorkList) Start(wg *sync.WaitGroup) error {
	if wg == nil {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	}

	defer wg.Done()

	if worklist.started {
		return nil
	}

	var workers []*Worker
	for i := 1; i <= worklist.config.MaxWorkers; i++ {
		worklist.logger.Infof("starting worker [ %d ]", i)
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
	if wg == nil {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	}

	defer wg.Done()

	if !worklist.started {
		return nil
	}

	for _, worker := range worklist.workers {
		worklist.logger.Infof("stopping worker [ %d: %s ]", worker.id, worker.name)
		worker.Stop()
	}

	worklist.started = false

	return nil
}

// Started ...
func (worklist *SimpleWorkList) Started() bool {
	return worklist.started
}

// AddWork ...
func (worklist *SimpleWorkList) AddWork(id string, data interface{}) {
	worklist.logger.Infof("adding work to the list [ name: %s ]", worklist.name)
	work := NewWork(id, data)
	worklist.list.Add(id, work)
}
