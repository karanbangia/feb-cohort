package main

import (
	"fmt"
	"sync"
)

type ConcurrentQueue struct {
	queue []int32
	mu    sync.Mutex
}

func NewConcurrentQueue() *ConcurrentQueue {
	return &ConcurrentQueue{
		queue: make([]int32, 0),
	}
}

func (cq *ConcurrentQueue) Enqueue(val int32) {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	cq.queue = append(cq.queue, val)
}

func (cq *ConcurrentQueue) Dequeue() int32 {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	if len(cq.queue) == 0 {
		return -1
	}
	val := cq.queue[0]
	cq.queue = cq.queue[1:]
	return val
}

func (cq *ConcurrentQueue) Size() int {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	return len(cq.queue)

}

func main() {
	q1 := NewConcurrentQueue()
	q1.Enqueue(1)
	q1.Enqueue(2)
	q1.Enqueue(3)
	fmt.Println(q1.Dequeue())
	fmt.Println(q1.Dequeue())
}
