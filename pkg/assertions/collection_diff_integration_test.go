package assertions

import (
	"strings"
	"testing"
)

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
				"expected different length",
				"got length: 5",
				"want length: 3",
				"collection content:", // Should show what's in the collection
				"[1 2 3 4 5]",        // Formatted collection display
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
			name: "Contains with large slice - truncated display",
			setupAndAssert: func(assert *Assert) {
				got := make([]int, 50) // Large slice
				for i := range got {
					got[i] = i
				}
				want := 999
				assert.Contains(got, want)
			},
			expectErrorContains: []string{
				"expected to contain element",
				"missing from collection:",
				"999",
				"... (showing first", // Should truncate large collections
			},
		},
		{
			name: "Len with empty slice vs non-empty expectation",
			setupAndAssert: func(assert *Assert) {
				got := []string{}
				want := 2
				assert.Len(got, want)
			},
			expectErrorContains: []string{
				"expected different length",
				"got length: 0",
				"want length: 2",
				"collection is empty",
			},
		},
		{
			name: "Contains with string - enhanced character diff",
			setupAndAssert: func(assert *Assert) {
				got := "hello world"
				want := 'z' // Character not in string
				assert.Contains(got, want)
			},
			expectErrorContains: []string{
				"expected to contain character",
				"character 'z' not found in string",
				"string content: \"hello world\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a dummy testing.T that captures failures
			dummyT := &capturingT{}
			assert := New(dummyT)

			// Execute the assertion (should fail)
			tt.setupAndAssert(assert)

			errorMsg := assert.Error()
			if errorMsg == "" {
				t.Fatalf("Expected error message but got none")
			}

			// Check that all expected parts are in the error message
			for _, expected := range tt.expectErrorContains {
				if !strings.Contains(errorMsg, expected) {
					t.Errorf("Error message missing expected content %q\nFull error message:\n%s", expected, errorMsg)
				}
			}
		})
	}
}

// TestCollectionDiffFormats tests different output formats for collection diffs
func TestCollectionDiffFormats(t *testing.T) {
	tests := []struct {
		name           string
		setupAndAssert func(assert *Assert)
		expectFormat   string // "compact", "detailed", or "truncated"
	}{
		{
			name: "small slice uses compact format",
			setupAndAssert: func(assert *Assert) {
				got := []string{"a", "b", "c"}
				want := "d"
				assert.Contains(got, want)
			},
			expectFormat: "compact",
		},
		{
			name: "large slice uses truncated format",
			setupAndAssert: func(assert *Assert) {
				got := make([]int, 100)
				want := 999
				assert.Contains(got, want)
			},
			expectFormat: "truncated",
		},
		{
			name: "map with few keys uses detailed format",
			setupAndAssert: func(assert *Assert) {
				got := map[string]int{"key1": 1, "key2": 2}
				want := "missing"
				assert.Contains(got, want)
			},
			expectFormat: "detailed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dummyT := &capturingT{}
			assert := New(dummyT)

			tt.setupAndAssert(assert)

			errorMsg := assert.Error()
			if errorMsg == "" {
				t.Fatalf("Expected error message but got none")
			}

			// Verify the format matches expectation
			switch tt.expectFormat {
			case "compact":
				if !strings.Contains(errorMsg, "[") || strings.Contains(errorMsg, "...") {
					t.Errorf("Expected compact format but got: %s", errorMsg)
				}
			case "truncated":
				if !strings.Contains(errorMsg, "...") {
					t.Errorf("Expected truncated format but got: %s", errorMsg)
				}
			case "detailed":
				if !strings.Contains(errorMsg, "available keys:") {
					t.Errorf("Expected detailed format but got: %s", errorMsg)
				}
			}
		})
	}
}