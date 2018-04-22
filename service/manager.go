package gomanager

import (
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/joaosoft/go-log/service"

	"fmt"
)

// Manager ...
type Manager struct {
	processes       map[string]IProcess
	configs         map[string]IConfig
	redis           map[string]IRedis
	nsqProducers    map[string]INSQProducer
	nsqConsumers    map[string]INSQConsumer
	dbs             map[string]IDB
	webs            map[string]IWeb
	gateways        map[string]IGateway
	worklist        map[string]IWorkList
	runInBackground bool
	config          *ManagerConfig
	isLogExternal   bool

	quit    chan int
	started bool
}

// NewManager ...
func NewManager(options ...managerOption) *Manager {
	gomanager := &Manager{
		processes:    make(map[string]IProcess),
		configs:      make(map[string]IConfig),
		redis:        make(map[string]IRedis),
		nsqProducers: make(map[string]INSQProducer),
		nsqConsumers: make(map[string]INSQConsumer),
		dbs:          make(map[string]IDB),
		webs:         make(map[string]IWeb),
		gateways:     make(map[string]IGateway),
		worklist:     make(map[string]IWorkList),
		quit:         make(chan int),
	}

	gomanager.Reconfigure(options...)

	// load configuration file
	configApp := &appConfig{}
	if _, err := readFile(fmt.Sprintf("/config/app.%s.json", getEnv()), configApp); err != nil {
		log.Error(err)
	} else {
		level, _ := golog.ParseLevel(configApp.GoManager.Log.Level)
		log.Debugf("setting log level to %s", level)
		log.Reconfigure(golog.WithLevel(level))
	}
	gomanager.config = &configApp.GoManager

	return gomanager
}

// Started ...
func (manager *Manager) Started() bool {
	return manager.started
}

// Start ...
func (manager *Manager) Start() error {
	if manager.runInBackground {
		go manager.executeStart()
	} else {
		return manager.executeStart()
	}

	return nil
}

// Stop ...
func (manager *Manager) Stop() error {
	if manager.started {
		log.Infof("stopping...")

		executeAction("stop", manager.processes)
		executeAction("stop", manager.worklist)
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

func (manager *Manager) executeStart() error {
	log.Info("starting...")

	// listen for termination signals
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	if err := executeAction("start", manager.processes); err != nil {
		return err
	}
	if err := executeAction("start", manager.worklist); err != nil {
		return err
	}
	if err := executeAction("start", manager.webs); err != nil {
		return err
	}
	if err := executeAction("start", manager.nsqProducers); err != nil {
		return err
	}
	if err := executeAction("start", manager.nsqConsumers); err != nil {
		return err
	}
	if err := executeAction("start", manager.dbs); err != nil {
		return err
	}
	if err := executeAction("start", manager.redis); err != nil {
		return err
	}

	manager.started = true
	log.Infof("started")

	select {
	case <-termChan:
		log.Infof("received term signal")
	case <-manager.quit:
		log.Infof("received shutdown signal")
	}

	return manager.Stop()
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
