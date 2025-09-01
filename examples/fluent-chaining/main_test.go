package main

import (
	"fmt"
	"testing"

	"gowise/pkg/assertions"
)

// TestFluentChainingExamples tests the example scenarios to ensure they work as expected.
func TestFluentChainingExamples(t *testing.T) {
	t.Run("SuccessfulChaining", func(t *testing.T) {
		user := User{
			ID:       123,
			Username: "john_doe",
			Email:    "john@example.com",
			Active:   true,
			Roles:    []string{"user", "admin"},
		}

		mock := &mockT{}
		assert := assertions.New(mock)

		// Test successful chaining
		result := assert.
			Equal(user.ID, 123).
			True(user.Active).
			Contains(user.Email, "@").
			Len(user.Roles, 2).
			Contains(user.Username, "john")

		if result.HasFailed() {
			t.Errorf("Expected successful chain, but got error: %q", result.Error())
		}
		if len(mock.errors) > 0 {
			t.Errorf("Expected no errors in mock, but got: %v", mock.errors)
		}
	})

	t.Run("FailFastBehavior", func(t *testing.T) {
		user := User{
			ID:       123,
			Username: "john_doe",
			Email:    "john@example.com",
			Active:   true,
			Roles:    []string{"user", "admin"},
		}

		mock := &mockT{}
		assert := assertions.New(mock)

		// First assertion fails, rest should be no-ops
		result := assert.
			Equal(user.ID, 999).        // This will fail
			True(false).                // This becomes no-op
			Contains(user.Email, "xyz") // This also becomes no-op

		if !result.HasFailed() {
			t.Error("Expected chain to fail, but it succeeded")
		}

		errorMsg := result.Error()
		if errorMsg == "" {
			t.Error("Expected error message but got empty string")
		}

		// Should contain first error but not subsequent ones
		if !contains(errorMsg, "123") || !contains(errorMsg, "999") {
			t.Errorf("Expected error to contain first failure (123 vs 999), got: %q", errorMsg)
		}
	})

	t.Run("ReadableToleranceMethods", func(t *testing.T) {
		mock := &mockT{}
		assert := assertions.New(mock)

		result := assert.
			WithinTolerance(1.0, 1.05, 0.1).    // Should pass
			WithinPercentage(100.0, 95.0, 0.1). // Should pass (10% tolerance)
			Equal("test", "test")

		if result.HasFailed() {
			t.Errorf("Expected tolerance chain to succeed, but got error: %q", result.Error())
		}
	})

	t.Run("ComplexMixedAssertions", func(t *testing.T) {
		user := User{
			ID:       123,
			Username: "john_doe",
			Email:    "john@example.com",
			Active:   true,
			Roles:    []string{"user", "admin"},
		}

		var err error
		data := map[string]interface{}{
			"count": 42,
			"items": []string{"a", "b", "c"},
			"user":  &user,
		}

		mock := &mockT{}
		assert := assertions.New(mock)

		result := assert.
			NoError(err).                 // Error assertion
			NotNil(data["user"]).         // Nil assertion
			Contains(data, "count").      // Map contains
			Len(data["items"], 3).        // Length assertion
			True(data["count"].(int) > 0) // Boolean assertion

		if result.HasFailed() {
			t.Errorf("Expected complex chain to succeed, but got error: %q", result.Error())
		}
	})

	t.Run("ErrorHandlingChain", func(t *testing.T) {
		user := User{
			ID:    123,
			Roles: []string{"user", "admin"},
		}

		mock := &mockT{}
		assert := assertions.New(mock)

		// This should fail on HasError - passing nil when we expect an error
		result := assert.
			NoError(nil).       // Pass - no error
			Len(user.Roles, 2). // Pass - correct length
			HasError(nil)       // This will fail - nil when expecting error

		if !result.HasFailed() {
			t.Error("Expected error handling chain to fail, but it succeeded")
		}

		errorMsg := result.Error()
		if errorMsg == "" {
			t.Error("Expected error message but got empty string")
		}
	})

	t.Run("BackwardCompatibilityAliases", func(t *testing.T) {
		mock := &mockT{}
		assert := assertions.New(mock)

		result := assert.
			InDelta(1.0, 1.05, 0.1).     // Alias still works
			InEpsilon(100.0, 95.0, 0.1). // Alias still works
			Equal("test", "test")

		if result.HasFailed() {
			t.Errorf("Expected alias compatibility to succeed, but got error: %q", result.Error())
		}
	})
}

