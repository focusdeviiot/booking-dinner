package memory

import (
	"errors"
	"sync"

	"booking-dinner/internal/domain/models"
)

// RestaurantRepository represents an in-memory storage for restaurant data
type RestaurantRepository struct {
	tables        int
	bookings      map[string]models.Booking
	mutex         sync.RWMutex
	isInitialized bool
}

// NewRestaurantRepository creates a new instance of RestaurantRepository
func NewRestaurantRepository() *RestaurantRepository {
	return &RestaurantRepository{
		bookings: make(map[string]models.Booking),
	}
}

// InitializeTables sets the initial number of tables in the restaurant
func (r *RestaurantRepository) InitializeTables(tables int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.isInitialized {
		return errors.New("tables have already been initialized")
	}

	r.tables = tables
	r.isInitialized = true
	return nil
}

// ReserveTables reserves tables for a booking
func (r *RestaurantRepository) ReserveTables(booking models.Booking) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.isInitialized {
		return errors.New("tables have not been initialized")
	}

	if r.tables < booking.TablesBooked {
		return errors.New("not enough tables available")
	}

	r.bookings[booking.ID] = booking
	r.tables -= booking.TablesBooked
	return nil
}

// CancelReservation cancels a booking and frees up the tables
func (r *RestaurantRepository) CancelReservation(bookingID string) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	booking, exists := r.bookings[bookingID]
	if !exists {
		return 0, errors.New("booking not found")
	}

	delete(r.bookings, bookingID)
	r.tables += booking.TablesBooked
	return booking.TablesBooked, nil
}

// GetAvailableTables returns the number of available tables
func (r *RestaurantRepository) GetAvailableTables() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.tables
}

// IsInitialized checks if the tables have been initialized
func (r *RestaurantRepository) IsInitialized() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.isInitialized
}
