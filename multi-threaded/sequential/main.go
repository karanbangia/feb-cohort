package main

import (
	"fmt"
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
	totalPrimes++
	return true
}

func main() {

	start := time.Now()

	for i := 2; i < 100; i++ {
		checkPrime(i)
	}
	elapsed := time.Since(start)
	fmt.Print("Total primes: ", totalPrimes, " Time taken: ", elapsed)
}
