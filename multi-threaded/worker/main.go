package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	start := time.Now()

	numWorkers := runtime.NumCPU() // Use all available CPU cores
	ch := generateNumbers(10000)   // Buffered channel
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ch, &wg)
	}

	wg.Wait() // Ensure all workers finish
	elapsed := time.Since(start)
	fmt.Println("Time taken:", elapsed)
}

func generateNumbers(bufferSize int) <-chan int32 {
	ch := make(chan int32, bufferSize)
	go func() {
		defer close(ch)
		for i := 3; i < 100000000; i++ {
			ch <- int32(i)
		}
	}()
	return ch
}

func worker(ch <-chan int32, wg *sync.WaitGroup) {
	defer wg.Done()
	localCount := int32(0)
	for num := range ch {
		if checkPrime(num) {
			localCount++
		}
	}
	fmt.Println(localCount)
}
func checkPrime(n int32) bool {
	if n < 2 || (n > 2 && n%2 == 0) {
		return false
	}
	for i := int32(3); i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}
