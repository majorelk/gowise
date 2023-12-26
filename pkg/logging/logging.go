// pkg/logging/logging.go
package logging

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	// Initialize the logger
	logger = log.New(os.Stdout, "GoWise ", log.Ldate|log.Ltime)
}

// LogInfo logs informational messages
func LogInfo(message string) {
	logger.Printf("INFO: %s\n", message)
}

// LogError logs error messages
func LogError(err error) {
	logger.Printf("ERROR: %v\n", err)
}

