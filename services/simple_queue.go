package gomanager

type IWork interface {
	Run() error
}

// Queue ...
type Queue struct {
	name string
}

// NewQueue ...
func NewQueue(name string) *Queue {
	queue := Queue{}
	return &queue
}

// AddWork ...
func (queue *Queue) AddWork(work IWork) error {
	log.Infof("Adding work to queue %s", queue.name)

	return nil
}
