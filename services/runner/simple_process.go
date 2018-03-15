package gomanager

import (
	"fmt"

	"go-manager/services/nsq"
	"go-manager/services/process"

	"github.com/labstack/gommon/log"
)

// -------------- PROCESS CLIENTS --------------
// NewNSQConsumer ... creates a new nsq consumer
func (manager *Manager) NewNSQConsumer(config *nsq.Config, handler nsq.IHandler) (nsq.IConsumer, error) {
	return nsq.NewConsumer(config, handler)
}

// NewNSQConsumer ... creates a new nsq producer
func (manager *Manager) NewNSQProducer(config *nsq.Config) (nsq.IProducer, error) {
	return nsq.NewProducer(config)
}

// -------------- METHODS --------------
// GetProcess ... get a process with key
func (manager *Manager) GetProcess(key string) process.IProcess {
	return manager.ProcessController[key].Process
}

// AddProcess ... add a process with key
func (manager *Manager) AddProcess(key string, prc process.IProcess) error {
	if manager.Started {
		panic("Manager, can not add processes after start")
	}

	manager.ProcessController[key] = &process.ProcessController{
		Process: prc,
		Control: make(chan int),
	}
	log.Infof(fmt.Sprintf("Manager, process '%s' added", key))

	return nil
}

// RemProcess ... remove the process by bey
func (manager *Manager) RemProcess(key string) (process.IProcess, error) {
	// get process
	controller := manager.ProcessController[key]

	// delete process
	delete(manager.ProcessController, key)
	log.Infof(fmt.Sprintf("Manager, process '%s' removed", key))

	return controller.Process, nil
}

// launch ... starts a process
func (manager *Manager) launch(name string, controller *process.ProcessController) error {
	if err := controller.Process.Start(); err != nil {
		log.Error(err, fmt.Sprintf("Manager, error launching process [process:%s]", name))
		manager.Stop()
		controller.Control <- 0
	} else {
		log.Infof(fmt.Sprintf("Manager, launched process [process:%s]", name))
		controller.Started = true
		controller.Control <- 0
	}

	return nil
}
