package assertions

import (
	"testing"
)

// TestFluentChaining tests the method chaining functionality.
func TestFluentChaining(t *testing.T) {
	t.Run("SuccessfulChaining", func(t *testing.T) {
		mock := &mockT{}
		assert := New(mock)

		// Chain multiple successful assertions
		result := assert.Equal(1, 1).True(true).NotNil(&assert).Equal("hello", "hello")

		if result != assert {
			t.Error("Chaining should return the same Assert instance")
		}
		if result.HasFailed() {
			t.Errorf("Expected no failure, but HasFailed() returned true")
		}
		if result.Error() != "" {
			t.Errorf("Expected no error message, got: %q", result.Error())
		}
	})

	t.Run("FailFastBehavior", func(t *testing.T) {
		mock := &mockT{}
		assert := New(mock)

		// First assertion fails, subsequent should be no-ops
		result := assert.Equal(1, 2).True(false).Equal(3, 4)

		if !result.HasFailed() {
			t.Error("Expected failure, but HasFailed() returned false")
		}

		errorMsg := result.Error()
		if errorMsg == "" {
			t.Error("Expected error message but got empty string")
		}

		// Should contain the first error (1 != 2), not subsequent ones
		if !chainTestContainsString(errorMsg, "1") || !chainTestContainsString(errorMsg, "2") {
			t.Errorf("Expected error to contain first failure (1 vs 2), got: %q", errorMsg)
		}

		// Should NOT contain errors from subsequent assertions
		if chainTestContainsString(errorMsg, "3") || chainTestContainsString(errorMsg, "4") {
			t.Errorf("Error message should not contain subsequent failed assertions, got: %q", errorMsg)
		}
	})

	t.Run("ChainWithMixedMethods", func(t *testing.T) {
		mock := &mockT{}
		assert := New(mock)

		var nilPtr *int
		slice := []string{"a", "b", "c"}

		// Mix different assertion types in chain
		result := assert.NotNil(&slice).Nil(nilPtr).Contains(slice, "b").Len(slice, 3).NoError(nil)

		if result.HasFailed() {
			t.Errorf("Expected success, but got error: %q", result.Error())
		}
	})

	t.Run("ChainFailureInMiddle", func(t *testing.T) {
		mock := &mockT{}
		assert := New(mock)

		slice := []string{"a", "b", "c"}

		// Second assertion fails, rest should be no-ops
		result := assert.Len(slice, 3).Contains(slice, "missing").True(true)

		if !result.HasFailed() {
			t.Error("Expected failure but chain succeeded")
		}

		errorMsg := result.Error()
		// Should contain error about missing element, not about boolean assertion
		if !chainTestContainsString(errorMsg, "missing") {
			t.Errorf("Expected error about missing element, got: %q", errorMsg)
		}
	})
}

// TestReturnTypeConsistency verifies all assertion methods return *Assert for chaining.
func TestReturnTypeConsistency(t *testing.T) {
	mock := &mockT{}
	assert := New(mock)

	// Test that all updated methods return *Assert (spot check key ones)
	var result *Assert

	result = assert.Equal(1, 1)
	if result == nil {
		t.Error("Equal should return *Assert")
	}

	result = assert.True(true)
	if result == nil {
		t.Error("True should return *Assert")
	}

	result = assert.NotNil(&assert)
	if result == nil {
		t.Error("NotNil should return *Assert")
	}

	result = assert.Contains([]int{1, 2, 3}, 2)
	if result == nil {
		t.Error("Contains should return *Assert")
	}

	result = assert.Len("hello", 5)
	if result == nil {
		t.Error("Len should return *Assert")
	}

	result = assert.NoError(nil)
	if result == nil {
		t.Error("NoError should return *Assert")
	}
}

// TestChainingPreservesFirstError verifies that only the first error is preserved in chaining.
func TestChainingPreservesFirstError(t *testing.T) {
	mock := &mockT{}
	assert := New(mock)

	// Create a chain where multiple assertions fail
	result := assert.Equal("first", "FIRST").Equal("second", "SECOND").Equal("third", "THIRD")

	if !result.HasFailed() {
		t.Fatal("Expected failure")
	}

	errorMsg := result.Error()

	// Should contain the first error
	if !chainTestContainsString(errorMsg, "first") || !chainTestContainsString(errorMsg, "FIRST") {
		t.Errorf("Expected first error in message, got: %q", errorMsg)
	}

	// Should NOT contain subsequent errors
	if chainTestContainsString(errorMsg, "second") || chainTestContainsString(errorMsg, "SECOND") ||
		chainTestContainsString(errorMsg, "third") || chainTestContainsString(errorMsg, "THIRD") {
		t.Errorf("Should only contain first error, got: %q", errorMsg)
	}
}

// chainTestContainsString checks if a string contains a substring (helper function).
func chainTestContainsString(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) &&
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
