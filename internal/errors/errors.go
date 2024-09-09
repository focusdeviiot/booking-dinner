package errors

import (
	"errors"
	"fmt"
)

var (
	ErrTableInitialized     = errors.New("tables have already been initialized")
	ErrTableNotInitialized  = errors.New("tables have not been initialized")
	ErrInsufficientTables   = errors.New("not enough tables available for the reservation")
	ErrInvalidBookingID     = errors.New("invalid booking ID")
	ErrInvalidCustomerCount = errors.New("invalid customer count")
	ErrMaxTablesExceeded    = errors.New("maximum number of tables exceeded")
)

type RestaurantError struct {
	Code    string
	Message string
}

// Error returns the error message
func (e *RestaurantError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewRestaurantError creates a new RestaurantError
func NewRestaurantError(code string, message string) *RestaurantError {
	return &RestaurantError{
		Code:    code,
		Message: message,
	}
}

// Common error codes
const (
	ErrCodeInitialization = "INIT_ERROR"
	ErrCodeReservation    = "RESERVATION_ERROR"
	ErrCodeCancellation   = "CANCELLATION_ERROR"
	ErrCodeValidation     = "VALIDATION_ERROR"
)

// Helper functions to create specific errors
func NewInitializationError(msg string) *RestaurantError {
	return NewRestaurantError(ErrCodeInitialization, msg)
}

func NewReservationError(msg string) *RestaurantError {
	return NewRestaurantError(ErrCodeReservation, msg)
}

func NewCancellationError(msg string) *RestaurantError {
	return NewRestaurantError(ErrCodeCancellation, msg)
}

func NewValidationError(msg string) *RestaurantError {
	return NewRestaurantError(ErrCodeValidation, msg)
}
