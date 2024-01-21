// Package testoutput provides a struct and methods for representing test output.
package testoutput

import (
	"encoding/json"
	"errors"
	"io"
)

// TestOutput represents a unit of output from a test to a specific output stream.
// It contains the text of the output, the stream to which the output was written,
// the ID and name of the test, and the status of the test. This struct is used to
// provide a structured representation of test output, which can be useful for
// inspecting the output or passing it to other functions.
type TestOutput struct {
	Text     string `json:"text"`
	Stream   string `json:"stream"`
	TestID   string `json:"testid,omitempty"`
	TestName string `json:"testname,omitempty"`
	Status   string `json:"status,omitempty"`
}

// NewTestOutput constructs a TestOutput with the given text, stream, test ID, test name, and status.
// The text parameter is the text of the output, the stream parameter is the stream to which the output
// was written, the testID parameter is the ID of the test, the testName parameter is the name of the test,
// and the status parameter is the status of the test. This function is used to create a new TestOutput
// instance with the specified properties.
func NewTestOutput(text, stream, testID, testName string, status string) TestOutput {
	return TestOutput{
		Text:     text,
		Stream:   stream,
		TestID:   testID,
		TestName: testName,
		Status:   status,
	}
}

// WithText sets the Text field of the TestOutput to the given text and returns the TestOutput.
// This method is used to change the text of the TestOutput after it has been created.
func (to *TestOutput) WithText(text string) *TestOutput {
	to.Text = text
	return to
}

// ToJSON converts the TestOutput object to a JSON string and returns the string.
// This method is used to serialize the TestOutput to JSON format, which can be useful
// for storing the TestOutput or sending it over a network.
func (to TestOutput) ToJSON() string {
	data, _ := json.MarshalIndent(to, "", "  ")
	return string(data)
}

// ToJSONWriter writes the TestOutput object to the given io.Writer in JSON format.
// This method is used to serialize the TestOutput to JSON format and write it to an
// io.Writer, which can be useful for writing the TestOutput to a file or a network connection.
func (to TestOutput) ToJSONWriter(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(to); err != nil {
		return errors.New("Error encoding TestOutput to JSON: " + err.Error())
	}
	return nil
}
