package mgr

import (
	"database/sql"
	"fmt"
	"github.com/joaosoft/go-Manager/config"
	"github.com/joaosoft/go-Manager/elastic"
	"github.com/joaosoft/go-Manager/gateway"
	"github.com/joaosoft/go-Manager/nsq"
	"github.com/joaosoft/go-Manager/process"
	"github.com/joaosoft/go-Manager/sqlcon"
	"github.com/joaosoft/go-Manager/web"
	"github.com/joaosoft/go-Manager/workqueue"
	"github.com/labstack/gommon/log"
	"io"
	"os"
	"os/signal"
	"syscall"
)

// Manager ... Manager structure
type Manager struct {
	processController     map[string]*process.ProcessController
	configController      map[string]*config.ConfigController
	sqlConController      map[string]*sqlcon.SQLConController
	gatewayController     map[string]*gateway.Gateway
	elasticController     map[string]*elastic.ElasticController
	workerQueueController map[string]*workqueue.QueueController

	control chan int
	Started bool
}

// NewManager ... create a new Manager
func NewManager() (*Manager, error) {

	return &Manager{
		processController:     make(map[string]*process.ProcessController),
		configController:      make(map[string]*config.ConfigController),
		sqlConController:      make(map[string]*sqlcon.SQLConController),
		gatewayController:     make(map[string]*gateway.Gateway),
		elasticController:     make(map[string]*elastic.ElasticController),
		workerQueueController: make(map[string]*workqueue.QueueController),

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
	for name, process := range instance.processController {
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

		for key, controller := range instance.processController {
			if controller.Started {
				log.Infof("Manager, stopping process [process:%s]", key)
				if err := controller.Process.Stop(); err != nil {
					log.Error(err, fmt.Sprintf("error stopping process [process:%s]", key))
				}
				log.Infof("Manager, close channel [process:%s]", key)
				<-controller.Control
				close(controller.Control)
				delete(instance.processController, key)
				log.Infof("Manager, stopped process [process:%s]", key)
			}
		}

		instance.Started = false
		log.Infof("Manager, stopped")
	}

	return nil
}
