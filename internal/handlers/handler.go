package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// Handler interface defines the common methods for all handlers
type Handler interface {
	InitializeTables(c *fiber.Ctx) error
	ReserveTables(c *fiber.Ctx) error
	CancelReservation(c *fiber.Ctx) error
}

// Response is a generic response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(message string, data interface{}) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(message string, err string) Response {
	return Response{
		Success: false,
		Message: message,
		Error:   err,
	}
}
