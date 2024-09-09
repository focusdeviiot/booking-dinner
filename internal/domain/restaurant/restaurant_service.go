package restaurant

import (
	"fmt"
	"math"
	"time"

	"booking-dinner/internal/domain/models"
	"booking-dinner/internal/errors"

	"golang.org/x/exp/rand"
)

type service struct {
	repo          Repository
	seatsPerTable int
	maxTables     int
	charsetCode   string
	lengthCode    int
}

// NewService creates a new instance of restaurant service
func NewService(repo Repository, seatsPerTable int, maxTables int, charsetCode string, lengthCode int) Service {
	return &service{
		repo:          repo,
		seatsPerTable: seatsPerTable,
		maxTables:     maxTables,
		charsetCode:   charsetCode,
		lengthCode:    lengthCode,
	}
}

func (s *service) InitializeTables(numTables int) error {
	if numTables <= 0 || numTables > s.maxTables {
		return errors.NewInitializationError(fmt.Sprintf("Number of tables must be between 1 and %d", s.maxTables))
	}

	if s.repo.IsInitialized() {
		return errors.ErrTableInitialized
	}

	return s.repo.InitializeTables(numTables)
}

func (s *service) ReserveTables(numCustomers int) (string, int, int, error) {
	if !s.repo.IsInitialized() {
		return "", 0, 0, errors.ErrTableNotInitialized
	}

	if numCustomers <= 0 {
		return "", 0, 0, errors.NewValidationError("Number of customers must be positive")
	}

	tablesNeeded := int(math.Ceil(float64(numCustomers) / float64(s.seatsPerTable)))
	availableTables := s.repo.GetAvailableTables()

	if tablesNeeded > availableTables {
		return "", 0, 0, errors.ErrInsufficientTables
	}

	bookingID := s.generateBookingID()
	booking := models.NewBooking(bookingID, "", numCustomers, tablesNeeded)

	err := s.repo.ReserveTables(*booking)
	if err != nil {
		return "", 0, 0, errors.NewReservationError(err.Error())
	}

	return bookingID, tablesNeeded, availableTables - tablesNeeded, nil
}

func (s *service) CancelReservation(bookingID string) (int, int, error) {
	if !s.repo.IsInitialized() {
		return 0, 0, errors.ErrTableNotInitialized
	}

	tablesFreed, err := s.repo.CancelReservation(bookingID)
	if err != nil {
		return 0, 0, errors.ErrInvalidBookingID
	}

	availableTables := s.repo.GetAvailableTables()
	return tablesFreed, availableTables, nil
}

func (s *service) GetAvailableTables() int {
	return s.repo.GetAvailableTables()
}

func (s *service) generateBookingID() string {
	// Generate seeds for random generator
	rand.Seed(uint64(time.Now().UnixNano()))

	result := make([]byte, s.lengthCode)
	for i := range result {
		result[i] = s.charsetCode[rand.Intn(len(s.charsetCode))]
	}

	return string(result)
}
