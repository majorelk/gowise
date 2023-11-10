// testoutput_test.go
package testoutput

import (
	"strings"
	"testing"
)


var expectedJSON = `{
  "text": "Hello, World!",
  "stream": "stdout",
  "testid": "test123",
  "testname": "ExampleTest"
}`

 // TestToJSON Function
func TestToJSON(t *testing.T) {

	output := NewTestOutput("Hello, World!", "stdout", "test123", "ExampleTest")
	actualJSON := output.ToJSON()

	if strings.TrimSpace(actualJSON) != strings.TrimSpace(expectedJSON) {
		t.Errorf("Expected JSON:\n%s\n\nActual JSON:\n%s", expectedJSON, actualJSON)
	}
}

// TestToJSONWriter Function
func TestToJSONWriter(t *testing.T) {

	output := NewTestOutput("Hello, World!", "stdout", "test123", "ExampleTest")

	var actualBuilder strings.Builder
	err := output.ToJSONWriter(&actualBuilder)
	if err != nil {
		t.Errorf("Error writing to JSON writer: %v", err)
		return
	}

	actualJSON := actualBuilder.String()

	if strings.TrimSpace(actualJSON) != strings.TrimSpace(expectedJSON) {
		t.Errorf("Expected JSON:\n%s\n\nActual JSON:\n%s", expectedJSON, actualJSON)
	}
}

