package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

var totalPrimes int32 = 0

func checkPrime(n int) bool {
	if n&1 == 0 {
		return false
	}
	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	atomic.AddInt32(&totalPrimes, 1)
	fmt.Println("Prime: ", n)
	return true
}

func main() {
	start := time.Now()
	wg := sync.WaitGroup{}
	for i := 2; i < math.MaxInt; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			checkPrime(i)
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Print("Total primes: ", totalPrimes, " Time taken: ", elapsed)
}
