package mgr

import (
	"fmt"

	"github.com/joaosoft/go-manager/services/nsq"
	"github.com/joaosoft/go-manager/services/process"
	"github.com/labstack/gommon/log"
)

// -------------- PROCESS CLIENTS --------------
// NewNSQConsumer ... creates a new nsq consumer
func (instance *Manager) NewNSQConsumer(config *nsq.Config, handler nsq.IHandler) (nsq.IConsumer, error) {
	return nsq.NewConsumer(config, handler)
}

// NewNSQConsumer ... creates a new nsq producer
func (instance *Manager) NewNSQProducer(config *nsq.Config) (nsq.IProducer, error) {
	return nsq.NewProducer(config)
}

// -------------- METHODS --------------
// GetProcess ... get a process with key
func (instance *Manager) GetProcess(key string) process.IProcess {
	return instance.ProcessController[key].Process
}

// AddProcess ... add a process with key
func (instance *Manager) AddProcess(key string, prc process.IProcess) error {
	if instance.Started {
		panic("Manager, can not add processes after start")
	}

	instance.ProcessController[key] = &process.ProcessController{
		Process: prc,
		Control: make(chan int),
	}
	log.Infof(fmt.Sprintf("Manager, process '%s' added", key))

	return nil
}

// RemProcess ... remove the process by bey
func (instance *Manager) RemProcess(key string) (process.IProcess, error) {
	// get process
	controller := instance.ProcessController[key]

	// delete process
	delete(instance.ProcessController, key)
	log.Infof(fmt.Sprintf("Manager, process '%s' removed", key))

	return controller.Process, nil
}

// launch ... starts a process
func (instance *Manager) launch(name string, controller *process.ProcessController) error {
	if err := controller.Process.Start(); err != nil {
		log.Error(err, fmt.Sprintf("Manager, error launching process [process:%s]", name))
		instance.Stop()
		controller.Control <- 0
	} else {
		log.Infof(fmt.Sprintf("Manager, launched process [process:%s]", name))
		controller.Started = true
		controller.Control <- 0
	}

	return nil
}
