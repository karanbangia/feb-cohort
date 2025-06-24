package main

import (
	"fmt"
	"sync"
)

var count = 0

func main() {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			count++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Print(count)
}
