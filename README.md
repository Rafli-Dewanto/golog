# GoLog - Simple and Flexible Logging Package for Go

GoLog is a lightweight, feature-rich logging package for Go applications that provides structured logging, multiple log levels, file rotation, and field-based contextual logging.

## Features

- Multiple log levels (DEBUG, INFO, WARNING, ERROR)
- File-based logging with automatic rotation
- Console output for DEBUG and INFO levels
- Structured logging with JSON format support
- Field-based contextual logging
- Thread-safe operations

## Installation

To install GoLog, use `go get`:

```bash
go get github.com/Rafli-Dewanto/golog
```

## Usage

### Basic Logging

```go
package main

import "github.com/Rafli-Dewanto/golog"

func main() {
    // Create a new logger instance
    logger, err := golog.NewLogger("app.log")
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    // Log messages with different levels
    logger.Debug("Debug message")    // Outputs to console
    logger.Info("Info message")     // Outputs to console
    logger.Warning("Warning message") // Outputs to file
    logger.Error("Error message")   // Outputs to file
}
```

### Structured Logging with Fields

```go
func main() {
    logger, _ := golog.NewLogger("app.log")
    defer logger.Close()

    // Add contextual fields to your logs
    fields := map[string]interface{}{
        "user_id": 123,
        "action": "login",
    }

    // Create a new logger instance with fields
    contextLogger := logger.WithFields(fields)

    // Log with fields - output will be in JSON format
    contextLogger.Info("User logged in")
}
```

## Log Format

### Regular Logs

Regular logs without fields are formatted as:

```
2006-01-02 15:04:05 [LEVEL] Message
```

### Structured Logs

When using fields, logs are formatted as JSON:

```json
{
  "timestamp": "2006-01-02 15:04:05",
  "level": "INFO",
  "message": "User logged in",
  "user_id": 123,
  "action": "login"
}
```

## Log Levels

- `DEBUG`: Development-level information (console output)
- `INFO`: General operational information (console output)
- `WARNING`: Warning messages (file output)
- `ERROR`: Error messages (file output)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
