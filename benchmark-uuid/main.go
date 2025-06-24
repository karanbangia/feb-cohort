package main

import (
	"database/sql"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/google/uuid"
)

const (
	maxOpenConn = 100
)

func main() {
	db := dbConnection()
	AddUuid(db)
	AddInteger(db)
}

func AddUuid(db *sql.DB) {
	wg := &sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := db.Exec("INSERT INTO UUID_BENCHMARK (id) VALUES (?)", uuid.New().String())
			if err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()
}

func AddInteger(db *sql.DB) {
	wg := &sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := db.Exec("INSERT INTO NUMBER_BENCHMARK (stub) VALUES ('a')")
			if err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()
}

func dbConnection() *sql.DB {
	dsn := "dbeaver:dbeaver@tcp(127.0.0.1:3306)/mysql"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db
}
