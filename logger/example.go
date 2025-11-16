package logger

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

// ExampleUsage demonstrates how to use the logger package
func ExampleUsage() {
	// Initialize logger (text format, good for development)
	Init()

	// Or initialize with JSON format (good for production)
	// InitJSON()

	// Set log level (optional)
	SetLevel("debug")

	// Set default fields that will be included in all logs (optional)
	SetDefaultFields(logrus.Fields{
		"app":     "quiz-app-tools",
		"version": "1.0.0",
	})

	// Simple logging
	Info("Application started")
	Debug("Debug information")
	Warn("This is a warning")
	Error("An error occurred")

	// Formatted logging
	Infof("User %s logged in at %s", "john_doe", time.Now().Format(time.RFC3339))
	Errorf("Failed to connect to database: %v", errors.New("connection timeout"))

	// Logging with fields (structured logging)
	WithFields(logrus.Fields{
		"user_id":   123,
		"action":    "login",
		"ip":        "192.168.1.1",
		"timestamp": time.Now(),
	}).Info("User logged in successfully")

	// Logging with a single field
	WithField("request_id", "req-12345").Info("Processing request")

	// Chaining fields
	WithField("service", "auth").
		WithField("method", "POST").
		WithField("path", "/api/login").
		Info("API request received")

	// Error logging with context
	err := errors.New("database connection failed")
	WithFields(logrus.Fields{
		"error":     err.Error(),
		"component": "database",
		"retries":   3,
	}).Error("Failed to establish database connection")
}
