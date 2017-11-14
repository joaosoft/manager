package mgr

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joaosoft/go-manager/services/config"
	"github.com/joaosoft/go-manager/services/elastic"
	"github.com/joaosoft/go-manager/services/gateway"
	"github.com/joaosoft/go-manager/services/process"
	"github.com/joaosoft/go-manager/services/sqlcon"
	"github.com/joaosoft/go-manager/services/workqueue"
	"github.com/labstack/gommon/log"
)

// Manager ... Manager structure
type Manager struct {
	ProcessController     map[string]*process.ProcessController
	ConfigController      map[string]*config.ConfigController
	SqlConController      map[string]*sqlcon.SQLConController
	GatewayController     map[string]*gateway.Gateway
	ElasticController     map[string]*elastic.ElasticController
	WorkerQueueController map[string]*workqueue.QueueController

	control chan int
	Started bool
}

// NewManager ... create a new Manager
func NewManager() (*Manager, error) {

	return &Manager{
		ProcessController:     make(map[string]*process.ProcessController),
		ConfigController:      make(map[string]*config.ConfigController),
		SqlConController:      make(map[string]*sqlcon.SQLConController),
		GatewayController:     make(map[string]*gateway.Gateway),
		ElasticController:     make(map[string]*elastic.ElasticController),
		WorkerQueueController: make(map[string]*workqueue.QueueController),

		control: make(chan int),
	}, nil
}

// Start ... starts and blocks until it receives a signal in its control channel or a SIGTERM,
func (instance *Manager) Start() error {
	log.Infof("Manager, starting")
	instance.Started = true

	// listen for termination signals
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	// launch every process in a separeted process
	for name, process := range instance.ProcessController {
		log.Infof("Manager, starting process [process:%s]", name)

		go instance.launch(name, process)

		log.Infof("Manager, started process [process:%s]", name)
	}

	select {
	case <-termChan:
		log.Infof("Manager, received term signal")
	case <-instance.control:
		log.Infof("Manager, received shutdown signal")
	}

	instance.Stop()

	return nil
}

// Stop ... stop all processes and stops the Manager
func (instance *Manager) Stop() error {
	if instance.Started {
		log.Infof("Manager, stopping")

		for key, controller := range instance.ProcessController {
			if controller.Started {
				log.Infof("Manager, stopping process [process:%s]", key)
				if err := controller.Process.Stop(); err != nil {
					log.Error(err, fmt.Sprintf("error stopping process [process:%s]", key))
				}
				log.Infof("Manager, close channel [process:%s]", key)
				<-controller.Control
				close(controller.Control)
				delete(instance.ProcessController, key)
				log.Infof("Manager, stopped process [process:%s]", key)
			}
		}

		instance.Started = false
		log.Infof("Manager, stopped")
	}

	return nil
}
