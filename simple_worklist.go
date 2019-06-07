package manager

import (
	"sync"

	"github.com/joaosoft/logger"
)

// SimpleWorkList ...
type SimpleWorkList struct {
	name                            string
	config                          *WorkListConfig
	handler                         WorkHandler
	workRecoverHandler              WorkRecoverHandler
	workRecoverWastedRetriesHandler WorkRecoverWastedRetriesHandler
	list                            IList
	workers                         []*Worker
	logger                          logger.ILogger
	started                         bool
}

// NewSimpleWorkList ...
func (manager *Manager) NewSimpleWorkList(config *WorkListConfig, handler WorkHandler, workRecoverHandler WorkRecoverHandler, workRecoverWastedRetriesHandler WorkRecoverWastedRetriesHandler) IWorkList {
	return &SimpleWorkList{
		name:                            config.Name,
		list:                            manager.NewQueue(WithMode(config.Mode)),
		config:                          config,
		handler:                         handler,
		workRecoverHandler:              workRecoverHandler,
		workRecoverWastedRetriesHandler: workRecoverWastedRetriesHandler,
		logger:                          manager.logger,
	}
}

// Start ...
func (worklist *SimpleWorkList) Start(waitGroup ...*sync.WaitGroup) error {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
	}

	defer wg.Done()

	if worklist.started {
		return nil
	}

	var workers []*Worker
	for i := 1; i <= worklist.config.MaxWorkers; i++ {
		worklist.logger.Infof("starting worker [ %d ]", i)
		worker := NewWorker(i, worklist.config, worklist.handler, worklist.list, worklist.workRecoverHandler, worklist.workRecoverWastedRetriesHandler, worklist.logger)
		worker.Start()
		workers = append(workers, worker)
	}
	worklist.workers = workers

	worklist.started = true

	return nil
}

// Stop ...
func (worklist *SimpleWorkList) Stop(waitGroup ...*sync.WaitGroup) error {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
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
	work := NewWork(id, data, worklist.logger)
	worklist.list.Add(id, work)
}
