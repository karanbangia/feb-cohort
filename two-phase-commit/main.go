package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/feb-cohort/two-phase-commit/delivery"
	"github.com/feb-cohort/two-phase-commit/store"
)

type OrderCoordinator struct {
	storeService    *store.StoreService
	deliveryService *delivery.DeliveryService
}

func NewOrderCoordinator() (*OrderCoordinator, error) {
	storeSvc, err := store.NewStoreService()
	if err != nil {
		return nil, fmt.Errorf("failed to create store service: %v", err)
	}

	deliverySvc, err := delivery.NewDeliveryService()
	if err != nil {
		return nil, fmt.Errorf("failed to create delivery service: %v", err)
	}

	return &OrderCoordinator{
		storeService:    storeSvc,
		deliveryService: deliverySvc,
	}, nil
}

func (c *OrderCoordinator) PlaceOrder(orderID int, items []string) error {
	// Phase 1: Prepare
	if err := c.storeService.PrepareOrder(orderID, items); err != nil {
		return fmt.Errorf("store prepare failed: %v", err)
	}

	if err := c.deliveryService.PrepareDelivery(orderID); err != nil {
		// Rollback store
		c.storeService.RollbackOrder(orderID)
		return fmt.Errorf("delivery prepare failed: %v", err)
	}

	// Phase 2: Commit
	if err := c.storeService.CommitOrder(orderID); err != nil {
		// Rollback both
		c.storeService.RollbackOrder(orderID)
		c.deliveryService.RollbackDelivery(orderID)
		return fmt.Errorf("store commit failed: %v", err)
	}

	if err := c.deliveryService.CommitDelivery(orderID); err != nil {
		// Rollback both
		c.storeService.RollbackOrder(orderID)
		c.deliveryService.RollbackDelivery(orderID)
		return fmt.Errorf("delivery commit failed: %v", err)
	}

	return nil
}

func main() {
	coordinator, err := NewOrderCoordinator()
	if err != nil {
		log.Fatalf("Failed to create coordinator: %v", err)
	}

	var wg sync.WaitGroup
	// Simulate multiple orders
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(orderID int) {
			defer wg.Done()
			items := []string{"Pizza", "Burger", "Fries"}
			if err := coordinator.PlaceOrder(orderID, items); err != nil {
				log.Printf("Order %d failed: %v", orderID, err)
			} else {
				log.Printf("Order %d placed successfully", orderID)
			}
		}(i)
		time.Sleep(100 * time.Millisecond) // Small delay between orders
	}

	wg.Wait()
}
