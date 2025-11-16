package handlers

import (
	"errors"

	"github.com/YOUR_USERNAME/quiz-app-tools/logger" // Just import - no passing needed!

	"example/services"
)

type UserHandler struct {
	userService *services.UserService
	// Notice: No logger field here! We don't need to pass it.
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	// Notice: No logger parameter needed!
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(email string) error {
	// Use logger directly - it's the same instance from main.go
	// Default fields ("app", "env") are automatically included!
	logger.WithFields(map[string]interface{}{
		"handler": "CreateUser",
		"email":   email,
	}).Info("Creating new user")

	// Call service - logger is available there too!
	err := h.userService.CreateUser(email)
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"handler": "CreateUser",
			"email":   email,
			"error":   err.Error(),
		}).Error("Failed to create user")
		return err
	}

	logger.WithField("email", email).Info("User created successfully")
	return nil
}

func (h *UserHandler) GetUser(userID string) (*User, error) {
	logger.WithField("user_id", userID).Debug("Fetching user")

	user, err := h.userService.GetUser(userID)
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		}).Warn("User not found")
		return nil, errors.New("user not found")
	}

	logger.WithField("user_id", userID).Info("User retrieved successfully")
	return user, nil
}

type User struct {
	ID    string
	Email string
}

