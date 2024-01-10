package reporter

import (
	"fmt"
	"io"
	"os"

	"gowise/pkg/interfaces/testattachment"
	"gowise/pkg/interfaces/testmessage"
	"gowise/pkg/interfaces/testoutput"
	"gowise/pkg/interfaces/teststatus"
)

type Reporter struct {
	writer io.Writer
}

func NewReporter(outputFilePath string) (*Reporter, error) {
	file, err := os.Create(outputFilePath)
	if err != nil {
		return nil, fmt.Errorf("error creating output file: %w", err)
	}

	return &Reporter{
		writer: file,
	}, nil
}

func (r *Reporter) ReportTestOutput(to testoutput.TestOutput) error {
	_, err := r.writer.Write([]byte(to.ToJSON() + "\n"))
	return err
}

func (r *Reporter) ReportTestMessage(tm testmessage.TestMessage) error {
	json, err := tm.ToJSON()
	if err != nil {
		return err
	}

	_, err = r.writer.Write([]byte(json + "\n"))
	return err
}

func (r *Reporter) ReportTestAttachment(ta testattachment.TestAttachment) error {
	// Here you can write the file path and description of the attachment to the report.
	// If you want to include the content of the attachment in the report, you'll need to read the file and write its content to the report.
	_, err := r.writer.Write([]byte(fmt.Sprintf("Attachment: %s, Description: %s\n", ta.FilePath, ta.Description)))
	return err
}

func (r *Reporter) Close() error {
	if closer, ok := r.writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// TestReport represents a test report.
type TestReport struct {
	Total   int
	Passed  int
	Failed  int
	Results []teststatus.TestStatus
}

// NewTestReport creates a new TestReport.
func NewTestReport() *TestReport {
	return &TestReport{
		Total:   0,
		Passed:  0,
		Failed:  0,
		Results: []teststatus.TestStatus{},
	}
}

// AddResult adds a test result to the report.
func (r *TestReport) AddResult(result teststatus.TestStatus) {
	r.Total++
	if result.GetResult() == teststatus.Passed.GetResult() {
		r.Passed++
	} else {
		r.Failed++
	}
	r.Results = append(r.Results, result)
}

// ReporterInterface represents the interface for a reporter.
type ReporterInterface interface {
	ReportTestOutput(to testoutput.TestOutput) error
	ReportTestMessage(tm testmessage.TestMessage) error
	ReportTestAttachment(ta testattachment.TestAttachment) error
	Close() error
}
