package gomanager

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
