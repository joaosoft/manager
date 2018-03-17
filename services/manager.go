package gomanager

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// GoManager ...
type GoManager struct {
	processes    map[string]IProcess
	configs      map[string]IConfig
	nsqProducers map[string]INSQProducer
	nsqConsumers map[string]INSQConsumer

	dbs      map[string]*DB
	webs     map[string]IWeb
	gateways map[string]IGateway
	queues   map[string]IQueue

	control chan int
	started bool
}

// NewManager ...
func NewManager() (*GoManager, error) {
	return &GoManager{
		processes:    make(map[string]IProcess),
		configs:      make(map[string]IConfig),
		nsqProducers: make(map[string]INSQProducer),
		nsqConsumers: make(map[string]INSQConsumer),
		dbs:          make(map[string]*DB),
		webs:         make(map[string]IWeb),
		gateways:     make(map[string]IGateway),
		queues:       make(map[string]IQueue),
		control:      make(chan int),
	}, nil
}

// Started ...
func (manager *GoManager) Started() bool {
	return manager.started
}

// Start ...
func (manager *GoManager) Start() error {
	log.Infof("starting...")
	manager.started = true

	// listen for termination signals
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	for key, process := range manager.processes {
		log.Infof("starting process [ process: %s ]", key)
		if !process.Started() {
			go process.Start()
		}
		log.Infof("started process [ process: %s ]", key)
	}

	for key, queue := range manager.queues {
		log.Infof("starting queue [ queue: %s ]", key)
		if !queue.Started() {
			go queue.Start()
		}
		log.Infof("started queue [ queue: %s ]", key)
	}

	for key, web := range manager.webs {
		log.Infof("starting web [ web: %s ]", key)
		if !web.Started() {
			go web.Start()
		}
		log.Infof("started web [ web: %s ]", key)
	}

	for key, consumer := range manager.nsqConsumers {
		log.Infof("starting web [ web: %s ]", key)
		if !consumer.Started() {
			go consumer.Start()
		}
		log.Infof("started web [ web: %s ]", key)
	}

	select {
	case <-termChan:
		log.Infof("received term signal")
	case <-manager.control:
		log.Infof("received shutdown signal")
	}

	manager.Stop()

	return nil
}

// Stop ...
func (manager *GoManager) Stop() error {
	if manager.started {
		log.Infof("stopping...")

		for key, process := range manager.processes {
			if process.Started() {
				log.Infof("stopping process [ process: %s ]", key)
				if err := process.Stop(); err != nil {
					log.Error(err, fmt.Sprintf("error stopping process [process:%s]", key))
				}
				log.Infof("close channel [ process: %s ]", key)
				delete(manager.processes, key)
				log.Infof("stopped process [ process: %s ]", key)
			}
		}

		manager.started = true
		log.Infof("stopped")
	}

	return nil
}
