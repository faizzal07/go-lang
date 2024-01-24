package main

import (
	"barbershop/handler"
	"fmt"
	"time"
)

const (
	closingTime = 10 * time.Second
)

func main() {
	shop := handler.NewBarberShop()

	go shop.RunBarbers()
	go shop.GenerateCustomers()

	time.Sleep(closingTime)
	shop.CloseShop()

	shop.Wg.Wait()
	fmt.Println("Barber shop is closed.")
}
