package manager

import (
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"sync"

	"github.com/joaosoft/logger"
)

// Manager ...
type Manager struct {
	processes         map[string]IProcess
	configs           map[string]IConfig
	redis             map[string]IRedis
	nsqProducers      map[string]INSQProducer
	nsqConsumers      map[string]INSQConsumer
	rabbitmqProducers map[string]IRabbitmqProducer
	rabbitmqConsumers map[string]IRabbitmqConsumer
	dbs               map[string]IDB
	webs              map[string]IWeb
	gateways          map[string]IGateway
	worklist          map[string]IWorkList
	runInBackground   bool
	config            *ManagerConfig
	logger            logger.ILogger
	isLogExternal     bool

	quit    chan int
	started bool
}

// NewManager ...
func NewManager(options ...ManagerOption) *Manager {
	config, _, err := NewConfig()
	log := logger.NewLogDefault("manager", logger.WarnLevel)

	service := &Manager{
		processes:         make(map[string]IProcess),
		configs:           make(map[string]IConfig),
		redis:             make(map[string]IRedis),
		nsqProducers:      make(map[string]INSQProducer),
		nsqConsumers:      make(map[string]INSQConsumer),
		rabbitmqProducers: make(map[string]IRabbitmqProducer),
		rabbitmqConsumers: make(map[string]IRabbitmqConsumer),
		dbs:               make(map[string]IDB),
		webs:              make(map[string]IWeb),
		gateways:          make(map[string]IGateway),
		worklist:          make(map[string]IWorkList),
		quit:              make(chan int),
		logger:            log,
		config:            config.Manager,
	}

	if err != nil {
		service.logger.Error(err.Error())
	} else if config.Manager != nil {
		level, _ := logger.ParseLevel(config.Manager.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	}

	service.Reconfigure(options...)

	return service
}

// Started ...
func (manager *Manager) Started() bool {
	return manager.started
}

// Start ...
func (manager *Manager) Start() error {
	if !manager.started {
		manager.started = true
		c := make(chan bool, 1)
		if manager.runInBackground {
			go manager.executeStart(c)
			<-c
		} else {
			return manager.executeStart(c)
		}
	}

	return nil
}

// Stop ...
func (manager *Manager) Stop() error {
	if manager.started {
		manager.started = false
		c := make(chan bool)
		if manager.runInBackground {
			go manager.executeStop(c)
			<-c
		} else {
			return manager.executeStop(c)
		}

		return nil
	}

	return nil
}

func (manager *Manager) executeStart(c chan bool) error {
	manager.logger.Info("starting...")

	// listen for termination signals
	termChan := make(chan os.Signal, 1)

	if !manager.runInBackground {
		signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	}

	var wg sync.WaitGroup

	if err := manager.executeAction("start", manager.dbs, &wg); err != nil {
		return err
	}
	if err := manager.executeAction("start", manager.nsqProducers, &wg); err != nil {
		return err
	}
	if err := manager.executeAction("start", manager.nsqConsumers, &wg); err != nil {
		return err
	}
	if err := manager.executeAction("start", manager.rabbitmqProducers, &wg); err != nil {
		return err
	}
	if err := manager.executeAction("start", manager.rabbitmqConsumers, &wg); err != nil {
		return err
	}
	if err := manager.executeAction("start", manager.redis, &wg); err != nil {
		return err
	}
	if err := manager.executeAction("start", manager.worklist, &wg); err != nil {
		return err
	}
	if err := manager.executeAction("start", manager.processes, &wg); err != nil {
		return err
	}
	if err := manager.executeAction("start", manager.webs, &wg); err != nil {
		return err
	}

	wg.Wait()

	if manager.runInBackground {
		c <- true
	}

	manager.logger.Infof("started")

	select {
	case <-termChan:
		manager.logger.Infof("received term signal")
	case <-manager.quit:
		manager.logger.Infof("received shutdown signal")
	}

	return manager.Stop()
}

func (manager *Manager) executeStop(c chan bool) error {
	manager.logger.Info("stopping...")

	var wg sync.WaitGroup

	if err := manager.executeAction("stop", manager.processes, &wg); err != nil {
		return err
	}
	if err := manager.executeAction("stop", manager.worklist, &wg); err != nil {
		return err
	}
	if err := manager.executeAction("stop", manager.webs, &wg); err != nil {
		return err
	}
	if err := manager.executeAction("stop", manager.nsqProducers, &wg); err != nil {
		return err
	}
	if err := manager.executeAction("stop", manager.nsqConsumers, &wg); err != nil {
		return err
	}
	if err := manager.executeAction("stop", manager.rabbitmqProducers, &wg); err != nil {
		return err
	}
	if err := manager.executeAction("stop", manager.rabbitmqConsumers, &wg); err != nil {
		return err
	}
	if err := manager.executeAction("stop", manager.redis, &wg); err != nil {
		return err
	}
	if err := manager.executeAction("stop", manager.dbs, &wg); err != nil {
		return err
	}

	wg.Wait()

	if manager.runInBackground {
		c <- true
	}

	manager.logger.Infof("stopped")

	return nil
}

func (manager *Manager) executeAction(action string, obj interface{}, wg *sync.WaitGroup) error {
	wg.Add(1)
	defer wg.Done()

	objMap := reflect.ValueOf(obj)

	if objMap.Kind() == reflect.Map {
		var wgProcess sync.WaitGroup
		for _, key := range objMap.MapKeys() {
			value := objMap.MapIndex(key)

			started := reflect.ValueOf(value.Interface()).MethodByName("Started").Call([]reflect.Value{})[0]
			switch action {
			case "start":
				if !started.Bool() {
					wgProcess.Add(1)
					go reflect.ValueOf(value.Interface()).MethodByName("Start").Call([]reflect.Value{reflect.ValueOf(&wgProcess)})
					manager.logger.Infof("started [ process: %s ]", key)
				}
			case "stop":
				if started.Bool() {
					wgProcess.Add(1)
					go reflect.ValueOf(value.Interface()).MethodByName("Stop").Call([]reflect.Value{reflect.ValueOf(&wgProcess)})
					manager.logger.Infof("stopped [ process: %s ]", key)
				}
			}
		}
		wgProcess.Wait()
	}

	return nil
}
