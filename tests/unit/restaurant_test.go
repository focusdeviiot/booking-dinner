package unit

import (
	"testing"

	"booking-dinner/internal/domain/models"
	"booking-dinner/internal/domain/restaurant"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) InitializeTables(numTables int) error {
	args := m.Called(numTables)
	return args.Error(0)
}

func (m *MockRepository) ReserveTables(booking models.Booking) error {
	args := m.Called(booking)
	return args.Error(0)
}

func (m *MockRepository) CancelReservation(bookingID string) (int, error) {
	args := m.Called(bookingID)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) GetAvailableTables() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockRepository) IsInitialized() bool {
	args := m.Called()
	return args.Bool(0)
}

func TestInitializeTables(t *testing.T) {
	mockRepo := new(MockRepository)
	service := restaurant.NewService(mockRepo, 4, 20, "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 6)

	mockRepo.On("IsInitialized").Return(false)
	mockRepo.On("InitializeTables", 10).Return(nil)

	err := service.InitializeTables(10)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestReserveTables(t *testing.T) {
	mockRepo := new(MockRepository)
	service := restaurant.NewService(mockRepo, 4, 20, "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 6)

	mockRepo.On("IsInitialized").Return(true)
	mockRepo.On("GetAvailableTables").Return(10)
	mockRepo.On("ReserveTables", mock.AnythingOfType("models.Booking")).Return(nil)

	bookingID, tablesBooked, remaining, err := service.ReserveTables(3)
	assert.NoError(t, err)
	assert.NotEmpty(t, bookingID)
	assert.Equal(t, 1, tablesBooked)
	assert.Equal(t, 9, remaining)

	mockRepo.AssertExpectations(t)
}

func TestCancelReservation(t *testing.T) {
	mockRepo := new(MockRepository)
	service := restaurant.NewService(mockRepo, 4, 20, "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 6)

	mockRepo.On("IsInitialized").Return(true)
	mockRepo.On("CancelReservation", "BOOK55").Return(1, nil)
	mockRepo.On("GetAvailableTables").Return(10)

	tablesFreed, remaining, err := service.CancelReservation("BOOK55")
	assert.NoError(t, err)
	assert.Equal(t, 1, tablesFreed)
	assert.Equal(t, 10, remaining)

	mockRepo.AssertExpectations(t)
}

// Add more test cases for edge cases and error scenarios
