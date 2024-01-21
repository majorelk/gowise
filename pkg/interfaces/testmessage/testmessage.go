package testmessage

import (
	"encoding/json"
	"fmt"
)

// TestMessage represents a message sent by a test to all listeners. It contains the destination
// of the message, the message itself, and an optional test ID. This struct is used to provide a
// structured representation of a test message, which can be useful for inspecting the message or
// passing it to other functions.
type TestMessage struct {
	Destination string `json:"destination"`
	Message     string `json:"message"`
	TestID      string `json:"testId,omitempty"`
}

// NewTestMessage constructs a TestMessage with the given destination, message, and test ID.
// The destination parameter is the destination of the message, the message parameter is the
// message itself, and the testID parameter is an optional test ID. This function is used to
// create a new TestMessage instance with the specified properties.
func NewTestMessage(destination, message, testID string) TestMessage {
	return TestMessage{
		Destination: destination,
		Message:     message,
		TestID:      testID,
	}
}

// ToString converts the TestMessage object to a string in the format "Destination: Message".
// This method is used to serialize the TestMessage to a human-readable string format, which
// can be useful for logging or displaying the TestMessage.
func (tm TestMessage) ToString() string {
	return fmt.Sprintf("%s: %s", tm.Destination, tm.Message)
}

// ToJSON converts the TestMessage object to a JSON string. If the TestMessage cannot be
// marshaled to JSON, ToJSON returns an error. This method is used to serialize the TestMessage
// to JSON format, which can be useful for storing the TestMessage or sending it over a network.
func (tm TestMessage) ToJSON() (string, error) {
	data, err := json.Marshal(tm)
	if err != nil {
		return "", fmt.Errorf("error marshalling TestMessage to JSON: %w", err)
	}
	return string(data), nil
}
