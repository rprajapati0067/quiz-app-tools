# quiz-app-tools

A comprehensive logging package built on logrus for Go applications.

## Features

- 🚀 **Easy to Use**: Simple API that can be used anywhere in your application
- 📝 **Structured Logging**: Support for structured logging with fields
- 🎨 **Multiple Formats**: Text format for development, JSON format for production
- 🔧 **Configurable**: Set log levels, default fields, and formatters
- 🔄 **Singleton Pattern**: One logger instance accessible throughout your app
- ✅ **Auto-initialization**: Logger initializes automatically if not initialized

## Installation

```bash
go get github.com/sirupsen/logrus
```

## Quick Start

```go
package main

import (
    "quiz-app-tools/logger"
)

func main() {
    // Initialize logger
    logger.Init()
    
    // Use anywhere in your app
    logger.Info("Application started")
    logger.Debug("Debug information")
    logger.Warn("This is a warning")
    logger.Error("An error occurred")
}
```

## Usage

### Initialization

The logger can be initialized in two ways:

#### Text Format (Development)
```go
logger.Init()
```

#### JSON Format (Production)
```go
logger.InitJSON()
```

### Basic Logging

```go
logger.Debug("Debug message")
logger.Info("Info message")
logger.Warn("Warning message")
logger.Error("Error message")
logger.Fatal("Fatal message") // Exits the program
logger.Panic("Panic message") // Panics
```

### Formatted Logging

```go
logger.Debugf("User %s logged in", username)
logger.Infof("Processing %d items", count)
logger.Errorf("Error: %v", err)
```

### Structured Logging

#### With Fields
```go
logger.WithFields(map[string]interface{}{
    "user_id": 123,
    "action":  "login",
    "ip":      "192.168.1.1",
}).Info("User logged in")
```

#### With Single Field
```go
logger.WithField("request_id", "req-12345").Info("Processing request")
```

#### Chaining Fields
```go
logger.WithField("service", "auth").
    WithField("method", "POST").
    WithField("path", "/api/login").
    Info("API request received")
```

### Configuration

#### Set Log Level
```go
// Available levels: trace, debug, info, warn, error, fatal, panic
logger.SetLevel("debug")
```

#### Set Default Fields
Default fields are included in all log entries:

```go
logger.SetDefaultFields(map[string]interface{}{
    "app":     "quiz-app-tools",
    "version": "1.0.0",
    "env":     "production",
})
```

Now all log messages will automatically include these fields.

### Using in Different Packages

The logger can be imported and used in any package:

```go
package handlers

import "quiz-app-tools/logger"

func HandleRequest() {
    logger.WithField("handler", "HandleRequest").Info("Request received")
}
```

## Example

See `examples/main.go` for a complete example.

```bash
go run examples/main.go
```

## Log Levels

- **Trace**: Very detailed logging, usually only enabled when debugging specific problems
- **Debug**: Detailed information for debugging
- **Info**: General informational messages
- **Warn**: Warning messages for potentially harmful situations
- **Error**: Error messages for error events that might still allow the application to continue
- **Fatal**: Fatal messages that cause the application to exit
- **Panic**: Panic messages that cause the application to panic

## Best Practices

1. **Initialize Early**: Call `logger.Init()` or `logger.InitJSON()` at the start of your `main()` function
2. **Use Appropriate Levels**: Use Debug for development info, Info for general events, Error for errors
3. **Add Context**: Use `WithFields()` to add relevant context to your logs
4. **Use JSON in Production**: Use `InitJSON()` in production for better log aggregation tools
5. **Set Default Fields**: Use `SetDefaultFields()` to add application-wide context like app name, version, environment

## Integration with Other Projects

Once this package is pushed to GitHub, you can use it in your other Go projects:

```bash
go get github.com/YOUR_USERNAME/quiz-app-tools/logger
```

### Integration Guides

- **[Integration Guide](INTEGRATION_GUIDE.md)**: Complete guide on pushing to GitHub and integrating in any project
- **[User Management API Example](USER_MANAGEMENT_API_EXAMPLE.md)**: Quick-start guide specifically for API projects with complete examples
- **[Cross-Package Usage](CROSS_PACKAGE_USAGE.md)**: How to use logger across different packages without passing it around

## License

MIT