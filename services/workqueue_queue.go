package gomanager

import (
	"encoding/json"
	"fmt"
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
	id       string      `json:"id"`
	data     interface{} `json:"data"`
	next     *node       `json:"next"`
	previous *node       `json:"previous"`
}

// Queue ...
type Queue struct {
	size    int              `json:"size"`
	start   *node            `json:"start"`
	end     *node            `json:"end"`
	mode    Mode             `json:"mode"`
	maxSize int              `json:"max_size"`
	mux     *sync.Mutex      `json:"-"`
	ids     map[string]*node `json:"-"`
}

// NewQueue ...
func NewQueue(options ...QueueOption) *Queue {
	queue := &Queue{
		ids: make(map[string]*node),
		mux: &sync.Mutex{},
	}
	queue.Reconfigure(options...)

	return queue
}

// Add ...
func (queue *Queue) Add(id string, data interface{}) error {
	queue.mux.Lock()
	defer queue.mux.Unlock()
	fmt.Printf("A ADICIONAR\n")

	if queue.maxSize > 0 && queue.size >= queue.maxSize {
		return fmt.Errorf("the queue is full with [ size: %d ]", queue.size)
	}

	nodeToAdd := &node{id: id, data: data, previous: queue.end}
	if queue.size == 0 {
		queue.start = nodeToAdd
		queue.end = nodeToAdd
	} else {
		queue.end.next = queue.end
		queue.end = nodeToAdd
	}
	queue.ids[id] = nodeToAdd
	queue.size++
	return nil
}

// Remove ...
func (queue *Queue) Remove(ids ...string) interface{} {
	queue.mux.Lock()
	defer queue.mux.Unlock()
	fmt.Printf("A REMOVER\n")
	fmt.Println(queue.Dump())

	if queue.size == 0 {
		log.Error("the list is empty")
		return nil
	}
	var nodeToRemove *node
	if len(ids) == 0 {
		switch queue.mode {
		case FIFO:
			nodeToRemove = queue.end
			if queue.size > 1 {
				queue.end = nodeToRemove.previous
				queue.end.next = nil
			} else {
				queue.start = nil
				queue.end = nil
			}
			delete(queue.ids, nodeToRemove.id)
			queue.size--
			return nodeToRemove.data

		case LIFO:
			nodeToRemove = queue.start
			if queue.size > 1 {
				queue.start = queue.start.next
			} else {
				queue.start.next = nil
			}
			delete(queue.ids, nodeToRemove.id)
			queue.size--
			return nodeToRemove.data

		default:
			return nil
		}
	} else {
		var nodesRemoved []interface{}
		for _, id := range ids {
			nodeToRemove = queue.ids[id]
			nodeToRemove.previous.next = nodeToRemove.next
			delete(queue.ids, nodeToRemove.id)
			nodesRemoved = append(nodesRemoved, nodeToRemove.data)
			queue.size--
		}
		return nodesRemoved
	}

	return nil
}

// Size ...
func (queue *Queue) Size() int {
	queue.mux.Lock()
	defer queue.mux.Unlock()
	return queue.size
}

// String ...
func (queue *Queue) Dump() string {
	type queuePrint struct {
		Size    int              `json:"size"`
		Mode    Mode             `json:"mode"`
		MaxSize int              `json:"max_size"`
		Ids     map[string]*node `json:"ids"`
	}

	print := queuePrint{
		Size:    queue.size,
		Mode:    queue.mode,
		MaxSize: queue.maxSize,
		Ids:     queue.ids,
	}

	if json, err := json.Marshal(print); err != nil {
		log.Error(err)
		return ""
	} else {
		return string(json)
	}
}
