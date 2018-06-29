package manager

import "time"

// Work ...
type Work struct {
	id          string
	data        interface{}
	retries     int
	createdAt   time.Time
	elapsedTime time.Time
	endedAt     time.Time
}

// NewWork ...
func NewWork(id string, data interface{}) *Work {
	return &Work{
		id:        id,
		data:      data,
		createdAt: time.Now(),
	}
}

// ElapsedTime ...
func (work *Work) ElapsedTime() time.Duration {
	return time.Since(work.createdAt)
}
