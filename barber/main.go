package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	numBarbers       = 2
	numWaitingChairs = 2
	closingTime      = 10 * time.Second
)

type Customer struct {
	id int
}

type BarberShop struct {
	mu                sync.Mutex
	customers         chan *Customer
	waitingRoom       chan *Customer
	barbersAvailable  int
	barbersBusy       int
	closingTimeSignal chan bool
	wg                sync.WaitGroup
}

func NewBarberShop() *BarberShop {
	return &BarberShop{
		customers:         make(chan *Customer),
		waitingRoom:       make(chan *Customer, numWaitingChairs),
		barbersAvailable:  numBarbers,
		barbersBusy:       0,
		closingTimeSignal: make(chan bool),
	}
}

func (bs *BarberShop) runBarbers() {
	for {
		select {
		case customer := <-bs.customers:
			bs.mu.Lock()
			if bs.barbersAvailable > 0 {
				bs.barbersAvailable--
				bs.barbersBusy++
				bs.mu.Unlock()
				bs.wg.Add(1)
				go bs.barberHandler(customer)
			} else if len(bs.waitingRoom) < numWaitingChairs {
				bs.mu.Unlock()
				bs.handleOverflow(customer)
			} else {
				bs.mu.Unlock()
				fmt.Printf("Customer %d left because the waiting room is full\n", customer.id)
			}
		case <-bs.closingTimeSignal:
			close(bs.customers)
			return
		}
	}
}

func (bs *BarberShop) barberHandler(customer *Customer) {
	defer bs.wg.Done()
	time.Sleep(time.Duration(4) * time.Second)
	fmt.Printf("Barber %d finished cutting hair for Customer %d\n", rand.Intn(numBarbers), customer.id)

	bs.mu.Lock()
	bs.barbersAvailable++
	bs.barbersBusy--

	select {
	case waitingCustomer := <-bs.waitingRoom:
		bs.barbersAvailable--
		bs.barbersBusy++
		bs.mu.Unlock()
		bs.wg.Add(1)
		go bs.barberHandler(waitingCustomer)
	default:
		bs.mu.Unlock()
	}
}

func (bs *BarberShop) generateCustomers() {
	for i := 1; i < 10; i++ {
		select {
		case bs.customers <- &Customer{id: i}:
			time.Sleep(time.Duration(1) * time.Second)
		case <-bs.closingTimeSignal:
			close(bs.customers)
			return
		}
	}
}

func (bs *BarberShop) handleOverflow(customer *Customer) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if bs.barbersBusy < numBarbers {
		bs.barbersBusy++
		bs.barbersAvailable--
		bs.wg.Add(1)
		go bs.barberHandler(customer)
	} else if len(bs.waitingRoom) < numWaitingChairs {
		bs.waitingRoom <- customer
		fmt.Printf("Customer %d is waiting in the waiting room\n", customer.id)
	} else {
		fmt.Printf("Customer %d left because the waiting room is full\n", customer.id)
	}
}

func (bs *BarberShop) closeShop() {
	fmt.Println("Barber shop is closing.")
	close(bs.closingTimeSignal)
}

func main() {
	shop := NewBarberShop()

	go shop.runBarbers()
	go shop.generateCustomers()

	time.Sleep(closingTime)
	shop.closeShop()

	shop.wg.Wait()
	fmt.Println("Barber shop is closed.")
}
