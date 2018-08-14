package manager

import (
	"sync"
	"time"
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

// Worker ...
type Worker struct {
	id         int
	name       string
	handler    WorkHandler
	list       IList
	maxRetries int
	sleepTime  time.Duration
	quit       chan bool
	mux        *sync.Mutex
	started    bool
}

// NewWorker ...
func NewWorker(id int, config *WorkListConfig, handler WorkHandler, list IList) *Worker {
	worker := &Worker{
		id:         id,
		name:       config.Name,
		maxRetries: config.MaxRetries,
		sleepTime:  config.SleepTime,
		handler:    handler,
		list:       list,
		quit:       make(chan bool),
		mux:        &sync.Mutex{},
	}

	return worker
}

// Start ...
func (worker *Worker) Start() error {

	go func() error {
		worker.mux.Lock()
		worker.started = true
		worker.mux.Unlock()

		for {
			select {
			case <-worker.quit:
				log.Debugf("worker quited [name: %s, list size: %d ]", worker.name, worker.list.Size())

				return nil
			default:
				if worker.list.Size() > 0 {
					log.Debugf("worker starting [ name: %d, queue size: %d]", worker.name, worker.list.Size())
					worker.execute()
					log.Debugf("worker finished [ name: %s, queue size: %d]", worker.name, worker.list.Size())
				} else {
					//log.Infof("worker waiting for work to do... [ id: %d, name: %s ]", worker.id, worker.name)
					<-time.After(worker.sleepTime)
				}
			}
		}
	}()

	return nil
}

// Stop ...
func (worker *Worker) Stop() error {
	worker.mux.Lock()
	defer worker.mux.Unlock()

	if worker.list.Size() > 0 {
		log.Infof("stopping worker with tasks in the list [ list size: %d ]", worker.list.Size())
	}
	worker.quit <- true

	return nil
}

// AddWork ...
func (worker *Worker) AddWork(id string, data interface{}) error {
	work := NewWork(id, data)
	return worker.list.Add(id, work)
}

func (worker *Worker) execute() bool {
	var work *Work
	if tmp := worker.list.Remove(); tmp != nil {
		work = tmp.(*Work)
	} else {
		return false
	}

	if err := worker.handler(work.id, work.data); err != nil {
		if work.retries < worker.maxRetries {
			work.retries++
			if err := worker.list.Add(work.id, work); err != nil {
				log.Errorf("error processing the work. re-adding the work to the list [retries: %d, error: %s ]", work.retries, err)
			}
		} else {
			log.Errorf("work discarded of the queue [ retries: %d, error: %s ]", work.retries, err)
		}

		return false
	}
	return true
}
