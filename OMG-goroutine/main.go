package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Order struct {
	ID     int
	Status string
}

func main() {
	orders := generateOrders(10)

	go processOrders(orders)
	go updateOrderStatus(orders)

	fmt.Println("All operations completetd. Exiting")

}

func generateOrders(count int) []*Order {
	orders := make([]*Order, count)
	for i := 0; i < count; i++ {
		orders[i] = &Order{ID: i + 1, Status: "pending"}
	}
	return orders

}

func processOrders(orders []*Order) {
	for _, order := range orders {
		delay := time.Duration(rand.Intn(500)) * time.Millisecond
		time.Sleep(delay)
		fmt.Printf("processing order %d taken: %dms \n", order.ID, delay)
	}
}

func updateOrderStatus(orders []*Order) {

	for _, order := range orders {
		time.Sleep(
			time.Duration(rand.Intn(300)) * time.Millisecond,
		)
		status := []string{
			"processing", "shipped", "delivered",
		}[rand.Intn(3)]
		order.Status = status
		fmt.Printf("updated order %d status: %s\n", order.ID, status)
	}
}
