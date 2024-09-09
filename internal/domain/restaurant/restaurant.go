package restaurant

import (
	"booking-dinner/internal/domain/models"
)

// Service defines the interface for restaurant operations
type Service interface {
	InitializeTables(numTables int) error
	ReserveTables(numCustomers int) (string, int, int, error)
	CancelReservation(bookingID string) (int, int, error)
	GetAvailableTables() int
}

// Repository defines the interface for data storage operations
type Repository interface {
	InitializeTables(numTables int) error
	ReserveTables(booking models.Booking) error
	CancelReservation(bookingID string) (int, error)
	GetAvailableTables() int
	IsInitialized() bool
}
