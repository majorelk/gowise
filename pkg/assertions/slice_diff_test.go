package assertions

import (
	"fmt"
	"strings"
	"testing"
)

// TestSliceDiffBasic tests basic slice diff functionality with behaviour-focused testing
func TestSliceDiffBasic(t *testing.T) {
	tests := []struct {
		name                string
		got                 []int
		want                []int
		shouldPass          bool
		expectErrorContains []string
	}{
		{
			name:       "identical slices should pass",
			got:        []int{1, 2, 3},
			want:       []int{1, 2, 3},
			shouldPass: true,
		},
		{
			name:       "different slices should fail with diff",
			got:        []int{1, 2, 4},
			want:       []int{1, 2, 3},
			shouldPass: false,
			expectErrorContains: []string{
				"slices differ",
				"index 2",
				"got: 4",
				"want: 3",
			},
		},
		{
			name:       "different lengths should fail with context",
			got:        []int{1, 2},
			want:       []int{1, 2, 3},
			shouldPass: false,
			expectErrorContains: []string{
				"slices differ",
				"length",
				"got: 2",
				"want: 3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create behavioral mock that captures test behavior
			mock := &behaviorMockT{}
			assert := New(mock)

			// Execute the slice diff assertion (doesn't exist yet - will fail compilation)
			assert.SliceDiff(tt.got, tt.want)

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

// ExampleAssert_SliceDiff demonstrates proper usage of basic slice diff assertion
func ExampleAssert_SliceDiff() {
	assert := New(&silentT{})

	// Test that two integer slices are identical
	got := []int{1, 2, 3}
	want := []int{1, 2, 3}
	assert.SliceDiff(got, want)

	fmt.Println("No error:", true)
	// Output: No error: true
}
