// Package logging provides a flexible logging mechanism with support for multiple logger instances,
// log levels, and customizable output destinations.
package logging

import (
	"log"
	"os"
	"time"
)

// Level represents the severity leval of a log message.
type LogLevel int

// Log levels supported the logging package.
const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

type Logger struct {
	logger  *log.Logger
	logLevel LogLevel
}

// NewLogger creates a new logger instance with the specified log level.
func NewLogger(logLevel LogLevel) *Logger {
	return &Logger{
		logger:  log.New(os.Stdout, "GoWise ", log.Ldate|log.Ltime),
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

// Helper function to format the date string for logs
func logDateString() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

// Example usage of the logging package.
// func Example() {
	// Create logger instances with different log levels
	// infoLogger := logging.NewLogger(logging.INFO)
	// debugLogger := logging.NewLogger(logging.DEBUG)

	// Log messages using the respective loggers
	// infoLogger.LogInfo("This is an informational message")
	// debugLogger.LogInfo("This is a debug message")
	// debugLogger.LogError(errors.New("This is an error"))
// }

