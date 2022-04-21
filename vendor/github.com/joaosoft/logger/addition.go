package logger

import (
	"github.com/joaosoft/errors"
)

type Addition struct {
	message string
}

// NewAddition ...
func NewAddition(message string) IAddition {
	addition := &Addition{
		message: message,
	}

	return addition
}

// ToError
func (addition *Addition) ToError() error {
	return errors.New(errors.ErrorLevel, 0, addition.message)
}
