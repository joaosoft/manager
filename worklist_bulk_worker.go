package manager

import (
	"sync"
	"time"

	"github.com/joaosoft/logger"

	"github.com/labstack/gommon/log"
)

// BulkWorkHandler ...
type BulkWorkHandler func([]*Work) error

// Worker ...
type BulkWorker struct {
	id         int
	name       string
	handler    BulkWorkHandler
	list       IList
	maxWorks   int
	maxRetries int
	sleepTime  time.Duration
	quit       chan bool
	mux        *sync.Mutex
	logger     logger.ILogger
	started    bool
}

// NewBulkWorker ...
func NewBulkWorker(id int, config *BulkWorkListConfig, handler BulkWorkHandler, list IList, logger logger.ILogger) *BulkWorker {
	bulkWorker := &BulkWorker{
		id:         id,
		name:       config.Name,
		maxWorks:   config.MaxWorks,
		maxRetries: config.MaxRetries,
		sleepTime:  config.SleepTime,
		handler:    handler,
		list:       list,
		quit:       make(chan bool),
		mux:        &sync.Mutex{},
		logger:     logger,
	}

	return bulkWorker
}

// Start ...
func (bulkWorker *BulkWorker) Start() error {
	go func() error {
		bulkWorker.mux.Lock()
		bulkWorker.mux.Unlock()

		for {
			select {
			case <-bulkWorker.quit:
				log.Debugf("worker quited [name: %s, list size: %d ]", bulkWorker.name, bulkWorker.list.Size())

				return nil
			default:
				if bulkWorker.list.Size() > 0 {
					log.Debugf("worker starting [ name: %d, queue size: %d]", bulkWorker.name, bulkWorker.list.Size())
					bulkWorker.execute()
					log.Debugf("worker finished [ name: %s, queue size: %d]", bulkWorker.name, bulkWorker.list.Size())
				} else {
					//log.Infof("worker waiting for work to do... [ Id: %d, name: %s ]", bulkWorker.Id, bulkWorker.name)
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
		log.Infof("stopping worker with tasks in the list [ list size: %d ]", bulkWorker.list.Size())
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

func (bulkWorker *BulkWorker) execute() bool {
	var work *Work

	var works []*Work
	for i := 0; i < bulkWorker.maxWorks; i++ {
		if tmp := bulkWorker.list.Remove(); tmp != nil {
			works = append(works, tmp.(*Work))
		} else {
			return false
		}
	}

	if err := bulkWorker.handler(works); err != nil {
		if work.retries < bulkWorker.maxRetries {
			work.retries++
			if err := bulkWorker.list.Add(work.Id, work); err != nil {
				log.Errorf("error processing the work. re-adding the work to the list [retries: %d, error: %s ]", work.retries, err)
			}
		} else {
			log.Errorf("work discarded of the queue [ retries: %d, error: %s ]", work.retries, err)
		}

		return false
	}
	return true
}
