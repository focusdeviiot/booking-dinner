package middleware

import (
	"fmt"
	"runtime/debug"
	"time"

	"smart_electricity_tracker_backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Logger returns a middleware that logs HTTP requests
func Logger() fiber.Handler {
	log, _ := logger.New(false) // Assuming we're using development logger

	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log after request is processed
		duration := time.Since(start)
		status := c.Response().StatusCode()
		method := c.Method()
		path := c.Path()
		ip := c.IP()

		log.Info("HTTP Request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("duration", duration),
			zap.String("ip", ip),
		)

		return err
	}
}

// Recover returns a middleware that recovers from panics
func Recover() fiber.Handler {
	log, _ := logger.New(false) // Assuming we're using development logger

	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}

				stack := debug.Stack()

				log.Error("Recovered from panic",
					zap.Error(err),
					zap.String("stack", string(stack)),
				)

				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Internal Server Error",
				})
			}
		}()

		return c.Next()
	}
}

// RequestID returns a middleware that adds a unique request ID to each request
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID() // You need to implement this function
		}
		c.Set("X-Request-ID", requestID)
		return c.Next()
	}
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
