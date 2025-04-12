package golog

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewLogger(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "logger_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test cases
	tests := []struct {
		name     string
		filePath string
		wantErr  bool
	}{
		{
			name:     "Valid log file path",
			filePath: filepath.Join(tmpDir, "test.log"),
			wantErr:  false,
		},
		{
			name:     "Invalid directory path",
			filePath: filepath.Join(tmpDir, "nonexistent", "test.log"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewLogger(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && logger == nil {
				t.Error("Expected logger to be non-nil")
			}
			if logger != nil {
				logger.Close()
			}
		})
	}
}

func TestLogLevels(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "logger_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logFile := filepath.Join(tmpDir, "test.log")
	logger, err := NewLogger(logFile)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Test different log levels
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warning("warning message")
	logger.Error("error message")

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Verify console output (DEBUG and INFO)
	if !strings.Contains(output, "DEBUG: debug message") {
		t.Error("Debug message not found in console output")
	}
	if !strings.Contains(output, "INFO: info message") {
		t.Error("Info message not found in console output")
	}

	// Read log file content
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	fileOutput := string(content)

	// Verify file output (WARNING and ERROR)
	if !strings.Contains(fileOutput, "WARNING: warning message") {
		t.Error("Warning message not found in file output")
	}
	if !strings.Contains(fileOutput, "ERROR: error message") {
		t.Error("Error message not found in file output")
	}
}

func TestLogRotation(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "logger_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logFile := filepath.Join(tmpDir, "test.log")
	logger, err := NewLogger(logFile)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Set a small max file size for testing
	logger.maxFileSize = 100

	// Write enough data to trigger rotation
	for i := 0; i < 10; i++ {
		logger.Warning(fmt.Sprintf("test message %d with some padding to increase size", i))
	}

	// Check if rotation files exist
	_, err = os.Stat(logFile + ".1")
	if os.IsNotExist(err) {
		t.Error("Expected rotation file .1 to exist")
	}
}

func TestWithFields(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "logger_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logFile := filepath.Join(tmpDir, "test.log")
	logger, err := NewLogger(logFile)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Add fields and log
	fields := map[string]interface{}{
		"user_id": 123,
		"action":  "test",
	}
	loggerWithFields := logger.WithFields(fields)
	loggerWithFields.Warning("test message")

	// Read log file content
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	fileOutput := string(content)

	// Verify structured fields in output
	if !strings.Contains(fileOutput, `"user_id":123`) {
		t.Error("user_id field not found in structured output")
	}
	if !strings.Contains(fileOutput, `"action":"test"`) {
		t.Error("action field not found in structured output")
	}
}
