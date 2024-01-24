package handler

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
	Id int
}

type BarberShop struct {
	mu                sync.Mutex
	customers         chan *Customer
	waitingRoom       chan *Customer
	barbersAvailable  int
	barbersBusy       int
	closingTimeSignal chan bool
	Wg                sync.WaitGroup
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

func (bs *BarberShop) RunBarbers() {
	for {
		select {
		case customer := <-bs.customers:
			bs.mu.Lock()
			if bs.barbersAvailable > 0 {
				bs.barbersAvailable--
				bs.barbersBusy++
				bs.mu.Unlock()
				bs.Wg.Add(1)
				go bs.barberHandler(customer)
			} else if len(bs.waitingRoom) < numWaitingChairs {
				bs.mu.Unlock()
				bs.handleOverflow(customer)
			} else {
				bs.mu.Unlock()
				fmt.Printf("Customer %d left because the waiting room is full\n", customer.Id)
			}
		case <-bs.closingTimeSignal:
			close(bs.customers)
			return
		}
	}
}

func (bs *BarberShop) barberHandler(customer *Customer) {
	defer bs.Wg.Done()
	time.Sleep(time.Duration(4) * time.Second)
	fmt.Printf("Barber %d finished cutting hair for Customer %d\n", rand.Intn(numBarbers), customer.Id)

	bs.mu.Lock()
	bs.barbersAvailable++
	bs.barbersBusy--

	select {
	case waitingCustomer := <-bs.waitingRoom:
		bs.barbersAvailable--
		bs.barbersBusy++
		bs.mu.Unlock()
		bs.Wg.Add(1)
		go bs.barberHandler(waitingCustomer)
	default:
		bs.mu.Unlock()
	}
}

func (bs *BarberShop) GenerateCustomers() {
	for i := 1; i < 10; i++ {
		select {
		case bs.customers <- &Customer{Id: i}:
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
		bs.Wg.Add(1)
		go bs.barberHandler(customer)
	} else if len(bs.waitingRoom) < numWaitingChairs {
		bs.waitingRoom <- customer
		fmt.Printf("Customer %d is waiting in the waiting room\n", customer.Id)
	} else {
		fmt.Printf("Customer %d left because the waiting room is full\n", customer.Id)
	}
}

func (bs *BarberShop) CloseShop() {
	fmt.Println("Barber shop is closing.")
	close(bs.closingTimeSignal)
}
