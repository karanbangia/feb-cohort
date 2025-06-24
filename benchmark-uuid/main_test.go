package main

import "testing"

// benchmark AddInteger
func BenchmarkAddInteger(b *testing.B) {
	db := dbConnection()
	b.ReportAllocs() // Reports memory allocations
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			AddInteger(db)
		}
	})
}

func BenchmarkAddUuid(b *testing.B) {
	db := dbConnection()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			AddUuid(db)
		}
	})
}
