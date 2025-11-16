# Using Logger Across Different Packages

## Key Concept: Singleton Pattern

The logger uses a **singleton pattern** - this means once you initialize it in `main.go`, it's **automatically available to ALL packages** that import it. You **don't need to pass it around**!

## How It Works

1. Initialize once in `main.go` → Sets up the global logger instance
2. Import the logger package in any file → Access the same logger instance
3. Default fields set in `main.go` → Automatically included in all logs everywhere

## Example Project Structure

```
user-management-api/
├── main.go              ← Initialize logger here
├── handlers/
│   ├── user_handler.go  ← Just import and use
│   └── auth_handler.go  ← Just import and use
├── services/
│   └── user_service.go  ← Just import and use
├── middleware/
│   └── logging.go       ← Just import and use
└── models/
    └── user.go
```

## Step-by-Step Usage

### 1. Initialize in main.go (Once)

```go
// main.go
package main

import (
    "os"
    "net/http"
    
    "github.com/YOUR_USERNAME/quiz-app-tools/logger"
    
    "your-api/handlers"
    "your-api/services"
    "your-api/middleware"
)

func main() {
    // Initialize logger ONCE at startup
    if os.Getenv("ENV") == "production" {
        logger.InitJSON()
    } else {
        logger.Init()
    }
    
    // Set default fields ONCE - they'll be available everywhere
    logger.SetDefaultFields(map[string]interface{}{
        "app": "user-management-api",
        "env": os.Getenv("ENV"),
    })
    
    logger.Info("Application starting...")
    
    // Now all your packages can use the logger!
    // You don't need to pass it to them
    userService := services.NewUserService()
    userHandler := handlers.NewUserHandler(userService)
    
    router := setupRouter(userHandler)
    http.ListenAndServe(":8080", router)
}
```

### 2. Use in Handlers Package (No Passing Required!)

```go
// handlers/user_handler.go
package handlers

import (
    "encoding/json"
    "net/http"
    
    "github.com/YOUR_USERNAME/quiz-app-tools/logger"  // Just import!
    
    "your-api/services"
)

type UserHandler struct {
    userService *services.UserService
    // Notice: No logger field needed!
}

func NewUserHandler(userService *services.UserService) *UserHandler {
    return &UserHandler{
        userService: userService,
    }
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // Use logger directly - it's already initialized!
    // Default fields from main.go are automatically included
    logger.WithFields(map[string]interface{}{
        "method": r.Method,
        "path":   r.URL.Path,
        "ip":     r.RemoteAddr,
    }).Info("Creating new user")
    
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        logger.WithField("error", err.Error()).Error("Failed to decode user data")
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    err := h.userService.CreateUser(&user)
    if err != nil {
        logger.WithFields(map[string]interface{}{
            "email": user.Email,
            "error": err.Error(),
        }).Error("Failed to create user")
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    
    logger.WithField("user_id", user.ID).Info("User created successfully")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
```

### 3. Use in Services Package (No Passing Required!)

```go
// services/user_service.go
package services

import (
    "errors"
    
    "github.com/YOUR_USERNAME/quiz-app-tools/logger"  // Just import!
)

type UserService struct {
    repository UserRepository
    // Notice: No logger field needed!
}

func NewUserService(repository UserRepository) *UserService {
    return &UserService{
        repository: repository,
    }
}

func (s *UserService) CreateUser(user *User) error {
    // Use logger directly - default fields from main.go included automatically
    logger.WithField("email", user.Email).Debug("Creating user in service")
    
    // Check if user exists
    exists, err := s.repository.UserExists(user.Email)
    if err != nil {
        logger.WithFields(map[string]interface{}{
            "email": user.Email,
            "error": err.Error(),
        }).Error("Database error while checking user existence")
        return err
    }
    
    if exists {
        logger.WithField("email", user.Email).Warn("User already exists")
        return errors.New("user already exists")
    }
    
    // Create user
    err = s.repository.Create(user)
    if err != nil {
        logger.WithFields(map[string]interface{}{
            "email": user.Email,
            "error": err.Error(),
        }).Error("Failed to create user in database")
        return err
    }
    
    logger.WithFields(map[string]interface{}{
        "user_id": user.ID,
        "email":   user.Email,
    }).Info("User created successfully in database")
    
    return nil
}

func (s *UserService) AuthenticateUser(email, password string) (*User, error) {
    logger.WithField("email", email).Debug("Authenticating user")
    
    user, err := s.repository.FindByEmail(email)
    if err != nil {
        logger.WithFields(map[string]interface{}{
            "email": email,
            "error": err.Error(),
        }).Error("Failed to find user")
        return nil, err
    }
    
    if !s.verifyPassword(user.PasswordHash, password) {
        logger.WithField("email", email).Warn("Invalid password attempt")
        return nil, errors.New("invalid credentials")
    }
    
    logger.WithFields(map[string]interface{}{
        "user_id": user.ID,
        "email":   email,
    }).Info("User authenticated successfully")
    
    return user, nil
}
```

