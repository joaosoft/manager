package writers

// QueueOption ...
type QueueOption func(queue *Queue)

// Reconfigure ...
func (queue *Queue) Reconfigure(options ...QueueOption) {
	for _, option := range options {
		option(queue)
	}
}

// WithMode ...
func WithMode(mode Mode) QueueOption {
	return func(queue *Queue) {
		queue.mode = mode
	}
}

// WithMaxSize ...
func WithMaxSize(size int) QueueOption {
	return func(queue *Queue) {
		queue.maxSize = size
	}
}
