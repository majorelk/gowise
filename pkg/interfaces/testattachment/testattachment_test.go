package testattachment

import "testing"

func TestNewTestAttachment(t *testing.T) {
	filePath := "/path/to/attachment.txt"
	description := "This is a test attachment."

	attachment := NewTestAttachment(filePath, description)

	if attachment.FilePath != filePath {
		t.Errorf("Expected FilePath to be %s, but got %s", filePath, attachment.FilePath)
	}

	if attachment.Description != description {
		t.Errorf("Expected Description to be %s, but got %s", description, attachment.Description)
	}
}

