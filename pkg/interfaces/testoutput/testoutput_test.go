// testoutput_test.go

package testoutput

import (
	"encoding/json"
	"testing"
)

func TestToJSON(t *testing.T) {
	expectedJSON := `{
		"text": "Hello, World!",
		"stream": "stdout",
		"testid": "test123",
		"testname": "ExampleTest"
	}`

	output := NewTestOutput("Hello, World!", "stdout", "test123", "ExampleTest")
	actualJSON := output.ToJSON()

	if actualJSON != expectedJSON {
		t.Errorf("Expected JSON:\n%s\n\nActual JSON:\n%s", expectedJSON, actualJSON)
	}
}

func TestToJSONWriter(t *testing.T) {
	expectedJSON := `{
		"text": "Hello, World!",
		"stream": "stdout",
		"testid": "test123",
		"testname": "ExampleTest"
	}`

	output := NewTestOutput("Hello, World!", "stdout", "test123", "ExampleTest")

	var actualBuilder strings.Builder
	err := output.ToJSONWriter(&actualBuilder)
	if err != nil {
		t.Errorf("Error writing to JSON writer: %v", err)
		return
	}

	actualJSON := actualBuilder.String()

	if actualJSON != expectedJSON {
		t.Errorf("Expected JSON:\n%s\n\nActual JSON:\n%s", expectedJSON, actualJSON)
	}
}

