package assertions

import (
	"fmt"
	"strings"
	"testing"
)

// TestMapDiff tests map diff functionality with behaviour-focused testing
func TestMapDiff(t *testing.T) {
	tests := []struct {
		name                string
		setupAndAssert      func(assert *Assert)
		expectErrorContains []string
		shouldPass          bool
	}{
		{
			name: "identical maps should pass",
			setupAndAssert: func(assert *Assert) {
				got := map[string]int{"a": 1, "b": 2, "c": 3}
				want := map[string]int{"a": 1, "b": 2, "c": 3}
				assert.MapDiff(got, want)
			},
			shouldPass: true,
		},
		{
			name: "missing key should fail with details",
			setupAndAssert: func(assert *Assert) {
				got := map[string]int{"a": 1, "b": 2}
				want := map[string]int{"a": 1, "b": 2, "c": 3}
				assert.MapDiff(got, want)
			},
			expectErrorContains: []string{
				"maps differ",
				"missing key",
				"c",
				"expected value: 3",
			},
			shouldPass: false,
		},
		{
			name: "extra key should fail with details",
			setupAndAssert: func(assert *Assert) {
				got := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
				want := map[string]int{"a": 1, "b": 2, "c": 3}
				assert.MapDiff(got, want)
			},
			expectErrorContains: []string{
				"maps differ",
				"unexpected key",
				"d",
				"got value: 4",
			},
			shouldPass: false,
		},
		{
			name: "different values should fail with details",
			setupAndAssert: func(assert *Assert) {
				got := map[string]int{"a": 1, "b": 5, "c": 3}
				want := map[string]int{"a": 1, "b": 2, "c": 3}
				assert.MapDiff(got, want)
			},
			expectErrorContains: []string{
				"maps differ",
				"key \"b\"",
				"got: 5",
				"want: 2",
			},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create behavioral mock that captures test behavior
			mock := &behaviorMockT{}
			assert := New(mock)

			// Execute the map diff assertion
			tt.setupAndAssert(assert)

			if tt.shouldPass {
				if len(mock.errorCalls) != 0 {
					t.Errorf("Expected assertion to pass but got %d errors: %v", len(mock.errorCalls), mock.errorCalls)
				}
			} else {
				if len(mock.errorCalls) != 1 {
					t.Fatalf("Expected assertion to fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
				}

				errorMsg := mock.errorCalls[0]
				for _, expected := range tt.expectErrorContains {
					if !strings.Contains(errorMsg, expected) {
						t.Errorf("Error message missing expected content %q\nFull error message:\n%s", expected, errorMsg)
					}
				}
			}
		})
	}
}

// ExampleAssert_MapDiff demonstrates proper usage of map diff assertion
func ExampleAssert_MapDiff() {
	assert := New(&silentT{})

	// Test that two maps are identical
	got := map[string]int{"alice": 30, "bob": 25}
	want := map[string]int{"alice": 30, "bob": 25}
	assert.MapDiff(got, want)

	fmt.Println("No error:", true)
	// Output: No error: true
}
