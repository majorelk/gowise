package testmessage

import (
	"encoding/json"
	"fmt"
)

// TestMessage holds a message sent by a test to all listeners
type TestMessage struct {
	Destination string `json:"destination"`
	Message     string `json:"message"`
	TestID      string `json:"testId,omitempty"`
}

// NewTestMessage constructs a TestMessage with the given destination, message, and test ID
func NewTestMessage(destination, message, testID string) TestMessage {
	return TestMessage{
		Destination: destination,
		Message:     message,
		TestID:      testID,
	}
}

// ToString converts TestMessage object to string
func (tm TestMessage) ToString() string {
	return fmt.Sprintf("%s: %s", tm.Destination, tm.Message)
}

// ToJSON converts TestMessage object to JSON string
func (tm TestMessage) ToJSON() (string, error) {
	data, err := json.Marshal(tm)
	if err != nil {
		log.Printf("Error marshalling TestMessage to JSON: %v", err)
		return "", err
	}
	return string(data), nil
}

