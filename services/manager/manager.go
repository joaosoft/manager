package gomanager

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

// Manager ...
type Manager struct {
	processes map[string]IProcessManager
	configs   map[string]ISetupManager
	sqls      map[string]*sql.DB
	webs      map[string]*echo.Echo
	gateways  map[string]*http.Client
	workers   map[string]IWorkManager

	control chan int
	started bool
}

// NewManager ...
func NewManager() (*Manager, error) {
	return &Manager{
		processes: make(map[string]IProcessManager),
		configs:   make(map[string]ISetupManager),
		sqls:      make(map[string]*sql.DB),
		webs:      make(map[string]*echo.Echo),
		gateways:  make(map[string]*http.Client),
		workers:   make(map[string]IWorkManager),
		control:   make(chan int),
	}, nil
}

// Started ...
func (manager *Manager) Started() bool {
	return manager.started
}

// Start ...
func (manager *Manager) Start() error {
	log.Infof("Manager, starting...")
	manager.started = true

	// listen for termination signals
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	// launch every process in a separeted process
	for name, process := range manager.processes {
		log.Infof("Manager, starting process [process:%s]", name)

		go manager.launch(name, process)

		log.Infof("Manager, started process [process:%s]", name)
	}

	select {
	case <-termChan:
		log.Infof("Manager, received term signal")
	case <-manager.control:
		log.Infof("Manager, received shutdown signal")
	}

	manager.Stop()

	return nil
}

// Stop ...
func (manager *Manager) Stop() error {
	if manager.started {
		log.Infof("Manager, stopping")

		for key, process := range manager.processes {
			if process.Started() {
				log.Infof("Manager, stopping process [ process: %s ]", key)
				if err := process.Stop(); err != nil {
					log.Error(err, fmt.Sprintf("error stopping process [process:%s]", key))
				}
				log.Infof("Manager, close channel [ process: %s ]", key)
				delete(manager.processes, key)
				log.Infof("Manager, stopped process [ process: %s ]", key)
			}
		}

		manager.started = true
		log.Infof("Manager, stopped")
	}

	return nil
}
