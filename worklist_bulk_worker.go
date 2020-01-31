package manager

import (
	"sync"
	"time"

	"github.com/joaosoft/logger"
)

// BulkWorkHandler ...
type BulkWorkHandler func([]*Work) error

// BulkWorkRecoverHandler ...
type BulkWorkRecoverHandler func(list IList) error

// BulkWorkRecoverWastedRetriesHandler ...
type BulkWorkRecoverWastedRetriesHandler func(id string, data interface{}) error

// Worker ...
type BulkWorker struct {
	id                          int
	name                        string
	handler                     BulkWorkHandler
	recoverHandler              BulkWorkRecoverHandler
	recoverWastedRetriesHandler BulkWorkRecoverWastedRetriesHandler
	list                        IList
	maxWorks                    int
	maxRetries                  int
	sleepTime                   time.Duration
	quit                        chan bool
	mux                         *sync.Mutex
	logger                      logger.ILogger
	started                     bool
}

// NewBulkWorker ...
func NewBulkWorker(id int, config *BulkWorkListConfig, handler BulkWorkHandler, list IList, bulkWorkRecoverHandler BulkWorkRecoverHandler, bulkWorkRecoverOneHandler BulkWorkRecoverWastedRetriesHandler, logger logger.ILogger) *BulkWorker {
	bulkWorker := &BulkWorker{
		id:                          id,
		name:                        config.Name,
		maxWorks:                    config.MaxWorks,
		maxRetries:                  config.MaxRetries,
		sleepTime:                   config.SleepTime,
		handler:                     handler,
		recoverHandler:              bulkWorkRecoverHandler,
		recoverWastedRetriesHandler: bulkWorkRecoverOneHandler,
		list:                        list,
		quit:                        make(chan bool),
		mux:                         &sync.Mutex{},
		logger:                      logger,
	}

	return bulkWorker
}

// Start ...
func (bulkWorker *BulkWorker) Start() error {
	go func() error {
		for {
			select {
			case <-bulkWorker.quit:
				logger.Debugf("worker quited [name: %s, list size: %d ]", bulkWorker.name, bulkWorker.list.Size())

				return nil
			default:
				if bulkWorker.list.Size() > 0 {
					logger.Debugf("worker starting [ name: %d, queue size: %d]", bulkWorker.name, bulkWorker.list.Size())
					bulkWorker.execute()
					logger.Debugf("worker finished [ name: %s, queue size: %d]", bulkWorker.name, bulkWorker.list.Size())
				} else {
					logger.Debugf("worker waiting for work to do... [ id: %d, name: %s ]", bulkWorker.id, bulkWorker.name)
					<-time.After(bulkWorker.sleepTime)
				}
			}
		}
	}()

	bulkWorker.started = true

	return nil
}

// Stop ...
func (bulkWorker *BulkWorker) Stop() error {

	bulkWorker.mux.Lock()
	defer bulkWorker.mux.Unlock()

	if bulkWorker.list.Size() > 0 {
		logger.Infof("stopping worker with tasks in the list [ list size: %d ]", bulkWorker.list.Size())
	}
	bulkWorker.quit <- true

	bulkWorker.started = false

	return nil
}

// AddWork ...
func (bulkWorker *BulkWorker) AddWork(id string, data interface{}) error {
	work := NewWork(id, data, bulkWorker.logger)
	return bulkWorker.list.Add(id, work)
}

func (bulkWorker *BulkWorker) execute() error {
	defer func() {
		if bulkWorker.recoverHandler != nil {
			if r := recover(); r != nil {
				logger.Debug("recovering worker data")
				if err := bulkWorker.recoverHandler(bulkWorker.list); err != nil {
					logger.Errorf("error processing recovering of worker. [ error: %s ]", err)
				}
			}
		}
	}()

	var works []*Work
	for i := 0; i < bulkWorker.maxWorks; i++ {
		if tmp := bulkWorker.list.Remove(); tmp != nil {
			if tmp == nil {
				break
			}
			works = append(works, tmp.(*Work))
		}
	}

	if err := bulkWorker.handler(works); err != nil {
		for _, work := range works {
			if work.retries < bulkWorker.maxRetries {
				work.retries++
				if err := bulkWorker.list.Add(work.Id, work); err != nil {
					logger.Errorf("error processing the work. re-adding the work to the list [retries: %d, error: %s ]", work.retries, err)
				}
				logger.Errorf("work requeued of the queue [ retries: %d, error: %s ]", work.retries, err).ToError()
			} else {
				if bulkWorker.recoverWastedRetriesHandler != nil {
					if err := bulkWorker.recoverWastedRetriesHandler(work.Id, work.Data); err != nil {
						logger.Errorf("error processing recovering one of worker. [ error: %s ]", err).ToError()
					}
				}
				logger.Errorf("work discarded of the queue [ retries: %d, error: %s ]", work.retries, err)
			}
		}

		return nil
	}
	return nil
}
