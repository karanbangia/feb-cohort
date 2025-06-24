package main

import (
	"database/sql"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID   int
	Name string
}

func main() {
	// get number of threads which can be spawned
	numThreads := runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("Number of threads that can be spawned: %d\n", numThreads)

	dsn := "dbeaver:dbeaver@tcp(127.0.0.1:3306)/mysql"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	// airline seat reset
	seatReset(db)

	// get all users
	users := getUsers(db)

	wg := sync.WaitGroup{}

	// users book seats
	startTime := time.Now()
	for _, user := range users {
		wg.Add(1)
		user := user
		go func() {
			defer wg.Done()
			seatName, err := bookings(db, user)
			if err != nil {
				log.Printf("User %s could not book a seat\n", user.Name)
			}
			log.Printf("User %s with id:%d booked seat %s\n", user.Name, user.ID, seatName)
		}()
	}
	wg.Wait()
	// how many seats booked
	fmt.Printf("Seats booked: %d\n", seatsBooked(db))
	fmt.Printf("Time taken: %v\n", time.Since(startTime))
}

func seatsBooked(db *sql.DB) int {
	// get the number of seats booked
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM seats WHERE user_id IS NOT NULL").Scan(&count)
	if err != nil {
		panic(err)
	}
	return count
}

func bookings(db *sql.DB, user User) (string, error) {
	// start a transaction
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	// get a seat
	var seatName string
	var seatId int
	err = tx.QueryRow("SELECT id, name FROM seats WHERE user_id IS NULL ORDER BY id LIMIT 1 FOR UPDATE").Scan(&seatId, &seatName)

	// update the seat with the user id
	_, err = tx.Exec("UPDATE seats SET user_id = ? WHERE id = ?", user.ID, seatId)
	if err != nil {
		return "", err
	}
	err = tx.Commit()
	if err != nil {
		return "", err
	}
	return seatName, nil
}

func getUsers(db *sql.DB) []User {
	// get all users from the database
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	// create a slice of users
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	return users
}

func seatReset(db *sql.DB) {
	// update user id to null
	_, err := db.Exec("UPDATE seats SET user_id = NULL")
	if err != nil {
		panic(err)
	}
}
