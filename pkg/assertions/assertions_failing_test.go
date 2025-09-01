// TestAssertions_Failing contains intentionally failing tests.

//go:build failing_tests
// +build failing_tests

package assertions

import (
	"fmt"
	"testing"
)

// behaviorMockT is defined in assertions_passing_test.go - shared across test files

func TestAssertions_Failing(t *testing.T) {
	t.Run("Equal", func(t *testing.T) {
		assert := New(t)

		// Test cases
		testCases := []struct {
			expected, actual interface{}
			pass             bool
		}{
			{42, 23, true},                       // different values int should fail
			{"hello", "world", true},             // different string should fail
			{false, true, true},                  // differet bool should fail
			{"hello", "hello", false},            // same value marked fals should fail
			{true, false, true},                  // opposite different bool should fail
			{42, "42", false},                    // Different types
			{[]int{1, 2, 3}, []int{1, 2}, false}, // Different slice lengths
			{struct{ X, Y int }{1, 2}, struct{ X, Y int }{1}, false}, // Different struct fields
			{map[string]int{"a": 1}, map[string]int{"b": 2}, false},  // Different map key-value pairs
		}

		for i, tc := range testCases {
			t.Run(fmt.Sprintf("Test case %d", i+1), func(t *testing.T) {
				// Behavioral test: verify framework behavior through TestingT interface
				mock := &behaviorMockT{}
				assert := New(mock)
				assert.Equal(tc.expected, tc.actual)

				if tc.pass {
					// Expected failure: assertion should call TestingT.Errorf exactly once
					if len(mock.errorCalls) != 1 {
						t.Errorf("Test case %d: expected assertion failure (1 Errorf call), got %d: %v", i+1, len(mock.errorCalls), mock.errorCalls)
					}
				} else {
					// Expected success: assertion should call no TestingT methods
					if len(mock.errorCalls) != 0 {
						t.Errorf("Test case %d: expected assertion success (no Errorf calls), got %d: %v", i+1, len(mock.errorCalls), mock.errorCalls)
					}
				}
			})
		}
	})

	t.Run("NotEqual", func(t *testing.T) {
		assert := New(t)

		// Test cases for NotEqual assertion
		testCases := []struct {
			expected, actual interface{}
			pass             bool
		}{
			{42, 42, true},                   // Identical values
			{[]int{1, 2}, []int{1, 2}, true}, // Identical slices
			{struct{ X, Y int }{1, 2}, struct{ X, Y int }{1, 2}, true}, // Identical structs
			{map[string]int{"a": 1}, map[string]int{"a": 1}, true},     // Identical maps
		}

		for i, tc := range testCases {
			t.Run(fmt.Sprintf("Test case %d", i+1), func(t *testing.T) {
				// Behavioral test: verify framework behavior through TestingT interface
				mock := &behaviorMockT{}
				assert := New(mock)
				assert.NotEqual(tc.expected, tc.actual)

				if tc.pass {
					// Expected failure: assertion should call TestingT.Errorf exactly once
					if len(mock.errorCalls) != 1 {
						t.Errorf("Test case %d: expected assertion failure (1 Errorf call), got %d: %v", i+1, len(mock.errorCalls), mock.errorCalls)
					}
				} else {
					// Expected success: assertion should call no TestingT methods
					if len(mock.errorCalls) != 0 {
						t.Errorf("Test case %d: expected assertion success (no Errorf calls), got %d: %v", i+1, len(mock.errorCalls), mock.errorCalls)
					}
				}
			})
		}
	})

	t.Run("True", func(t *testing.T) {
		assert := New(t)

		// Test cases for True assertion
		testCases := []struct {
			value bool
			pass  bool
		}{
			{true, false}, // Expected false, but actual is true
			{false, true}, // Expected true, but actual is false
			{0, true},     // Expected true, but actual is 0
			{"", true},    // Expected true, but actual is an empty string
		}

		for i, tc := range testCases {
			t.Run(fmt.Sprintf("Test case %d", i+1), func(t *testing.T) {
				// Behavioral test: verify framework behavior through TestingT interface
				mock := &behaviorMockT{}
				assert := New(mock)
				assert.True(tc.value)

				if tc.pass {
					// Expected failure: assertion should call TestingT.Errorf exactly once
					if len(mock.errorCalls) != 1 {
						t.Errorf("Test case %d: expected assertion failure (1 Errorf call), got %d: %v", i+1, len(mock.errorCalls), mock.errorCalls)
					}
				} else {
					// Expected success: assertion should call no TestingT methods
					if len(mock.errorCalls) != 0 {
						t.Errorf("Test case %d: expected assertion success (no Errorf calls), got %d: %v", i+1, len(mock.errorCalls), mock.errorCalls)
					}
				}
			})
		}
	})

	t.Run("False", func(t *testing.T) {
		assert := New(t)

		// Test cases for False assertion
		testCases := []struct {
			value bool
			pass  bool
		}{
			{true, true},    // Expected false, but actual is true
			{1, true},       // Expected false, but actual is 1
			{"hello", true}, // Expected false, but actual is a non-empty string
		}

		for i, tc := range testCases {
			t.Run(fmt.Sprintf("Test case %d", i+1), func(t *testing.T) {
				// Behavioral test: verify framework behavior through TestingT interface
				mock := &behaviorMockT{}
				assert := New(mock)
				assert.False(tc.value)

				if tc.pass {
					// Expected failure: assertion should call TestingT.Errorf exactly once
					if len(mock.errorCalls) != 1 {
						t.Errorf("Test case %d: expected assertion failure (1 Errorf call), got %d: %v", i+1, len(mock.errorCalls), mock.errorCalls)
					}
				} else {
					// Expected success: assertion should call no TestingT methods
					if len(mock.errorCalls) != 0 {
						t.Errorf("Test case %d: expected assertion success (no Errorf calls), got %d: %v", i+1, len(mock.errorCalls), mock.errorCalls)
					}
				}
			})
		}
	})
}
