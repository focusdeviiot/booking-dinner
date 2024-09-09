package models

import (
	"time"
)

type Booking struct {
	ID           string
	CustomerName string
	NumCustomers int
	TablesBooked int
	BookingTime  time.Time
}

func NewBooking(id string, customerName string, numCustomers int, tablesBooked int) *Booking {
	return &Booking{
		ID:           id,
		CustomerName: customerName,
		NumCustomers: numCustomers,
		TablesBooked: tablesBooked,
		BookingTime:  time.Now(),
	}
}
