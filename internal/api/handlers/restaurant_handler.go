package handlers

import (
	"booking-dinner/internal/domain/restaurant"
	"booking-dinner/internal/errors"

	"github.com/gofiber/fiber/v2"
)

type RestaurantHandler struct {
	service restaurant.Service
}

func NewRestaurantHandler(service restaurant.Service) *RestaurantHandler {
	return &RestaurantHandler{
		service: service,
	}
}

func (h *RestaurantHandler) InitializeTables(c *fiber.Ctx) error {
	var request struct {
		Tables int `json:"tables"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(NewErrorResponse("Invalid request body", err.Error()))
	}

	if request.Tables <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(NewErrorResponse("Invalid number of tables", "Number of tables must be positive"))
	}

	err := h.service.InitializeTables(request.Tables)
	if err != nil {
		if err == errors.ErrTableInitialized {
			return c.Status(fiber.StatusBadRequest).JSON(NewErrorResponse("Initialization error", err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(NewErrorResponse("Initialization failed", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(NewSuccessResponse("Tables initialized successfully", nil))
}

func (h *RestaurantHandler) ReserveTables(c *fiber.Ctx) error {
	var request struct {
		Customers int `json:"customers"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(NewErrorResponse("Invalid request body", err.Error()))
	}

	if request.Customers <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(NewErrorResponse("Invalid number of customers", "Number of customers must be positive"))
	}

	bookingID, tablesBooked, remainingTables, err := h.service.ReserveTables(request.Customers)
	if err != nil {
		if err == errors.ErrInsufficientTables {
			return c.Status(fiber.StatusBadRequest).JSON(NewErrorResponse("Reservation failed", err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(NewErrorResponse("Reservation failed", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(NewSuccessResponse("Reservation successful", fiber.Map{
		"bookingID":       bookingID,
		"tablesBooked":    tablesBooked,
		"remainingTables": remainingTables,
	}))
}

func (h *RestaurantHandler) CancelReservation(c *fiber.Ctx) error {
	var request struct {
		BookingID string `json:"bookingID"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(NewErrorResponse("Invalid request body", err.Error()))
	}

	tablesFreed, remainingTables, err := h.service.CancelReservation(request.BookingID)
	if err != nil {
		if err == errors.ErrInvalidBookingID {
			return c.Status(fiber.StatusNotFound).JSON(NewErrorResponse("Cancellation failed", err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(NewErrorResponse("Cancellation failed", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(NewSuccessResponse("Reservation cancelled successfully", fiber.Map{
		"tablesFreed":     tablesFreed,
		"remainingTables": remainingTables,
	}))
}
