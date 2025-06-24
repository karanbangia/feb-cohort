package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-faker/faker/v4"
)

func main() {
	username := "dbeaver"
	password := "dbeaver"
	host := "127.0.0.1"
	port := "3306"

	// Connect without specifying a database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/DB_1", username, password, host, port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	//for i := 0; i < 40; i++ {
	//	// insert into seats table with seat number 1A 1B 1C 1D 1E 1F do not use faker where 1 should be variable
	//	_, err = db.Exec("INSERT INTO seats (name) VALUES (?)", fmt.Sprintf("%dA", i+1))
	//	if err != nil {
	//		return
	//	}
	//	db.Exec("INSERT INTO seats (name) VALUES (?)", fmt.Sprintf("%dB", i+1))
	//	db.Exec("INSERT INTO seats (name) VALUES (?)", fmt.Sprintf("%dC", i+1))
	//	db.Exec("INSERT INTO seats (name) VALUES (?)", fmt.Sprintf("%dD", i+1))
	//
	//}

	fakeNames(err, db)
}

func fakeNames(err error, db *sql.DB) {
	for i := 0; i < 100; i++ {
		name := faker.FirstName() + " " + faker.LastName()
		// insert this in sql
		_, err = db.Exec("INSERT INTO USER(id,name) VALUES (?,?)", i, name)
		if err != nil {
			panic(err)
		}

	}
}
