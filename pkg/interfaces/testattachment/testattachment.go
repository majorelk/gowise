package testattachment

import (
	"errors"
	"fmt"
	"os"
)

// TestAttachment represents a file attached to a TestResult with an optional description.
type TestAttachment struct {
	// FilePath is the absolute file path to the attachment file.
	FilePath string

	// Description is the user-specified description of the attachment. It may be empty.
	Description string
}

// NewTestAttachment creates a TestAttachment to represent a file attached to a test result.
func NewTestAttachment(filePath, description string) (TestAttachment, error) {
	// Perform any validation or checks here
	if err := validateFilePath(filePath); err != nil {
		return TestAttachment{}, fmt.Errorf("error creating TestAttachment: %w", err)
	}

	return TestAttachment{
		FilePath:    filePath,
		Description: description,
	}, nil
}

// validateFilePath is an example validation function, you can replace it with your actual validation logic.
func validateFilePath(filePath string) error {
	if filePath == "" {
		return errors.New("file path cannot be empty")
	}

	// Check if file already exists.
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("File does not exist at path: %s", filePath)
		}
		return err
	}

	// Check if the path points to a regular file
	fileInfo, err := os.Stat(filePath)
	if err!= nil {
		return err
	}
	if !fileInfo.Mode().IsRegular() {
		return errors.New("Path does not point to a regular file")
	}

	// Check if the program has permission to access the file.
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Error opening file: %w", err)
	}
	defer file.Close()

	return nil

}

