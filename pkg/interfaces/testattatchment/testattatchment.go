package testattachment

// TestAttachment represents a file attached to a TestResult with an optional description.
type TestAttachment struct {
	// FilePath is the absolute file path to the attachment file.
	FilePath string

	// Description is the user-specified description of the attachment. It may be empty.
	Description string
}

// NewTestAttachment creates a TestAttachment to represent a file attached to a test result.
func NewTestAttachment(filePath, description string) TestAttachment {
	return TestAttachment{
		FilePath:    filePath,
		Description: description,
	}
}

