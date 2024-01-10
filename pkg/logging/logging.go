// Package logging provides a flexible logging mechanism with support for multiple logger instances,
// log levels, and customizable output destinations.
package logging

import (
	"log"
	"os"
)

// LogLevel represents the severity level of a log message.
type LogLevel int

// Log levels supported by the logging package.
const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// LoggerInterface is an interface that defines the methods required for logging.
type LoggerInterface interface {
	LogInfo(message string)
	LogError(err error)
}

// Logger is a type that provides logging functionality.
type Logger struct {
	logger   *log.Logger
	logLevel LogLevel
}

// NewLogger creates a new logger instance with the specified log level.
func NewLogger(logLevel LogLevel) *Logger {
	return &Logger{
		logger:   log.New(os.Stdout, "GoWise ", log.Ldate|log.Ltime),
		logLevel: logLevel,
	}
}

// LogInfo logs an informational message.
func (l *Logger) LogInfo(message string) {
	if l.logLevel <= INFO {
		l.logger.Printf("INFO: %s\n", message)
	}
}

// LogError logs an error message.
func (l *Logger) LogError(err error) {
	if l.logLevel <= ERROR {
		l.logger.Printf("ERROR: %v\n", err)
	}
}

// MockLogger is a type that provides mock logging functionality for testing.
type MockLogger struct {
	InfoMessages  []string
	ErrorMessages []string
}

// NewMockLogger creates a new mock logger instance.
func NewMockLogger() *MockLogger {
	return &MockLogger{
		InfoMessages:  make([]string, 0),
		ErrorMessages: make([]string, 0),
	}
}

// LogInfo logs an informational message to the mock logger.
func (m *MockLogger) LogInfo(message string) {
	m.InfoMessages = append(m.InfoMessages, message)
}

// LogError logs an error message to the mock logger.
func (m *MockLogger) LogError(err error) {
	// Record the error message
	m.ErrorMessages = append(m.ErrorMessages, err.Error())
}
