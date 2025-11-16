package main

import (
	"os"

	"github.com/YOUR_USERNAME/quiz-app-tools/logger"

	"example/handlers"
	"example/services"
)

// This example demonstrates how the logger initialized here
// is automatically available in all packages without passing it

func main() {
	// STEP 1: Initialize logger ONCE here
	if os.Getenv("ENV") == "production" {
		logger.InitJSON()
	} else {
		logger.Init()
	}

	// STEP 2: Set default fields ONCE - they'll be available everywhere
	logger.SetDefaultFields(map[string]interface{}{
		"app": "user-management-api",
		"env": os.Getenv("ENV"),
	})

	logger.Info("Application starting...")

	// STEP 3: Create your services/handlers
	// Notice: We DON'T pass the logger - it's already available globally!
	userService := services.NewUserService()
	userHandler := handlers.NewUserHandler(userService)

	// Now use your handlers
	// The logger in handlers/service will automatically use the same
	// logger instance initialized here, including the default fields!

	_ = userHandler // Prevent unused variable warning
}

