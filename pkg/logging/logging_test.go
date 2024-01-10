package logging

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger(INFO)
	if logger.logLevel != INFO {
		t.Errorf("Expected log level %v, got %v", INFO, logger.logLevel)
	}
}

func logDateString() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

func TestLogInfo(t *testing.T) {
	// Create a logger with INFO level
	logger := NewLogger(INFO)

	// Capture log output
	var buf bytes.Buffer
	logger.logger.SetOutput(&buf)

	// Log an informational message
	logger.LogInfo("Test Info Message")

	// Verify log output
	expectedOutput := fmt.Sprintf("GoWise %s INFO: %s\n", logDateString(), "Test Info Message")
	if buf.String() != expectedOutput {
		t.Errorf("Expected log output %q, got %q", expectedOutput, buf.String())
	}
}

func TestLogError(t *testing.T) {
	// Create a logger with ERROR level
	logger := NewLogger(ERROR)

	// Capture log output
	var buf bytes.Buffer
	logger.logger.SetOutput(&buf)

	// Log an error message
	testError := errors.New("Test Error")
	logger.LogError(testError)

	// Verify log output
	expectedOutput := fmt.Sprintf("GoWise %s ERROR: %v\n", logDateString(), testError)
	if buf.String() != expectedOutput {
		t.Errorf("Expected log output %q, got %q", expectedOutput, buf.String())
	}
}
