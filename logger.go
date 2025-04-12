package golog

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR

	// Default values for log rotation
	defaultMaxSize    = 10 * 1024 * 1024 // 10MB
	defaultMaxBackups = 5
)

type Logger struct {
	debug   *log.Logger
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger
	file    *os.File

	// Configuration
	minLevel    LogLevel
	maxFileSize int64
	maxBackups  int
	filePath    string
	currentSize int64
	mu          sync.Mutex

	// Structured logging
	fields map[string]interface{}
}

// NewLogger initializes the logger and writes WARNING and ERROR logs to a file
func NewLogger(logFilePath string) (*Logger, error) {
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return &Logger{
		debug:       log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime),
		info:        log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
		warning:     log.New(file, "WARNING: ", log.Ldate|log.Ltime),
		error:       log.New(file, "ERROR: ", log.Ldate|log.Ltime),
		file:        file,
		minLevel:    DEBUG,
		maxFileSize: defaultMaxSize,
		maxBackups:  defaultMaxBackups,
		filePath:    logFilePath,
		fields:      make(map[string]interface{}),
	}, nil
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.minLevel <= DEBUG {
		message := l.formatMessage("DEBUG", format, v...)
		l.debug.Print(message)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.minLevel <= INFO {
		message := l.formatMessage("INFO", format, v...)
		l.info.Print(message)
	}
}

func (l *Logger) Warning(format string, v ...interface{}) {
	if l.minLevel <= WARNING {
		message := l.formatMessage("WARNING", format, v...)
		l.warning.Print(message)
		l.checkRotation(len(message))
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	if l.minLevel <= ERROR {
		message := l.formatMessage("ERROR", format, v...)
		l.error.Print(message)
		l.checkRotation(len(message))
	}
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level LogLevel) {
	l.minLevel = level
}

// WithFields adds structured fields to the log output
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
newLogger := Logger{
    debug:       l.debug,
    info:        l.info,
    warning:     l.warning,
    error:       l.error,
    file:        l.file,
    minLevel:    l.minLevel,
    maxFileSize: l.maxFileSize,
    maxBackups:  l.maxBackups,
    filePath:    l.filePath,
    currentSize: l.currentSize,
    fields:      make(map[string]interface{}),
}
	newLogger.fields = make(map[string]interface{}, len(l.fields)+len(fields))
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}
	for k, v := range fields {
		newLogger.fields[k] = v
	}
	return &newLogger
}

// formatMessage formats the log message with timestamp and structured fields
func (l *Logger) formatMessage(level string, format string, v ...interface{}) string {
	message := fmt.Sprintf(format, v...)
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	if len(l.fields) == 0 {
		return fmt.Sprintf("%s %s", timestamp, message)
	}

	data := map[string]interface{}{
		"timestamp": timestamp,
		"level":     level,
		"message":   message,
	}
	for k, v := range l.fields {
		data[k] = v
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Sprintf("%s %s", timestamp, message)
	}
	return string(jsonData)
}

// checkRotation checks if log rotation is needed and performs rotation if necessary
func (l *Logger) checkRotation(messageSize int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.currentSize += int64(messageSize)
	if l.currentSize >= l.maxFileSize {
		l.rotate()
	}
}

// rotate performs log file rotation
func (l *Logger) rotate() {
	// Close current file
	l.file.Close()

	// Rotate backup files
	for i := l.maxBackups - 1; i > 0; i-- {
		oldPath := fmt.Sprintf("%s.%d", l.filePath, i)
		newPath := fmt.Sprintf("%s.%d", l.filePath, i+1)
		os.Rename(oldPath, newPath)
	}

	// Rename current log file
	os.Rename(l.filePath, l.filePath+".1")

	// Create new log file
	file, err := os.OpenFile(l.filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err == nil {
		l.file = file
		l.warning.SetOutput(file)
		l.error.SetOutput(file)
		l.currentSize = 0
	}
}

// Close closes the log file
func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}
