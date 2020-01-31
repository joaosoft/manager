package manager

import (
	"sync"
	"time"

	"github.com/joaosoft/logger"
)

// IList ...
type IList interface {
	Add(id string, data interface{}) error
	Remove(ids ...string) interface{}
	Size() int
	IsEmpty() bool
	Dump() string
}

// WorkHandler ...
type WorkHandler func(id string, data interface{}) error

// WorkRecoverHandler ...
type WorkRecoverHandler func(list IList) error

// WorkRecoverWastedRetriesHandler ...
type WorkRecoverWastedRetriesHandler func(id string, data interface{}) error

// Worker ...
type Worker struct {
	id                              int
	name                            string
	handler                         WorkHandler
	workRecoverHandler              WorkRecoverHandler
	workRecoverWastedRetriesHandler WorkRecoverWastedRetriesHandler
	list                            IList
	maxRetries                      int
	sleepTime                       time.Duration
	quit                            chan bool
	mux                             *sync.Mutex
	logger                          logger.ILogger
	started                         bool
}

// NewWorker ...
func NewWorker(id int, config *WorkListConfig, handler WorkHandler, list IList, workRecoverHandler WorkRecoverHandler, workRecoverWastedRetriesHandler WorkRecoverWastedRetriesHandler, logger logger.ILogger) *Worker {
	worker := &Worker{
		id:                              id,
		name:                            config.Name,
		maxRetries:                      config.MaxRetries,
		sleepTime:                       config.SleepTime,
		handler:                         handler,
		workRecoverHandler:              workRecoverHandler,
		workRecoverWastedRetriesHandler: workRecoverWastedRetriesHandler,
		list:                            list,
		quit:                            make(chan bool),
		mux:                             &sync.Mutex{},
		logger:                          logger,
	}

	return worker
}

// Start ...
func (worker *Worker) Start() error {
	go func() error {
		for {
			select {
			case <-worker.quit:
				logger.Debugf("worker quited [name: %s, list size: %d ]", worker.name, worker.list.Size())

				return nil
			default:
				if worker.list.Size() > 0 {
					logger.Debugf("worker starting [ name: %d, queue size: %d]", worker.name, worker.list.Size())
					if err := worker.execute(); err != nil {
						logger.Errorf("worker errored [ name: %s, queue size: %d]", worker.name, worker.list.Size())
					}
					logger.Debugf("worker finished [ name: %s, queue size: %d]", worker.name, worker.list.Size())
				} else {
					logger.Debugf("worker waiting for work to do... [ id: %d, name: %s ]", worker.id, worker.name)
					<-time.After(worker.sleepTime)
				}
			}
		}
	}()

	worker.started = true

	return nil
}

// Stop ...
func (worker *Worker) Stop() error {

	worker.mux.Lock()
	defer worker.mux.Unlock()

	if worker.list.Size() > 0 {
		logger.Infof("stopping worker with tasks in the list [ list size: %d ]", worker.list.Size())
	}
	worker.quit <- true

	worker.started = false

	return nil
}

// AddWork ...
func (worker *Worker) AddWork(id string, data interface{}) error {
	work := NewWork(id, data, worker.logger)
	return worker.list.Add(id, work)
}

func (worker *Worker) execute() error {
	var work *Work

	defer func() {
		if worker.workRecoverHandler != nil {
			if r := recover(); r != nil {
				logger.Debug("recovering worker data")
				if err := worker.workRecoverHandler(worker.list); err != nil {
					logger.Errorf("error processing recovering of worker. [ error: %s ]", err)
				}
			}
		}
	}()

	if tmp := worker.list.Remove(); tmp != nil {
		work = tmp.(*Work)
	}

	if err := worker.handler(work.Id, work.Data); err != nil {
		if work.retries < worker.maxRetries {
			work.retries++
			if err := worker.list.Add(work.Id, work); err != nil {
				logger.Errorf("error processing the work. re-adding the work to the list [retries: %d, error: %s ]", work.retries, err).ToError()
			}
			logger.Errorf("work requeued of the queue [ retries: %d, error: %s ]", work.retries, err).ToError()

		} else {
			if worker.workRecoverWastedRetriesHandler != nil {
				if err := worker.workRecoverWastedRetriesHandler(work.Id, work.Data); err != nil {
					logger.Errorf("error processing recovering one of worker. [ error: %s ]", err).ToError()
				}
			}
			logger.Errorf("work discarded of the queue [ retries: %d, error: %s ]", work.retries, err).ToError()

		}

		return nil
	}

	return nil
}
