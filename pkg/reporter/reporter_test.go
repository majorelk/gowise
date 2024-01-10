package reporter

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"gowise/pkg/interfaces/testattachment"
	"gowise/pkg/interfaces/testmessage"
	"gowise/pkg/interfaces/testoutput"
)

func TestNewReporter(t *testing.T) {
	filePath := "test_report.txt"
	defer os.Remove(filePath)

	r, err := NewReporter(filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer r.Close()

	if r.writer == nil {
		t.Fatal("Expected writer to be initialized, got nil")
	}
}

func TestReportTestOutput(t *testing.T) {
	filePath := "test_report.txt"
	defer os.Remove(filePath)

	r, _ := NewReporter(filePath)
	defer r.Close()

	to := testoutput.NewTestOutput("text", "stream", "testID", "testName", "status")
	err := r.ReportTestOutput(to)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	content, _ := ioutil.ReadFile(filePath)
	if !strings.Contains(string(content), to.ToJSON()) {
		t.Fatalf("Expected file to contain test output, got %s", content)
	}
}

func TestReportTestMessage(t *testing.T) {
	filePath := "test_report.txt"
	defer os.Remove(filePath)

	r, _ := NewReporter(filePath)
	defer r.Close()

	tm := testmessage.NewTestMessage("destination", "message", "testID")
	err := r.ReportTestMessage(tm)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	content, _ := ioutil.ReadFile(filePath)
	json, _ := tm.ToJSON()
	if !strings.Contains(string(content), json) {
		t.Fatalf("Expected file to contain test message, got %s", content)
	}
}

func TestReportTestAttachment(t *testing.T) {
	filePath := "test_report.txt"
	defer os.Remove(filePath)

	r, _ := NewReporter(filePath)
	defer r.Close()

	ta, _ := testattachment.NewTestAttachment("filePath", "description")
	err := r.ReportTestAttachment(ta)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	content, _ := ioutil.ReadFile(filePath)
	expected := fmt.Sprintf("Attachment: %s, Description: %s\n", ta.FilePath, ta.Description)
	if !strings.Contains(string(content), expected) {
		t.Fatalf("Expected file to contain test attachment, got %s", content)
	}
}
