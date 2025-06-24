package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	start := time.Now()

	ch := generateNumbers()
	chs := fanout(ch, 8)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ch := range chs[i] {
				checkPrime(ch)
			}
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Print("Total primes: ", totalPrimes, " Time taken: ", elapsed)
}

func generateNumbers() <-chan int32 {
	ch := make(chan int32)
	go func() {
		defer close(ch)
		for i := 3; i < 100000000; i++ {
			ch <- int32(i)
		}
	}()
	return ch
}

var totalPrimes int32 = 0

func checkPrime(n int32) bool {
	if n&1 == 0 {
		return false
	}
	for i := int32(3); i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	atomic.AddInt32(&totalPrimes, 1)
	return true
}

func fanout(channel <-chan int32, breakup int32) []chan int32 {
	outputs := make([]chan int32, breakup)
	for i := range outputs {
		outputs[i] = make(chan int32, 10000)
	}
	go func() {
		defer func() {
			for i := int32(0); i < breakup; i++ {
				close(outputs[i])
			}
		}()
		for ch := range channel {
			outputs[rand.Intn(int(breakup))] <- ch
		}
	}()
	return outputs
}
