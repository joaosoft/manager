package gomanager

import (
	"sync"
)

// Mode ...
type Mode int

const (
	// First In First Out
	FIFO Mode = iota
	// Last In Last Out
	LIFO
)

type node struct {
	data     interface{}
	next     *node
	previous *node
}

// Queue ...
type Queue struct {
	length int
	start  *node
	end    *node
	mode   Mode
	mux    *sync.Mutex
}

// NewQueue ...
func NewQueue(mode Mode) *Queue {
	return &Queue{
		mode: mode,
	}
}

// Push ...
func (queue *Queue) Push(data interface{}) {
	queue.mux.Lock()
	defer queue.mux.Unlock()

	if queue.length == 0 {
		new := &node{data: data, previous: queue.end}
		queue.start = new
		queue.end = new
	} else {
		new := &node{data: data}
		queue.end.next = new
		queue.end = new
	}
	queue.length++
}

// Pop ...
func (queue *Queue) Pop() interface{} {
	queue.mux.Lock()
	defer queue.mux.Unlock()

	if queue.length == 0 {
		log.Error("the queue is empty")
		return nil
	}

	var toRemove *node
	switch queue.mode {
	case FIFO:
		toRemove = queue.end
		queue.end = toRemove.previous
		return toRemove.data

	case LIFO:
		toRemove = queue.start
		queue.start = queue.start.next
		return toRemove.data

	default:
		return nil
	}
	queue.length--

	return nil
}
