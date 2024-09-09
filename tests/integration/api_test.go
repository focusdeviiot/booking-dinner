package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"booking-dinner/internal/api"
	"booking-dinner/internal/api/handlers"
	"booking-dinner/internal/config"
	"booking-dinner/internal/domain/restaurant"
	"booking-dinner/internal/storage/memory"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	cfg := &config.Config{
		Restaurant: config.RestaurantConfig{
			MaxTables:     20,
			SeatsPerTable: 4,
			Code: config.CodeConfig{
				Charset: "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
				Length:  6,
			},
		}}
	repo := memory.NewRestaurantRepository()
	service := restaurant.NewService(repo, cfg.Restaurant.SeatsPerTable, cfg.Restaurant.MaxTables, cfg.Restaurant.Code.Charset, cfg.Restaurant.Code.Length)
	handler := handlers.NewRestaurantHandler(service)

	app := fiber.New()
	api.SetupRoutes(app, handler)

	return app
}

func TestInitializeTables(t *testing.T) {
	app := setupTestApp()

	// Test successful initialization
	req := httptest.NewRequest(http.MethodPost, "/api/v1/initialize", strings.NewReader(`{"tables": 10}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test initialization with invalid number of tables
	req = httptest.NewRequest(http.MethodPost, "/api/v1/initialize", strings.NewReader(`{"tables": 0}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestReserveTables(t *testing.T) {
	app := setupTestApp()

	// Initialize tables
	initReq := httptest.NewRequest(http.MethodPost, "/api/v1/initialize", strings.NewReader(`{"tables": 10}`))
	initReq.Header.Set("Content-Type", "application/json")
	app.Test(initReq)

	// Test successful reservation
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reserve", strings.NewReader(`{"customers": 3}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	// Check the structure of the response
	assert.True(t, result["success"].(bool))
	assert.Equal(t, "Reservation successful", result["message"])

	// Check the data field
	data, ok := result["data"].(map[string]interface{})
	assert.True(t, ok, "Data field should be a map")
	assert.Contains(t, data, "bookingID")
	assert.Contains(t, data, "tablesBooked")
	assert.Contains(t, data, "remainingTables")

	// Test reservation with insufficient tables
	req = httptest.NewRequest(http.MethodPost, "/api/v1/reserve", strings.NewReader(`{"customers": 50}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCancelReservation(t *testing.T) {
	app := setupTestApp()

	// Initialize and make a reservation
	initReq := httptest.NewRequest(http.MethodPost, "/api/v1/initialize", strings.NewReader(`{"tables": 10}`))
	initReq.Header.Set("Content-Type", "application/json")
	app.Test(initReq)

	reserveReq := httptest.NewRequest(http.MethodPost, "/api/v1/reserve", strings.NewReader(`{"customers": 3}`))
	reserveReq.Header.Set("Content-Type", "application/json")
	reserveResp, _ := app.Test(reserveReq)

	var reserveResult map[string]interface{}
	err := json.NewDecoder(reserveResp.Body).Decode(&reserveResult)
	assert.NoError(t, err)
	assert.True(t, reserveResult["success"].(bool))
	assert.Equal(t, "Reservation successful", reserveResult["message"])
	bookingID := reserveResult["data"].(map[string]interface{})["bookingID"].(string)

	// Test successful cancellation
	cancelReq := httptest.NewRequest(http.MethodPost, "/api/v1/cancel", strings.NewReader(`{"bookingID":"`+bookingID+`"}`))
	cancelReq.Header.Set("Content-Type", "application/json")
	cancelResp, _ := app.Test(cancelReq)

	assert.Equal(t, http.StatusOK, cancelResp.StatusCode)

	var cancelResult map[string]interface{}
	json.NewDecoder(cancelResp.Body).Decode(&cancelResult)
	data, ok := cancelResult["data"].(map[string]interface{})
	assert.True(t, ok, "Data field should be a map")
	assert.Contains(t, data, "tablesFreed")
	assert.Contains(t, data, "remainingTables")

	// Test cancellation with invalid booking ID
	invalidCancelReq := httptest.NewRequest(http.MethodPost, "/api/v1/cancel", strings.NewReader(`{"bookingID":"invalid-id"}`))
	invalidCancelReq.Header.Set("Content-Type", "application/json")
	invalidCancelResp, _ := app.Test(invalidCancelReq)

	assert.Equal(t, http.StatusNotFound, invalidCancelResp.StatusCode)
}

func TestEdgeCases(t *testing.T) {
	app := setupTestApp()

	// Test initializing with 0 tables
	zeroTablesReq := httptest.NewRequest(http.MethodPost, "/api/v1/initialize", strings.NewReader(`{"tables": 0}`))
	zeroTablesReq.Header.Set("Content-Type", "application/json")
	zeroTablesResp, _ := app.Test(zeroTablesReq)

	assert.Equal(t, http.StatusBadRequest, zeroTablesResp.StatusCode)

	// Initialize with valid number of tables
	initReq := httptest.NewRequest(http.MethodPost, "/api/v1/initialize", strings.NewReader(`{"tables": 10}`))
	initReq.Header.Set("Content-Type", "application/json")
	app.Test(initReq)

	// Test reserving for 0 customers
	zeroCustomersReq := httptest.NewRequest(http.MethodPost, "/api/v1/reserve", strings.NewReader(`{"customers": 0}`))
	zeroCustomersReq.Header.Set("Content-Type", "application/json")
	zeroCustomersResp, _ := app.Test(zeroCustomersReq)

	assert.Equal(t, http.StatusBadRequest, zeroCustomersResp.StatusCode)

	// Test reserving for more customers than available seats
	overCapacityReq := httptest.NewRequest(http.MethodPost, "/api/v1/reserve", strings.NewReader(`{"customers": 50}`))
	overCapacityReq.Header.Set("Content-Type", "application/json")
	overCapacityResp, _ := app.Test(overCapacityReq)

	assert.Equal(t, http.StatusBadRequest, overCapacityResp.StatusCode)

	// Test health check endpoint
	healthReq := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	healthResp, _ := app.Test(healthReq)

	assert.Equal(t, http.StatusOK, healthResp.StatusCode)

	var healthResult map[string]interface{}
	json.NewDecoder(healthResp.Body).Decode(&healthResult)
	assert.Equal(t, "ok", healthResult["status"])
}
