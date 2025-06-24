package main

import (
	"fmt"
	"sync"
	"time"
)

type Jobs func()

type ThreadPool struct {
	Worker chan Jobs
	wg     sync.WaitGroup
}

func NewThreadPool(size int) *ThreadPool {
	pool := &ThreadPool{
		Worker: make(chan Jobs),
	}
	for i := 0; i < size; i++ {
		pool.wg.Add(1)
		go func() {
			defer pool.wg.Done()
			for job := range pool.Worker {
				job()
			}
		}()
	}
	return pool
}

func (tp *ThreadPool) AddJob(job Jobs) {
	tp.Worker <- job
}

func (tp *ThreadPool) Wait() {
	close(tp.Worker)
	tp.wg.Wait()
}

func main() {

	pool := NewThreadPool(3)
	for i := 0; i < 10; i++ {

		pool.AddJob(func() {
			time.Sleep(1000 * time.Millisecond)
		})
		fmt.Println("Job added", i)
	}
	pool.Wait()
}