### 4. Use in Middleware Package (No Passing Required!)

```go
// middleware/logging.go
package middleware

import (
    "net/http"
    "time"
    
    "github.com/YOUR_USERNAME/quiz-app-tools/logger"  // Just import!
)

func RequestLoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Use logger directly - default fields included automatically
        logger.WithFields(map[string]interface{}{
            "method":     r.Method,
            "path":       r.URL.Path,
            "ip":         r.RemoteAddr,
            "user_agent": r.UserAgent(),
        }).Info("Incoming request")
        
        rw := &responseWriter{
            ResponseWriter: w,
            statusCode:     http.StatusOK,
        }
        
        next.ServeHTTP(rw, r)
        
        duration := time.Since(start)
        logger.WithFields(map[string]interface{}{
            "method":      r.Method,
            "path":        r.URL.Path,
            "status_code": rw.statusCode,
            "duration_ms": duration.Milliseconds(),
        }).Info("Request completed")
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

### 5. Use in Any Other Package

```go
// utils/validator.go
package utils

import (
    "github.com/YOUR_USERNAME/quiz-app-tools/logger"  // Just import!
)

func ValidateEmail(email string) bool {
    logger.WithField("email", email).Debug("Validating email")
    
    // Validation logic...
    
    isValid := // ... your validation logic
    
    if !isValid {
        logger.WithField("email", email).Warn("Invalid email format")
    }
    
    return isValid
}
```

## Important Points

### ✅ What Happens Automatically

1. **Default Fields**: When you set `logger.SetDefaultFields()` in `main.go`, these fields are **automatically included** in every log message across all packages.

2. **Same Instance**: All packages share the **same logger instance**. Configuration in `main.go` affects all packages.

3. **No Passing Needed**: You don't need to:
   - Add logger as a struct field
   - Pass logger as function parameter
   - Return logger from functions

### ✅ Example: How Default Fields Work

```go
// In main.go
logger.SetDefaultFields(map[string]interface{}{
    "app": "user-management-api",
    "env": "production",
})

// In handlers/user_handler.go
logger.Info("User created")
// Output: {"level":"info","msg":"User created","app":"user-management-api","env":"production",...}

// In services/user_service.go
logger.WithField("user_id", 123).Info("Processing user")
// Output: {"level":"info","msg":"Processing user","app":"user-management-api","env":"production","user_id":123,...}
```

Notice: `app` and `env` fields are automatically included everywhere!

## Alternative: Dependency Injection Pattern (Optional)

If you prefer explicit dependency injection instead of singleton, you can wrap the logger:

```go
// config/logger.go
package config

import (
    "github.com/YOUR_USERNAME/quiz-app-tools/logger"
)

type Logger interface {
    Info(args ...interface{})
    Error(args ...interface{})
    WithField(key string, value interface{}) *logrus.Entry
    // ... other methods
}

func GetLogger() Logger {
    return loggerWrapper{}
}

type loggerWrapper struct{}

func (l loggerWrapper) Info(args ...interface{}) {
    logger.Info(args...)
}

func (l loggerWrapper) Error(args ...interface{}) {
    logger.Error(args...)
}

func (l loggerWrapper) WithField(key string, value interface{}) *logrus.Entry {
    return logger.WithField(key, value)
}
```

Then use it in your structs:

```go
// services/user_service.go
type UserService struct {
    repository UserRepository
    logger     config.Logger  // Inject logger
}

func NewUserService(repository UserRepository, logger config.Logger) *UserService {
    return &UserService{
        repository: repository,
        logger:     logger,
    }
}

func (s *UserService) CreateUser(user *User) error {
    s.logger.WithField("email", user.Email).Info("Creating user")
    // ...
}
```

**However, this is optional!** The singleton pattern works perfectly fine and is simpler.

## Summary

**You don't need to pass the logger!** Just:

1. ✅ Initialize once in `main.go`
2. ✅ Import `"github.com/YOUR_USERNAME/quiz-app-tools/logger"` in any package
3. ✅ Use `logger.Info()`, `logger.Error()`, etc. directly
4. ✅ Default fields from `main.go` are automatically included everywhere

That's it! The singleton pattern handles everything for you. 🚀

