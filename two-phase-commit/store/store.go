package store

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type StoreService struct {
	db *sql.DB
}

type Order struct {
	OrderID      int
	RestaurantID int
	Items        []string
	Status       string // "PREPARING", "COMMITTED", "ROLLED_BACK"
}

func NewStoreService() (*StoreService, error) {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/zomato_store?sslmode=disable")
	if err != nil {
		return nil, err
	}

	// Create tables if they don't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			order_id SERIAL PRIMARY KEY,
			restaurant_id INT NOT NULL,
			items TEXT[] NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'PREPARING'
		)
	`)
	if err != nil {
		return nil, err
	}

	return &StoreService{db: db}, nil
}

func (s *StoreService) PrepareOrder(orderID int, items []string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO orders (order_id, items, status) VALUES ($1, $2, 'PREPARING')",
		orderID, items)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *StoreService) CommitOrder(orderID int) error {
	_, err := s.db.Exec("UPDATE orders SET status = 'COMMITTED' WHERE order_id = $1", orderID)
	return err
}

func (s *StoreService) RollbackOrder(orderID int) error {
	_, err := s.db.Exec("UPDATE orders SET status = 'ROLLED_BACK' WHERE order_id = $1", orderID)
	return err
}

func (s *StoreService) GetOrderStatus(orderID int) (string, error) {
	var status string
	err := s.db.QueryRow("SELECT status FROM orders WHERE order_id = $1", orderID).Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}
