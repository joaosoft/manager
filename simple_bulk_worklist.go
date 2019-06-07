package manager

import (
	"sync"

	"github.com/joaosoft/logger"
)

// SimpleBulkWorkList ...
type SimpleBulkWorkList struct {
	name                                string
	config                              *BulkWorkListConfig
	handler                             BulkWorkHandler
	bulkWorkRecoverWastedRetriesHandler BulkWorkRecoverWastedRetriesHandler
	bulkWorkRecoverHandler              BulkWorkRecoverHandler
	list                                IList
	workers                             []*BulkWorker
	logger                              logger.ILogger
	started                             bool
}

// NewSimpleBulkWorkList ...
func (manager *Manager) NewSimpleBulkWorkList(config *BulkWorkListConfig, handler BulkWorkHandler, bulkWorkRecoverHandler BulkWorkRecoverHandler, bulkWorkRecoverWastedRetriesHandler BulkWorkRecoverWastedRetriesHandler) IWorkList {
	return &SimpleBulkWorkList{
		name:                                config.Name,
		list:                                manager.NewQueue(WithMode(config.Mode)),
		config:                              config,
		handler:                             handler,
		bulkWorkRecoverHandler:              bulkWorkRecoverHandler,
		bulkWorkRecoverWastedRetriesHandler: bulkWorkRecoverWastedRetriesHandler,
		logger:                              manager.logger,
	}
}

// Start ...
func (bulkWorklist *SimpleBulkWorkList) Start(waitGroup ...*sync.WaitGroup) error {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
	}

	defer wg.Done()

	if bulkWorklist.started {
		return nil
	}

	var workers []*BulkWorker
	for i := 1; i <= bulkWorklist.config.MaxWorkers; i++ {
		bulkWorklist.logger.Infof("starting worker [ %d ]", i)
		worker := NewBulkWorker(i, bulkWorklist.config, bulkWorklist.handler, bulkWorklist.list, bulkWorklist.bulkWorkRecoverHandler, bulkWorklist.bulkWorkRecoverWastedRetriesHandler, bulkWorklist.logger)
		worker.Start()
		workers = append(workers, worker)
	}
	bulkWorklist.workers = workers

	bulkWorklist.started = true

	return nil
}

// Stop ...
func (bulkWorklist *SimpleBulkWorkList) Stop(waitGroup ...*sync.WaitGroup) error {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
	}

	defer wg.Done()

	if !bulkWorklist.started {
		return nil
	}

	for _, worker := range bulkWorklist.workers {
		bulkWorklist.logger.Infof("stopping worker [ %d: %s ]", worker.id, worker.name)
		worker.Stop()
	}

	bulkWorklist.started = false

	return nil
}

// Started ...
func (bulkWorklist *SimpleBulkWorkList) Started() bool {
	return bulkWorklist.started
}

// AddWork ...
func (bulkWorklist *SimpleBulkWorkList) AddWork(id string, data interface{}) {
	bulkWorklist.logger.Infof("adding work to the list [ name: %s ]", bulkWorklist.name)
	work := NewWork(id, data, bulkWorklist.logger)
	bulkWorklist.list.Add(id, work)
}
