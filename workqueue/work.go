package workqueue

type IWork interface {
	GetWork() (interface{}, error)
}

type Work struct {
	data interface{}
}

// NewWork ... create a new work
func NewWork(work interface{}) *Work {
	Work := Work{
		data: work,
	}

	return &Work
}

// GetWork ... get a work
func (work *Work) GetWork() (interface{}, error) {
	return work.data, nil
}
