// Package teststatus provides enums and interfaces for test statuses.
package teststatus

import (
	"testing"
	"gowise/pkg/assertions" // Import the assertions package
)

// TestResultGetResult tests the GetResult method of the Result type.
func TestResultGetResult(t *testing.T) {
	// Test for Passed result
	passed := Result(Passed)
	assert := assertions.New(t)
	assert.Equal("Passed", passed.GetResult())

	// Test for Failed result
	failed := Result(Failed)
	assert.Equal("Failed", failed.GetResult())

	// Test for Unknown result
	unknown := Result(42) // Some unknown result
	assert.Equal("Unknown", unknown.GetResult())
}

// Additional tests for other functionality in the teststatus package can be added similarly.

