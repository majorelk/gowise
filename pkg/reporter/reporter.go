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

// Reporter is responsible for writing test results to an output file.
// writer is an io.Writer used to write the test results.
type Reporter struct {
	writer io.Writer
}

// NewReporter creates a new Reporter that writes to the specified output file.
// outputFilePath is the path to the output file.
// The function returns a new Reporter or an error if the output file could not be created.
func NewReporter(outputFilePath string) (*Reporter, error) {
	file, err := os.Create(outputFilePath)
	if err != nil {
		return nil, fmt.Errorf("error creating output file: %w", err)
	}

	return &Reporter{
		writer: file,
	}, nil
}

// ReportTestOutput writes the JSON representation of a TestOutput to the output file.
// to is the TestOutput to report.
// The method returns an error if the TestOutput could not be written.
func (r *Reporter) ReportTestOutput(to testoutput.TestOutput) error {
	_, err := r.writer.Write([]byte(to.ToJSON() + "\n"))
	return err
}

// ReportTestMessage writes the JSON representation of a TestMessage to the output file.
// tm is the TestMessage to report.
// The method returns an error if the TestMessage could not be written.
func (r *Reporter) ReportTestMessage(tm testmessage.TestMessage) error {
	json, err := tm.ToJSON()
	if err != nil {
		return err
	}

	_, err = r.writer.Write([]byte(json + "\n"))
	return err
}

// ReportTestAttachment writes the file path and description of a TestAttachment to the output file.
// ta is the TestAttachment to report.
// The method returns an error if the TestAttachment could not be written.
func (r *Reporter) ReportTestAttachment(ta testattachment.TestAttachment) error {
	// Here you can write the file path and description of the attachment to the report.
	// If you want to include the content of the attachment in the report, you'll need to read the file and write its content to the report.
	_, err := r.writer.Write([]byte(fmt.Sprintf("Attachment: %s, Description: %s\n", ta.FilePath, ta.Description)))
	return err
}

// Close closes the output file if it implements the io.Closer interface.
// The method returns an error if the output file could not be closed.
func (r *Reporter) Close() error {
	if closer, ok := r.writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// TestReport represents a test report.
// Total is the total number of tests.
// Passed is the number of tests that passed.
// Failed is the number of tests that failed.
// Results is a slice of the results of all tests.
type TestReport struct {
	Total   int
	Passed  int
	Failed  int
	Results []teststatus.TestStatus
}

// NewTestReport creates a new TestReport.
// The function returns a new TestReport with Total, Passed, and Failed set to 0 and Results set to an empty slice.
func NewTestReport() *TestReport {
	return &TestReport{
		Total:   0,
		Passed:  0,
		Failed:  0,
		Results: []teststatus.TestStatus{},
	}
}

// AddResult adds a test result to the report.
// result is the result of a test.
// The method increments Total by 1, increments Passed by 1 if the test passed, increments Failed by 1 if the test failed, and appends the result to Results.
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
// It includes methods for reporting a TestOutput, a TestMessage, and a TestAttachment, and for closing the reporter.
type ReporterInterface interface {
	ReportTestOutput(to testoutput.TestOutput) error
	ReportTestMessage(tm testmessage.TestMessage) error
	ReportTestAttachment(ta testattachment.TestAttachment) error
	Close() error
}
