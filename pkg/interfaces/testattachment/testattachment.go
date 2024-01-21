package testattachment

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// TestAttachment represents a file attached to a TestResult. It contains the absolute file path to the attachment file,
// a user-specified description of the attachment (which may be empty), the file type, the file size, and the time the file was created.
// This struct is used to provide a structured representation of a test attachment, which can be useful for inspecting the attachment or passing it to other functions.
type TestAttachment struct {
	// FilePath is the absolute file path to the attachment file.
	FilePath string

	// Description is the user-specified description of the attachment. It may be empty.
	Description string
	FileType    string
	FileSize    int64
	CreatedAt   time.Time
}

// NewTestAttachment creates a TestAttachment to represent a file attached to a test result. It takes a file path and a description as parameters,
// validates the file path, gets the file info, determines the file type, and returns a TestAttachment with the specified properties.
// If the file path is invalid or the file info cannot be obtained, NewTestAttachment returns an error.
func NewTestAttachment(filePath, description string) (TestAttachment, error) {
	// Validate the file path
	if err := validateFilePath(filePath); err != nil {
		return TestAttachment{}, fmt.Errorf("error creating TestAttachment: %w", err)
	}

	// Get the file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return TestAttachment{}, fmt.Errorf("error getting file info: %w", err)
	}

	// Determine the file type
	fileType := filepath.Ext(filePath)

	return TestAttachment{
		FilePath:    filePath,
		Description: description,
		FileType:    fileType,
		FileSize:    fileInfo.Size(),
		CreatedAt:   fileInfo.ModTime(),
	}, nil
}

// validateFilePath validates the given file path. It checks if the file path is not empty, if the file exists, if the path points to a regular file,
// and if the program has permission to access the file. If any of these checks fail, validateFilePath returns an error.
// This function is used to ensure that a file path is valid before creating a TestAttachment with it.
func validateFilePath(filePath string) error {
	if filePath == "" {
		return errors.New("file path cannot be empty")
	}

	// Check if file already exists.
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist at path: %s", filePath)
		}
		return err
	}

	// Check if the path points to a regular file
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	if !fileInfo.Mode().IsRegular() {
		return errors.New("path does not point to a regular file")
	}

	// Check if the program has permission to access the file.
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	return nil

}
