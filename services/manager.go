package gomanager

import (
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

// GoManager ...
type GoManager struct {
	processes    map[string]IProcess
	configs      map[string]IConfig
	redis        map[string]IRedis
	nsqProducers map[string]INSQProducer
	nsqConsumers map[string]INSQConsumer
	dbs          map[string]IDB
	webs         map[string]IWeb
	gateways     map[string]IGateway
	workqueue    map[string]IWorkQueue

	control chan int
	started bool
}

// NewManager ...
func NewManager() *GoManager {
	return &GoManager{
		processes:    make(map[string]IProcess),
		configs:      make(map[string]IConfig),
		redis:        make(map[string]IRedis),
		nsqProducers: make(map[string]INSQProducer),
		nsqConsumers: make(map[string]INSQConsumer),
		dbs:          make(map[string]IDB),
		webs:         make(map[string]IWeb),
		gateways:     make(map[string]IGateway),
		workqueue:    make(map[string]IWorkQueue),
		control:      make(chan int),
	}
}

// Started ...
func (manager *GoManager) Started() bool {
	return manager.started
}

// Start ...
func (manager *GoManager) Start() error {
	log.Infof("starting...")

	// listen for termination signals
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	executeAction("start", manager.processes)
	executeAction("start", manager.workqueue)
	executeAction("start", manager.webs)
	executeAction("start", manager.nsqProducers)
	executeAction("start", manager.nsqConsumers)
	executeAction("start", manager.dbs)
	executeAction("start", manager.redis)

	manager.started = true
	log.Infof("started")

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

		executeAction("stop", manager.processes)
		executeAction("stop", manager.workqueue)
		executeAction("stop", manager.webs)
		executeAction("stop", manager.nsqProducers)
		executeAction("stop", manager.nsqConsumers)
		executeAction("stop", manager.dbs)
		executeAction("stop", manager.redis)

		manager.started = false
		log.Infof("stopped")
	}

	return nil
}

func executeAction(action string, obj interface{}) error {
	objMap := reflect.ValueOf(obj)

	if objMap.Kind() == reflect.Map {
		for _, key := range objMap.MapKeys() {
			value := objMap.MapIndex(key)

			started := reflect.ValueOf(value.Interface()).MethodByName("Started").Call([]reflect.Value{})[0]
			switch action {
			case "start":
				if !started.Bool() {
					go reflect.ValueOf(value.Interface()).MethodByName("Start").Call([]reflect.Value{})
					log.Infof("started [ process: %s ]", key)
				}
			case "stop":
				if started.Bool() {
					go reflect.ValueOf(value.Interface()).MethodByName("Stop").Call([]reflect.Value{})
					log.Infof("stopped [ process: %s ]", key)
				}
			}
		}
	}

	return nil
}
