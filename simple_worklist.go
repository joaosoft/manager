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
func (s *SimpleWorkList) Start(waitGroup ...*sync.WaitGroup) (err error) {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
	}

	defer wg.Done()

	if s.started {
		return nil
	}

	var workers []*Worker
	for i := 1; i <= s.config.MaxWorkers; i++ {
		s.logger.Infof("starting worker [ %d ]", i)
		worker := NewWorker(i, s.config, s.handler, s.list, s.workRecoverHandler, s.workRecoverWastedRetriesHandler, s.logger)

		if err = worker.Start(); err != nil {
			s.logger.Errorf("error starting worker [ %d: %s ]: %s", worker.id, worker.name, err)

			for _, w := range workers {
				if err := w.Stop(); err != nil {
					s.logger.Errorf("error stopping worker [ %d: %s ]: %s", w.id, w.name, err)
				}
			}

			return err
		}

		workers = append(workers, worker)
	}

	s.workers = workers
	s.started = true

	return nil
}

// Stop ...
func (s *SimpleWorkList) Stop(waitGroup ...*sync.WaitGroup) error {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
	}

	defer wg.Done()

	if !s.started {
		return nil
	}

	for _, worker := range s.workers {
		s.logger.Infof("stopping worker [ %d: %s ]", worker.id, worker.name)
		if err := worker.Stop(); err != nil {
			s.logger.Errorf("error stopping worker [ %d: %s ]: %s", worker.id, worker.name, err)
		}
	}

	s.started = false

	return nil
}

// Started ...
func (s *SimpleWorkList) Started() bool {
	return s.started
}

// AddWork ...
func (s *SimpleWorkList) AddWork(id string, data interface{}) {
	s.logger.Infof("adding work to the list [ name: %s ]", s.name)
	work := NewWork(id, data, s.logger)
	s.list.Add(id, work)
}
