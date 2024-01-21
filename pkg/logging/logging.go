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
// It allows for different implementations of logging, such as a standard logger
// that logs to an output stream, or a mock logger for testing.
type LoggerInterface interface {
	LogInfo(message string)
	LogError(err error)
}

// Logger is a struct that provides logging functionality to an output stream.
// It contains a log.Logger and a log level, which determines the severity of
// messages that will be logged.
type Logger struct {
	logger   *log.Logger
	logLevel LogLevel
}

// NewLogger creates a new Logger instance with the specified log level.
// The Logger will log messages of a severity equal to or greater than the
// specified log level to os.Stdout.
func NewLogger(logLevel LogLevel) *Logger {
	return &Logger{
		logger:   log.New(os.Stdout, "GoWise ", log.Ldate|log.Ltime),
		logLevel: logLevel,
	}
}

// LogInfo logs an informational message if the log level is INFO or lower.
// The message is prefixed with "INFO: ".
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

// MockLogger is a struct that provides mock logging functionality for testing.
// It records informational and error messages in separate slices instead of
// logging them to an output stream.
type MockLogger struct {
	InfoMessages  []string
	ErrorMessages []string
}

// NewMockLogger creates a new MockLogger instance.
// The MockLogger records messages instead of logging them, which can be useful
// for verifying that the correct messages were logged during testing.
func NewMockLogger() *MockLogger {
	return &MockLogger{
		InfoMessages:  make([]string, 0),
		ErrorMessages: make([]string, 0),
	}
}

// LogInfo records an informational message in the InfoMessages slice.
func (m *MockLogger) LogInfo(message string) {
	m.InfoMessages = append(m.InfoMessages, message)
}

// LogError logs an error message to the mock logger.
func (m *MockLogger) LogError(err error) {
	// Record the error message
	m.ErrorMessages = append(m.ErrorMessages, err.Error())
}
