package testattachment

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewTestAttachment(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "attachment_*.txt")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	description := "This is a test attachment."
	attachment, err := NewTestAttachment(tempFile.Name(), description)

	if err != nil {
		t.Errorf("Error creating TestAttachment: %v", err)
		return
	}

	if attachment.FilePath != tempFile.Name() {
		t.Errorf("Expected FilePath to be %s, but got %s", tempFile.Name(), attachment.FilePath)
	}

	if attachment.Description != description {
		t.Errorf("Expected Description to be '%s', but got %s", description, attachment.Description)
	}

	if attachment.FileType != filepath.Ext(tempFile.Name()) {
		t.Errorf("Expected FileType to be %s, but got %s", filepath.Ext(tempFile.Name()), attachment.FileType)
	}

	fileInfo, _ := os.Stat(tempFile.Name())
	if attachment.FileSize != fileInfo.Size() {
		t.Errorf("Expected FileSize to be %d, but got %d", fileInfo.Size(), attachment.FileSize)
	}

	if attachment.CreatedAt != fileInfo.ModTime() {
		t.Errorf("Expected CreatedAt to be %v, but got %v", fileInfo.ModTime(), attachment.CreatedAt)
	}
}
