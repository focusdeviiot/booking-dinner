package api

import (
	"booking-dinner/internal/api/handlers"
	"booking-dinner/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures the routes for the API
func SetupRoutes(app *fiber.App, handler handlers.Handler) {
	// API group
	api := app.Group("/api")
	api = api.Group("/v1")

	// Middleware
	api.Use(middleware.Logger())
	api.Use(middleware.Recover())

	// Routes

	// Health check
	api.Get("/health", HealthCheck)
}

// HealthCheck handler for the health check endpoint
func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "ok",
		"message": "Server is healthy",
	})
}
