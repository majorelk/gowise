package testattachment

import (
	"testing"
	"os"
	"io/ioutil"
)

func TestNewTestAttachment(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "attachment_*.txt")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	attachment, err := NewTestAttachment(tempFile.Name(), "This is a test attachment.")

	if err != nil {
		t.Errorf("Error creating TestAttachment: %v", err)
		return
	}

	if attachment.FilePath != tempFile.Name() {
		t.Errorf("Expected FilePath to be %s, but got %s", tempFile.Name(), attachment.FilePath)
	}

	if attachment.Description != "This is a test attachment." {
		t.Errorf("Expected Description to be 'This is a test attachment.', but got %s", attachment.Description)
	}
}


