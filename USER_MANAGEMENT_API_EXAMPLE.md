# Quick Start: Using Logger in user-management-api

This is a quick reference guide for integrating the logger into your `user-management-api` project.

## Prerequisites

1. Your `quiz-app-tools` repository is pushed to GitHub
2. Your GitHub repository URL is: `github.com/YOUR_USERNAME/quiz-app-tools` (or your organization path)

## Step-by-Step Integration

### Step 1: Install the Logger Package

In your `user-management-api` directory:

```bash
go get github.com/YOUR_USERNAME/quiz-app-tools/logger
```

**Replace `YOUR_USERNAME` with your actual GitHub username or organization name.**

### Step 2: Initialize Logger in main.go

```go
package main

import (
    "os"
    "net/http"
    
    "github.com/YOUR_USERNAME/quiz-app-tools/logger"
    
    // Your other imports
    // "your-api/handlers"
    // "your-api/middleware"
)

func main() {
    // Initialize logger based on environment
    env := os.Getenv("ENV")
    if env == "production" {
        logger.InitJSON()
        logger.SetLevel("info")
    } else {
        logger.Init()
        logger.SetLevel("debug")
    }
    
    // Set default fields (optional but recommended)
    logger.SetDefaultFields(map[string]interface{}{
        "app": "user-management-api",
        "env": env,
    })
    
    logger.Info("Starting user-management-api...")
    
    // Your router setup
    // router := setupRouter()
    
    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    logger.WithField("port", port).Info("Server starting")
    
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        logger.WithField("error", err.Error()).Fatal("Server failed to start")
    }
}
```

### Step 3: Add Logging to Handlers

Example: `handlers/user_handler.go`

```go
package handlers

import (
    "encoding/json"
    "net/http"
    
    "github.com/YOUR_USERNAME/quiz-app-tools/logger"
)

type UserHandler struct {
    // your dependencies
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
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
    
    // Validate and create user
    // ... your business logic ...
    
    logger.WithField("user_id", user.ID).Info("User created successfully")
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("id")
    
    logger.WithField("user_id", userID).Debug("Fetching user")
    
    // Fetch user logic
    user, err := h.userService.GetUser(userID)
    if err != nil {
        logger.WithFields(map[string]interface{}{
            "user_id": userID,
            "error":   err.Error(),
        }).Warn("User not found")
        
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    
    logger.WithField("user_id", userID).Info("User retrieved successfully")
    json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("id")
    
    logger.WithField("user_id", userID).Info("Updating user")
    
    var updates User
    if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
        logger.WithField("error", err.Error()).Error("Failed to decode update data")
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // Update logic here
    err := h.userService.UpdateUser(userID, updates)
    if err != nil {
        logger.WithFields(map[string]interface{}{
            "user_id": userID,
            "error":   err.Error(),
        }).Error("Failed to update user")
        
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    
    logger.WithField("user_id", userID).Info("User updated successfully")
    w.WriteHeader(http.StatusOK)
}
```

### Step 4: Add Logging to Service Layer

Example: `services/user_service.go`

```go
package services

import (
    "errors"
    
    "github.com/YOUR_USERNAME/quiz-app-tools/logger"
)

type UserService struct {
    // your dependencies
}

func (s *UserService) CreateUser(user *User) error {
    logger.WithField("email", user.Email).Debug("Creating user in service")
    
    // Validate email
    if !s.isValidEmail(user.Email) {
        logger.WithField("email", user.Email).Warn("Invalid email format")
        return errors.New("invalid email format")
    }
    
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

### Step 5: Add Request Logging Middleware

Create: `middleware/logging.go`

```go
package middleware

import (
    "net/http"
    "time"
    
    "github.com/YOUR_USERNAME/quiz-app-tools/logger"
)

func RequestLoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Log incoming request
        logger.WithFields(map[string]interface{}{
            "method": r.Method,
            "path":   r.URL.Path,
            "ip":     r.RemoteAddr,
            "user_agent": r.UserAgent(),
        }).Info("Incoming request")
        
        // Wrap response writer to capture status code
        rw := &responseWriter{
            ResponseWriter: w,
            statusCode:     http.StatusOK,
        }
        
        next.ServeHTTP(rw, r)
        
        // Log completed request
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

Use it in your router:

