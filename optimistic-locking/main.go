package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var value int32 = 0
	wg := sync.WaitGroup{}
	// Simulate concurrent updates
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {

				oldValue := atomic.LoadInt32(&value)
				newValue := oldValue + 1
				if atomic.CompareAndSwapInt32(&value, oldValue, newValue) {
					fmt.Printf("Updated value to %d\n", newValue)
					break
				}
			}
		}()
	}

	wg.Wait()

	fmt.Printf("Final value: %d\n", value)
}
