package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Niiaks/campusCart/internal/config"
	"github.com/Niiaks/campusCart/internal/logger"
	"github.com/Niiaks/campusCart/internal/middleware"
	"github.com/Niiaks/campusCart/internal/server"
	"github.com/rs/zerolog"
)

// mockDatabase implements database.Pinger interface for testing
type mockDatabase struct {
	pingFunc  func(ctx context.Context) error
	closeFunc func() error
}

func (m *mockDatabase) Ping(ctx context.Context) error {
	if m.pingFunc != nil {
		return m.pingFunc(ctx)
	}
	return nil
}

func (m *mockDatabase) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}

func setupTestServer(pingFunc func(ctx context.Context) error) *server.Server {
	cfg := &config.Config{
		Primary: config.PrimaryConfig{
			Env: "test",
		},
		Server: config.ServerConfig{
			Port:               "8080",
			ReadTimeout:        30,
			WriteTimeout:       30,
			IdleTimeout:        120,
			CorsAllowedOrigins: []string{"*"},
		},
	}

	log := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()
	loggerService := logger.NewLoggerService(config.DefaultObservabilityConfig())

	// Create mock database
	mockDB := &mockDatabase{
		pingFunc: pingFunc,
	}

	// Create server with mock database
	srv := &server.Server{
		Config:        cfg,
		Logger:        &log,
		LoggerService: loggerService,
		DB:            mockDB,
	}

	return srv
}

func TestHealthHandler_CheckHealth_Healthy(t *testing.T) {
	// Setup: Create a server with a healthy database (ping returns nil)
	srv := setupTestServer(func(ctx context.Context) error {
		return nil // Simulate healthy database
	})

	handler := NewHealthHandler(srv)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	// Add logger to context
	log := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()
	ctx := context.WithValue(req.Context(), middleware.LoggerKey, &log)
	req = req.WithContext(ctx)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute
	handler.CheckHealth(rr, req)

	// Assert status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Assert content type
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "application/json")
	}

	// Parse JSON response
	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Assert response structure
	if response["status"] != "healthy" {
		t.Errorf("expected status to be 'healthy', got %v", response["status"])
	}

	if response["environment"] != "test" {
		t.Errorf("expected environment to be 'test', got %v", response["environment"])
	}

	if response["timestamp"] == nil {
		t.Error("expected timestamp to be present")
	}

	// Assert checks field exists and has database check
	checks, ok := response["checks"].(map[string]interface{})
	if !ok {
		t.Fatal("expected checks to be a map")
	}

	dbCheck, ok := checks["database"].(map[string]interface{})
	if !ok {
		t.Fatal("expected database check to be present")
	}

	if dbCheck["status"] != "healthy" {
		t.Errorf("expected database status to be 'healthy', got %v", dbCheck["status"])
	}

	if dbCheck["response_time"] == nil {
		t.Error("expected database response_time to be present")
	}

	// Assert no error field in healthy check
	if _, hasError := dbCheck["error"]; hasError {
		t.Error("expected no error field in healthy database check")
	}
}

func TestHealthHandler_CheckHealth_UnhealthyDatabase(t *testing.T) {
	// Setup: Create a server with an unhealthy database (ping returns error)
	expectedError := errors.New("connection refused")
	srv := setupTestServer(func(ctx context.Context) error {
		return expectedError // Simulate database connection failure
	})

	handler := NewHealthHandler(srv)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	// Add logger to context
	log := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()
	ctx := context.WithValue(req.Context(), middleware.LoggerKey, &log)
	req = req.WithContext(ctx)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute
	handler.CheckHealth(rr, req)

	// Assert status code is 503 Service Unavailable
	if status := rr.Code; status != http.StatusServiceUnavailable {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusServiceUnavailable)
	}

	// Assert content type
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "application/json")
	}

	// Parse JSON response
	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Assert response structure
	if response["status"] != "unhealthy" {
		t.Errorf("expected status to be 'unhealthy', got %v", response["status"])
	}

	if response["environment"] != "test" {
		t.Errorf("expected environment to be 'test', got %v", response["environment"])
	}

	// Assert checks field exists and has database check
	checks, ok := response["checks"].(map[string]interface{})
	if !ok {
		t.Fatal("expected checks to be a map")
	}

	dbCheck, ok := checks["database"].(map[string]interface{})
	if !ok {
		t.Fatal("expected database check to be present")
	}

	// Assert database check shows unhealthy
	if dbCheck["status"] != "unhealthy" {
		t.Errorf("expected database status to be 'unhealthy', got %v", dbCheck["status"])
	}

	// Assert error message is present
	errorMsg, hasError := dbCheck["error"].(string)
	if !hasError {
		t.Error("expected error field in unhealthy database check")
	}

	if errorMsg != expectedError.Error() {
		t.Errorf("expected error message to be '%s', got '%s'", expectedError.Error(), errorMsg)
	}

	// Assert response_time is present
	if dbCheck["response_time"] == nil {
		t.Error("expected database response_time to be present even in unhealthy state")
	}
}

func TestHealthHandler_CheckHealth_DatabaseTimeout(t *testing.T) {
	// Setup: Create a server where database ping times out
	srv := setupTestServer(func(ctx context.Context) error {
		// Simulate a slow database that respects context timeout
		select {
		case <-time.After(10 * time.Second):
			return nil
		case <-ctx.Done():
			return ctx.Err() // Return context deadline exceeded
		}
	})

	handler := NewHealthHandler(srv)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	// Add logger to context
	log := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()
	ctx := context.WithValue(req.Context(), middleware.LoggerKey, &log)
	req = req.WithContext(ctx)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute
	handler.CheckHealth(rr, req)

	// Assert status code is 503 Service Unavailable
	if status := rr.Code; status != http.StatusServiceUnavailable {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusServiceUnavailable)
	}

	// Parse JSON response
	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Assert overall status is unhealthy
	if response["status"] != "unhealthy" {
		t.Errorf("expected status to be 'unhealthy', got %v", response["status"])
	}

	// Assert checks field has database check with error
	checks, ok := response["checks"].(map[string]interface{})
	if !ok {
		t.Fatal("expected checks to be a map")
	}

	dbCheck, ok := checks["database"].(map[string]interface{})
	if !ok {
		t.Fatal("expected database check to be present")
	}

	if dbCheck["status"] != "unhealthy" {
		t.Errorf("expected database status to be 'unhealthy', got %v", dbCheck["status"])
	}

	// Assert error contains context-related message
	errorMsg, hasError := dbCheck["error"].(string)
	if !hasError {
		t.Error("expected error field in database check after timeout")
	}

	if errorMsg == "" {
		t.Error("expected non-empty error message after timeout")
	}
}
