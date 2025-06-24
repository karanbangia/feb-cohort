package delivery

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DeliveryService struct {
	db *sql.DB
}

type Delivery struct {
	DeliveryID int
	OrderID    int
	Status     string // "PREPARING", "COMMITTED", "ROLLED_BACK"
}

func NewDeliveryService() (*DeliveryService, error) {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/zomato_delivery?sslmode=disable")
	if err != nil {
		return nil, err
	}

	// Create tables if they don't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS deliveries (
			delivery_id SERIAL PRIMARY KEY,
			order_id INT NOT NULL UNIQUE,
			status VARCHAR(20) NOT NULL DEFAULT 'PREPARING'
		)
	`)
	if err != nil {
		return nil, err
	}

	return &DeliveryService{db: db}, nil
}

func (d *DeliveryService) PrepareDelivery(orderID int) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO deliveries (order_id, status) VALUES ($1, 'PREPARING')",
		orderID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (d *DeliveryService) CommitDelivery(orderID int) error {
	_, err := d.db.Exec("UPDATE deliveries SET status = 'COMMITTED' WHERE order_id = $1", orderID)
	return err
}

func (d *DeliveryService) RollbackDelivery(orderID int) error {
	_, err := d.db.Exec("UPDATE deliveries SET status = 'ROLLED_BACK' WHERE order_id = $1", orderID)
	return err
}

func (d *DeliveryService) GetDeliveryStatus(orderID int) (string, error) {
	var status string
	err := d.db.QueryRow("SELECT status FROM deliveries WHERE order_id = $1", orderID).Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}
