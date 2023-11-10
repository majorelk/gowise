package testmessage

import (
	"encoding/json"
	"testing"
)

func TestToJSON(t *testing.T) {
	expectedJSON := `{"destination":"Console","message":"Hello, World!","testId":"test123"}`

	testMessage := NewTestMessage("Console", "Hello, World!", "test123")
	actualJSON, err := testMessage.ToJSON()
	if err != nil {
		t.Errorf("Error converting TestMessage to JSON: %v", err)
		return
	}

	if actualJSON != expectedJSON {
		t.Errorf("Expected JSON:\n%s\n\nActual JSON:\n%s", expectedJSON, actualJSON)
	}
}

func TestToJSONWithoutTestID(t *testing.T) {
	expectedJSON := `{"destination":"Console","message":"Hello, World!"}`

	testMessage := NewTestMessage("Console", "Hello, World!", "")
	actualJSON, err := testMessage.ToJSON()
	if err != nil {
		t.Errorf("Error converting TestMessage to JSON: %v", err)
		return
	}

	if actualJSON != expectedJSON {
		t.Errorf("Expected JSON:\n%s\n\nActual JSON:\n%s", expectedJSON, actualJSON)
	}
}

func TestToString(t *testing.T) {
	expectedString := "Console: Hello, World!"

	testMessage := NewTestMessage("Console", "Hello, World!", "test123")
	actualString := testMessage.ToString()

	if actualString != expectedString {
		t.Errorf("Expected string: %s\nActual string: %s", expectedString, actualString)
	}
}

