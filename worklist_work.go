package manager

import (
	"time"

	"github.com/joaosoft/logger"
)

// Work ...
type Work struct {
	id          string
	data        interface{}
	retries     int
	createdAt   time.Time
	elapsedTime time.Time
	endedAt     time.Time
	logger      logger.ILogger
}

// NewWork ...
func NewWork(id string, data interface{}, logger logger.ILogger) *Work {
	return &Work{
		id:        id,
		data:      data,
		createdAt: time.Now(),
		logger:    logger,
	}
}

// ElapsedTime ...
func (work *Work) ElapsedTime() time.Duration {
	return time.Since(work.createdAt)
}
