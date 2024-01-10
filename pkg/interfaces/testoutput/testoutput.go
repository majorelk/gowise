// Package testoutput provides a struct and methods for representing test output.
package testoutput

import (
	"encoding/json"
	"errors"
	"io"
)

// TestOutput holds a unit of output from a test to a specific output stream
type TestOutput struct {
	Text     string `json:"text"`
	Stream   string `json:"stream"`
	TestID   string `json:"testid,omitempty"`
	TestName string `json:"testname,omitempty"`
	Status   string `json:"status,omitempty"`
}

// NewTestOutput constructs a TestOutput with the given text, stream, test ID, and test name
func NewTestOutput(text, stream, testID, testName string, status string) TestOutput {
	return TestOutput{
		Text:     text,
		Stream:   stream,
		TestID:   testID,
		TestName: testName,
		Status:   status,
	}
}

func (to *TestOutput) WithText(text string) *TestOutput {
	to.Text = text
	return to
}

// ToJSON converts the TestOutput object to a JSON string
func (to TestOutput) ToJSON() string {
	data, _ := json.MarshalIndent(to, "", "  ")
	return string(data)
}

// ToJSONWriter writes the TestOutput object to a JSON writer
func (to TestOutput) ToJSONWriter(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(to); err != nil {
		return errors.New("Error encoding TestOutput to JSON: " + err.Error())
	}
	return nil
}
