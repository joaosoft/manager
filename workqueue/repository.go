package workqueue

import (
	"github.com/labstack/gommon/log"
)

// IRepository ... repository interface
type IRepository interface {
	DoSomething(data interface{}) error
}

// Repository ... repository structure
type Repository struct {
}

// NewRepository ...  create a new repository
func NewRepository() *Repository {
	repository := Repository{}

	return &repository
}

// DoSomething ... dummy method
func (repository *Repository) DoSomething(data interface{}) error {
	log.Infof("Repository: DoSomething")

	return nil
}
