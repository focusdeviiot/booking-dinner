package main

import (
	"fmt"
	"log"

	"booking-dinner/internal/api"
	"booking-dinner/internal/api/handlers"
	"booking-dinner/internal/config"
	"booking-dinner/internal/domain/restaurant"
	"booking-dinner/internal/storage/memory"
	"booking-dinner/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger, err := logger.New(cfg.Logger.Production)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Initialize repository
	repo := memory.NewRestaurantRepository()

	// Initialize service
	service := restaurant.NewService(repo, cfg.Restaurant.SeatsPerTable, cfg.Restaurant.MaxTables, cfg.Restaurant.Code.Charset, cfg.Restaurant.Code.Length)

	// Initialize handler
	handler := handlers.NewRestaurantHandler(service)

	// Initialize Fiber app
	app := fiber.New()

	// Setup routes
	api.SetupRoutes(app, handler)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info(fmt.Sprintf("Starting server on %s", addr))
	if err := app.Listen(addr); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to start server: %v", err))
	}
}
