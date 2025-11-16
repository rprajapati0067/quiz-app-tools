package main

import (
	"errors"
	"time"

	"quiz-app-tools/logger"
)

func main() {
	// Initialize logger (text format, good for development)
	logger.Init()

	// Or initialize with JSON format (good for production)
	// logger.InitJSON()

	// Set log level (optional)
	logger.SetLevel("debug")

	// Set default fields that will be included in all logs (optional)
	logger.SetDefaultFields(map[string]interface{}{
		"app":     "quiz-app-tools",
		"version": "1.0.0",
	})

	// Simple logging
	logger.Info("Application started")
	logger.Debug("Debug information")
	logger.Warn("This is a warning")
	logger.Error("An error occurred")

	// Formatted logging
	logger.Infof("User %s logged in at %s", "john_doe", time.Now().Format(time.RFC3339))
	logger.Errorf("Failed to connect to database: %v", errors.New("connection timeout"))

	// Logging with fields (structured logging)
	logger.WithFields(map[string]interface{}{
		"user_id":   123,
		"action":    "login",
		"ip":        "192.168.1.1",
		"timestamp": time.Now(),
	}).Info("User logged in successfully")

	// Logging with a single field
	logger.WithField("request_id", "req-12345").Info("Processing request")

	// Chaining fields
	logger.WithField("service", "auth").
		WithField("method", "POST").
		WithField("path", "/api/login").
		Info("API request received")

	// Error logging with context
	err := errors.New("database connection failed")
	logger.WithFields(map[string]interface{}{
		"error":     err.Error(),
		"component": "database",
		"retries":   3,
	}).Error("Failed to establish database connection")

	// Example: Using in different parts of your application
	simulateUserLogin("user123", "192.168.1.100")
	simulateDatabaseQuery("SELECT * FROM users")
}

func simulateUserLogin(userID, ip string) {
	logger.WithFields(map[string]interface{}{
		"user_id": userID,
		"ip":      ip,
		"action":  "login",
	}).Info("User login attempt")
}

func simulateDatabaseQuery(query string) {
	logger.WithField("query", query).Debug("Executing database query")
	
	// Simulate an error
	if query == "SELECT * FROM users" {
		logger.WithField("query", query).Error("Query execution failed: timeout")
	}
}

