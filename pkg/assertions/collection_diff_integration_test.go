package assertions

import (
	"strings"
	"testing"
)

// behaviorMockT is defined in assertions_passing_test.go - shared across test files

// TestCollectionDiffIntegration tests that collection assertions provide enhanced error messages
// through the public API, testing observable behaviour rather than internal implementation.
func TestCollectionDiffIntegration(t *testing.T) {
	tests := []struct {
		name                string
		setupAndAssert      func(assert *Assert)
		expectErrorContains []string
	}{
		{
			name: "Contains with slice - missing element shows diff",
			setupAndAssert: func(assert *Assert) {
				got := []string{"apple", "banana", "cherry"}
				want := "orange"
				assert.Contains(got, want)
			},
			expectErrorContains: []string{
				"expected to contain element",
				"missing from collection:",
				"orange",
				"apple", "banana", "cherry", // Should show what's actually in the collection
			},
		},
		{
			name: "Len with slice - wrong length shows collection content",
			setupAndAssert: func(assert *Assert) {
				got := []int{1, 2, 3, 4, 5}
				want := 3
				assert.Len(got, want)
			},
			expectErrorContains: []string{
				"got length: 5",
				"want length: 3",
				"collection content:", // Should show what's in the collection
				"[1 2 3 4 5]",         // Formatted collection display
			},
		},
		{
			name: "Contains with map - missing key shows available keys",
			setupAndAssert: func(assert *Assert) {
				got := map[string]int{"foo": 1, "bar": 2, "baz": 3}
				want := "missing_key"
				assert.Contains(got, want)
			},
			expectErrorContains: []string{
				"expected to contain key",
				"missing from map:",
				"missing_key",
				"available keys:", // Should show what keys exist
				"foo", "bar", "baz",
			},
		},
		{
			name: "Contains with string - missing substring shows string content",
			setupAndAssert: func(assert *Assert) {
				got := "hello world"
				want := "goodbye"
				assert.Contains(got, want)
			},
			expectErrorContains: []string{
				"expected to contain substring",
				"substring \"goodbye\" not found in string",
				"string content: \"hello world\"",
			},
		},
		{
			name: "Len with empty slice shows empty collection",
			setupAndAssert: func(assert *Assert) {
				got := []string{}
				want := 3
				assert.Len(got, want)
			},
			expectErrorContains: []string{
				"got length: 0",
				"want length: 3",
				"collection is empty",
			},
		},
		{
			name: "Contains with empty slice shows empty collection",
			setupAndAssert: func(assert *Assert) {
				got := []int{}
				want := 42
				assert.Contains(got, want)
			},
			expectErrorContains: []string{
				"expected to contain element",
				"missing from collection: 42",
				"collection is empty",
			},
		},
		{
			name: "Len with large slice shows truncated content",
			setupAndAssert: func(assert *Assert) {
				got := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
				want := 5
				assert.Len(got, want)
			},
			expectErrorContains: []string{
				"got length: 15",
				"want length: 5",
				"collection content:",
				"... (showing first 10 elements)", // Should indicate truncation
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a dummy testing.T that captures failures
			mock := &behaviorMockT{}
			assert := New(mock)

			// Execute the assertion
			tt.setupAndAssert(assert)

			// Framework behavior: FAIL = exactly 1 Errorf call with expected content
			if len(mock.errorCalls) != 1 {
				t.Fatalf("Expected assertion to fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
			}

			errorMsg := mock.errorCalls[0]

			// Check that all expected parts are in the error message
			for _, expected := range tt.expectErrorContains {
				if !strings.Contains(errorMsg, expected) {
					t.Errorf("Error message missing expected content %q\nFull error message:\n%s", expected, errorMsg)
				}
			}
		})
	}
}