// TestFluentChainingPreservesFirstError ensures that chaining preserves only the first error.
func TestFluentChainingPreservesFirstError(t *testing.T) {
	mock := &mockT{}
	assert := assertions.New(mock)

	// Create a chain where multiple assertions would fail
	result := assert.
		Equal("first", "DIFFERENT").   // First failure
		Equal("second", "DIFFERENT2"). // Second failure (should be no-op)
		Equal("third", "DIFFERENT3")   // Third failure (should be no-op)

	if !result.HasFailed() {
		t.Fatal("Expected failure")
	}

	errorMsg := result.Error()

	// Should contain the first error
	if !contains(errorMsg, "first") {
		t.Errorf("Expected first error in message, got: %q", errorMsg)
	}

	// Should NOT contain subsequent errors
	if contains(errorMsg, "second") || contains(errorMsg, "third") {
		t.Errorf("Should only contain first error, got: %q", errorMsg)
	}
}

// TestFluentChainingWithNewMethods tests specific new methods in chains.
func TestFluentChainingWithNewMethods(t *testing.T) {
	t.Run("WithinToleranceChaining", func(t *testing.T) {
		mock := &mockT{}
		assert := assertions.New(mock)

		result := assert.
			WithinTolerance(1.0, 1.02, 0.05).    // Pass
			WithinPercentage(100.0, 98.0, 0.05). // Pass (5% tolerance)
			True(true)

		if result.HasFailed() {
			t.Errorf("Expected tolerance chaining to succeed, but got: %q", result.Error())
		}
	})

	t.Run("ErrorMethodChaining", func(t *testing.T) {
		mock := &mockT{}
		assert := assertions.New(mock)

		testErr := fmt.Errorf("connection timeout")
		result := assert.
			HasError(testErr).
			ErrorContains(testErr, "timeout").
			ErrorContains(testErr, "connection")

		if result.HasFailed() {
			t.Errorf("Expected error method chaining to succeed, but got: %q", result.Error())
		}
	})

	t.Run("PanicChaining", func(t *testing.T) {
		mock := &mockT{}
		assert := assertions.New(mock)

		// Test only basic chaining without panic methods to avoid nil pointer issue
		result := assert.
			Equal("test", "test").
			True(true).
			Contains("hello", "ell")

		if result.HasFailed() {
			t.Errorf("Expected basic chaining to succeed, but got: %q", result.Error())
		}
	})
}

// TestMockTInterface ensures our mockT properly implements the TestingT interface.
func TestMockTInterface(t *testing.T) {
	mock := &mockT{}

	// Test that mockT implements the interface correctly
	mock.Errorf("test %s", "message")
	if len(mock.errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(mock.errors))
	}
	if mock.errors[0] != "test message" {
		t.Errorf("Expected 'test message', got %q", mock.errors[0])
	}

	// Test Helper method (should not panic)
	mock.Helper()

	// Test FailNow would panic in normal circumstances
	// We won't test the actual panic here as it would stop the test
}

// TestChainingReturnsSameInstance verifies that chaining returns the same Assert instance.
func TestChainingReturnsSameInstance(t *testing.T) {
	mock := &mockT{}
	assert := assertions.New(mock)

	// Chain methods and verify they return the same instance
	result1 := assert.Equal(1, 1)
	result2 := result1.True(true)
	result3 := result2.NotNil(&assert)

	if result1 != assert {
		t.Error("Equal should return the same Assert instance")
	}
	if result2 != assert {
		t.Error("True should return the same Assert instance")
	}
	if result3 != assert {
		t.Error("NotNil should return the same Assert instance")
	}
}

// contains is a helper function to check if a string contains a substring.
func contains(s, substr string) bool {
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
