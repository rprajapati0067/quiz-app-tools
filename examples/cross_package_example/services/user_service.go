package services

import (
	"errors"

	"github.com/YOUR_USERNAME/quiz-app-tools/logger" // Just import - no passing needed!
)

type UserService struct {
	// Notice: No logger field here! We don't need to pass it.
	// The logger is available globally once initialized in main.go
}

func NewUserService() *UserService {
	// Notice: No logger parameter needed!
	return &UserService{}
}

func (s *UserService) CreateUser(email string) error {
	// Use logger directly - it's the same instance from main.go
	// Default fields ("app", "env") are automatically included in all logs!
	logger.WithField("email", email).Debug("Creating user in service layer")

	// Validate email
	if !s.isValidEmail(email) {
		logger.WithField("email", email).Warn("Invalid email format")
		return errors.New("invalid email format")
	}

	// Check if user exists (simulated)
	exists := s.userExists(email)
	if exists {
		logger.WithField("email", email).Warn("User already exists")
		return errors.New("user already exists")
	}

	// Create user (simulated)
	logger.WithFields(map[string]interface{}{
		"email": email,
		"layer": "service",
	}).Info("User created successfully in service")

	return nil
}

func (s *UserService) GetUser(userID string) (*User, error) {
	logger.WithField("user_id", userID).Debug("Fetching user from service")

	// Simulate database lookup
	if userID == "" {
		logger.WithField("user_id", userID).Error("Invalid user ID")
		return nil, errors.New("invalid user ID")
	}

	// Simulate user not found
	if userID == "999" {
		logger.WithField("user_id", userID).Warn("User not found in database")
		return nil, errors.New("user not found")
	}

	logger.WithField("user_id", userID).Info("User retrieved from service")
	return &User{ID: userID, Email: "user@example.com"}, nil
}

// Helper methods
func (s *UserService) isValidEmail(email string) bool {
	return len(email) > 0 && email != ""
}

func (s *UserService) userExists(email string) bool {
	// Simulated check
	return email == "existing@example.com"
}

type User struct {
	ID    string
	Email string
}

