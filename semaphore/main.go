package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/sync/singleflight"
)

type User struct {
	ID   int
	Name string
}

type Service struct {
	db        *sql.DB
	group     singleflight.Group
	cache     sync.Map
	cacheHits int64
	cacheMiss int64

	dbHits     int64
	sharedHits int64
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) GetUsersWithoutSingleFlight(key string) ([]User, error) {
	// 1. Check cache
	if value, ok := s.cache.Load(key); ok {
		atomic.AddInt64(&s.cacheHits, 1)
		return value.([]User), nil
	}

	atomic.AddInt64(&s.cacheMiss, 1)

	// 2. No singleflight here.
	// So many goroutines can enter this DB call at same time.
	users, err := s.getUsersFromDB()
	if err != nil {
		return nil, err
	}
	fmt.Println(users)
	// 3. Store in cache
	s.cache.Store(key, users)

	return users, nil
}

func main() {
	dsn := "replication_user:pass@tcp(127.0.0.1:3308)/mydatabase"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(5 * time.Minute)

	service := NewService(db)

	requests := 1000
	key := "users:all"

	start := time.Now()

	var wg sync.WaitGroup

	for i := 0; i < requests; i++ {
		wg.Add(1)

		go func(requestID int) {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)

			//users, err := service.GetUsersWithSingleFlight(key)
			users, err := service.GetUsersWithoutSingleFlight(key)

			if err != nil {
				log.Printf("request %d failed: %v\n", requestID, err)
				return
			}

			if requestID < 5 {
				log.Printf("request %d got %d users\n", requestID, len(users))
			}
		}(i)
	}

	wg.Wait()

	fmt.Println("------------- RESULT -------------")
	fmt.Println("Total requests:", requests)
	fmt.Println("Cache hits:", atomic.LoadInt64(&service.cacheHits))
	fmt.Println("Cache misses:", atomic.LoadInt64(&service.cacheMiss))
	fmt.Println("Shared Hits:", atomic.LoadInt64(&service.sharedHits))
	fmt.Println("Actual DB hits:", atomic.LoadInt64(&service.dbHits))
	fmt.Println("Total time:", time.Since(start))
}

func (s *Service) GetUsersWithSingleFlight(key string) ([]User, error) {
	// 1. Check cache first
	if value, ok := s.cache.Load(key); ok {
		atomic.AddInt64(&s.cacheHits, 1)
		return value.([]User), nil
	}

	atomic.AddInt64(&s.cacheMiss, 1)

	// 2. singleflight: only one goroutine executes DB query for same key
	result, err, shared := s.group.Do(key, func() (interface{}, error) {
		// Double-check cache after entering singleflight
		if value, ok := s.cache.Load(key); ok {
			atomic.AddInt64(&s.cacheHits, 1)

			return value.([]User), nil
		}

		users, err := s.getUsersFromDB()
		if err != nil {
			return nil, err
		}

		s.cache.Store(key, users)

		return users, nil
	})

	if err != nil {
		return nil, err
	}

	if shared {
		atomic.AddInt64(&s.sharedHits, 1)

		// This means result was shared with other waiting goroutines
		// Useful for observing singleflight behavior
	}

	return result.([]User), nil
}

func (s *Service) getUsersFromDB() ([]User, error) {
	atomic.AddInt64(&s.dbHits, 1)

	rows, err := s.db.Query("SELECT id, name FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
