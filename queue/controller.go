package queue

import (
	"github.com/labstack/gommon/log"
)

// IQueueController queue controller interface
type IQueueController interface {
	Do(data interface{}) error
	Undo() error
}

// QueueController queue controller structure
type QueueController struct {
	Repository Repository
	Data       interface{}
}

// NewQueueController create a new queue controller
func NewQueueController(repository Repository) *QueueController {
	return &QueueController{
		Repository: repository,
	}
}

// Do ... do something
func (instance *QueueController) Do(data interface{}) error {

	return instance.Repository.DoSomething(data)
}

// Undo .. undo something
func (controller *QueueController) Undo() error {
	log.Infof("Undo()")

	return nil
}
