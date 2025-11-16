# Integration Guide: Using Logger in user-management-api

This guide explains how to push the logger package to GitHub and use it in your `user-management-api` project.

## Step 1: Push Logger Package to GitHub

### 1.1 Initialize Git Repository (if not already done)

```bash
cd /path/to/quiz-app-tools
git init
git add .
git commit -m "Initial commit: Add logger package"
```

### 1.2 Create GitHub Repository

1. Go to [GitHub](https://github.com) and create a new repository named `quiz-app-tools`
2. Make sure it's public (or ensure you have access if private)

### 1.3 Push to GitHub

```bash
git remote add origin https://github.com/YOUR_USERNAME/quiz-app-tools.git
git branch -M main
git push -u origin main
```

Replace `YOUR_USERNAME` with your GitHub username.

**Note:** If your GitHub repo URL is different (e.g., organization repo), use that URL instead.

## Step 2: Use Logger in user-management-api

### 2.1 Install the Logger Package

In your `user-management-api` project:

```bash
cd /path/to/user-management-api
go get github.com/YOUR_USERNAME/quiz-app-tools/logger
```

Or if it's under an organization:

```bash
go get github.com/YOUR_ORG/quiz-app-tools/logger
```

### 2.2 Update go.mod

The logger package and its dependencies will be automatically added to your `go.mod` file.

### 2.3 Initialize Logger in Your API

In your `main.go` or initialization file:

```go
package main

import (
    "os"
    
    "github.com/YOUR_USERNAME/quiz-app-tools/logger"
)

func main() {
    // Initialize logger based on environment
    if os.Getenv("ENV") == "production" {
        logger.InitJSON() // Use JSON format in production
    } else {
        logger.Init() // Use text format in development
    }
    
    // Set log level from environment variable (optional)
    if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
        logger.SetLevel(logLevel)
    }
    
    // Set default fields (optional but recommended)
    logger.SetDefaultFields(map[string]interface{}{
        "app":     "user-management-api",
        "version": "1.0.0",
        "env":     os.Getenv("ENV"),
    })
    
    logger.Info("Application starting...")
    
    // Your application code here
    // ...
}
```

### 2.4 Use Logger in Your Handlers

Example in a user handler:

```go
package handlers

import (
    "encoding/json"
    "net/http"
    
    "github.com/YOUR_USERNAME/quiz-app-tools/logger"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
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
    
    // Create user logic here
    logger.WithField("user_id", user.ID).Info("User created successfully")
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("id")
    
    logger.WithField("user_id", userID).Debug("Fetching user")
    
    // Fetch user logic here
    // If user not found:
    logger.WithFields(map[string]interface{}{
        "user_id": userID,
        "error":   "user not found",
    }).Warn("User not found")
    
    http.Error(w, "User not found", http.StatusNotFound)
}
```

### 2.5 Use Logger in Service Layer

Example in a service:

```go
package services

import (
    "errors"
    
    "github.com/YOUR_USERNAME/quiz-app-tools/logger"
)

type UserService struct {
    // your dependencies
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
    
    if !s.verifyPassword(user.Password, password) {
        logger.WithField("email", email).Warn("Invalid password attempt")
        return nil, errors.New("invalid credentials")
    }
    
    logger.WithField("user_id", user.ID).Info("User authenticated successfully")
    return user, nil
}
```

### 2.6 Use Logger in Middleware

Example authentication middleware:

```go
package middleware

import (
    "net/http"
    "time"
    
    "github.com/YOUR_USERNAME/quiz-app-tools/logger"
)

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        logger.WithFields(map[string]interface{}{
            "method": r.Method,
            "path":   r.URL.Path,
            "ip":     r.RemoteAddr,
        }).Info("Incoming request")
        
        // Create a response writer wrapper to capture status code
        rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        
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

### 2.7 Use Logger for Error Handling

```go
package handlers

import (
    "github.com/YOUR_USERNAME/quiz-app-tools/logger"
)

func HandleError(err error, context map[string]interface{}) {
    fields := map[string]interface{}{
        "error": err.Error(),
    }
    
    // Merge context fields
    for k, v := range context {
        fields[k] = v
    }
    
    logger.WithFields(fields).Error("Operation failed")
}

// Usage
func SomeHandler(w http.ResponseWriter, r *http.Request) {
    if err := someOperation(); err != nil {
        HandleError(err, map[string]interface{}{
            "handler": "SomeHandler",
            "path":    r.URL.Path,
        })
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
}
```

## Step 3: Environment Configuration

Create a `.env` file or set environment variables:

```bash
# .env file
ENV=development
LOG_LEVEL=debug
```

Or in production:

```bash
ENV=production
LOG_LEVEL=info
```

## Step 4: Complete Example Structure

Here's a complete example of how your `user-management-api` structure might look:

```
user-management-api/
├── main.go              # Initialize logger here
├── go.mod
├── go.sum
├── handlers/
│   └── user_handler.go  # Use logger in handlers
├── services/
│   └── user_service.go  # Use logger in services
├── middleware/
│   └── logging.go       # Use logger in middleware
└── models/
    └── user.go
```

## Step 5: Update Logger Package

If you need to update the logger package in your API:

```bash
cd /path/to/user-management-api
go get -u github.com/YOUR_USERNAME/quiz-app-tools/logger
go mod tidy
```

## Best Practices for API Integration

1. **Initialize Early**: Call `logger.Init()` or `logger.InitJSON()` at the very beginning of `main()`

2. **Use Appropriate Log Levels**:
   - `Debug`: Detailed info for debugging
   - `Info`: General application flow
   - `Warn`: Warning situations
   - `Error`: Error events that don't stop the app

3. **Add Context**: Always include relevant context with `WithFields()`:
   ```go
   logger.WithFields(map[string]interface{}{
       "user_id": userID,
       "action":  "update_profile",
   }).Info("User profile updated")
   ```

4. **Log HTTP Requests**: Use middleware to log all incoming requests

5. **Log Errors Properly**: Include full context when logging errors

6. **Use Default Fields**: Set application-wide fields once at startup

7. **Environment-Based Configuration**: Use JSON format in production, text in development

## Troubleshooting

### Issue: Cannot import logger package

**Solution**: Make sure:
- The GitHub repository is public (or you have access if private)
- The module path in `go.mod` matches your GitHub repository path
- You've run `go get` to download the package

### Issue: Logger not initialized error

**Solution**: The logger auto-initializes, but it's better to explicitly call `Init()` or `InitJSON()` in your `main()` function.

### Issue: Want to use a specific version

**Solution**: Use Go modules versioning:
```bash
go get github.com/YOUR_USERNAME/quiz-app-tools/logger@v1.0.0
```

## Next Steps

1. Push your logger package to GitHub
2. Initialize it in your `user-management-api` main.go
3. Add logging to your handlers, services, and middleware
4. Configure environment variables for different environments
5. Test your logging in development and production modes