```go
func setupRouter() *mux.Router {
    router := mux.NewRouter()
    
    // Apply logging middleware
    router.Use(middleware.RequestLoggingMiddleware)
    
    // Your routes
    // router.HandleFunc("/users", handlers.CreateUser).Methods("POST")
    
    return router
}
```

### Step 6: Environment Variables

Create `.env` file:

```bash
# Development
ENV=development
LOG_LEVEL=debug
PORT=8080

# Production (when deploying)
# ENV=production
# LOG_LEVEL=info
# PORT=8080
```

Load with a package like `godotenv`:

```go
import "github.com/joho/godotenv"

func init() {
    if err := godotenv.Load(); err != nil {
        logger.Warn("No .env file found, using system environment variables")
    }
}
```

## Complete Example: main.go

```go
package main

import (
    "os"
    "net/http"
    
    "github.com/YOUR_USERNAME/quiz-app-tools/logger"
    "github.com/joho/godotenv"
    
    "your-api/handlers"
    "your-api/middleware"
    "your-api/services"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        logger.Warn("No .env file found")
    }
    
    // Initialize logger
    env := os.Getenv("ENV")
    if env == "production" {
        logger.InitJSON()
        logger.SetLevel("info")
    } else {
        logger.Init()
        logger.SetLevel("debug")
    }
    
    logger.SetDefaultFields(map[string]interface{}{
        "app": "user-management-api",
        "env": env,
    })
    
    logger.Info("Starting user-management-api...")
    
    // Initialize dependencies
    userService := services.NewUserService(/* dependencies */)
    userHandler := handlers.NewUserHandler(userService)
    
    // Setup router
    router := setupRouter(userHandler)
    
    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    logger.WithField("port", port).Info("Server listening")
    
    if err := http.ListenAndServe(":"+port, router); err != nil {
        logger.WithField("error", err.Error()).Fatal("Server failed to start")
    }
}

func setupRouter(userHandler *handlers.UserHandler) *mux.Router {
    router := mux.NewRouter()
    
    // Apply middleware
    router.Use(middleware.RequestLoggingMiddleware)
    // router.Use(middleware.AuthMiddleware) // if needed
    
    // Routes
    api := router.PathPrefix("/api/v1").Subrouter()
    api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
    api.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
    api.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
    api.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")
    
    return router
}
```

## What Gets Logged

### Development Mode (Text Format)
```
INFO[2024-01-15T10:30:45Z] Starting user-management-api...    app=user-management-api env=development
INFO[2024-01-15T10:30:45Z] Server listening                   app=user-management-api env=development port=8080
INFO[2024-01-15T10:30:50Z] Incoming request                   app=user-management-api env=development ip=127.0.0.1 method=POST path=/api/v1/users
INFO[2024-01-15T10:30:50Z] Creating new user                  app=user-management-api env=development ip=127.0.0.1 method=POST path=/api/v1/users
INFO[2024-01-15T10:30:50Z] User created successfully          app=user-management-api env=development user_id=12345
INFO[2024-01-15T10:30:50Z] Request completed                  app=user-management-api env=development duration_ms=45 method=POST path=/api/v1/users status_code=201
```

### Production Mode (JSON Format)
```json
{"level":"info","msg":"Starting user-management-api...","app":"user-management-api","env":"production","time":"2024-01-15T10:30:45Z"}
{"level":"info","msg":"Server listening","app":"user-management-api","env":"production","port":"8080","time":"2024-01-15T10:30:45Z"}
{"level":"info","msg":"Incoming request","app":"user-management-api","env":"production","ip":"127.0.0.1","method":"POST","path":"/api/v1/users","time":"2024-01-15T10:30:50Z"}
{"level":"info","msg":"User created successfully","app":"user-management-api","env":"production","user_id":"12345","time":"2024-01-15T10:30:50Z"}
```

## Tips

1. **Always initialize logger first** in your `main()` function
2. **Use structured logging** with `WithFields()` for better log aggregation
3. **Include request context** (user_id, request_id, etc.) in your logs
4. **Use appropriate log levels** (Debug for dev, Info for normal flow, Error for errors)
5. **Add request ID** to trace requests across your application
6. **Log errors with full context** to make debugging easier

## Next Steps

1. Push your logger package to GitHub
2. Install it in your user-management-api project
3. Initialize it in main.go
4. Add logging to your handlers, services, and middleware
5. Test it in development mode first
6. Configure for production with JSON logging

Happy logging! 🚀

