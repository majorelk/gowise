package assertions

import (
	"strings"
	"testing"
)

// TestSliceDiffGeneric tests generic slice diff functionality with behaviour-focused testing
func TestSliceDiffGeneric(t *testing.T) {
	tests := []struct {
		name                string
		setupAndAssert      func(assert *Assert)
		expectErrorContains []string
		shouldPass          bool
	}{
		{
			name: "string slices - identical should pass",
			setupAndAssert: func(assert *Assert) {
				got := []string{"a", "b", "c"}
				want := []string{"a", "b", "c"}
				assert.SliceDiffGeneric(got, want)
			},
			shouldPass: true,
		},
		{
			name: "string slices - different elements should fail",
			setupAndAssert: func(assert *Assert) {
				got := []string{"a", "b", "x"}
				want := []string{"a", "b", "c"}
				assert.SliceDiffGeneric(got, want)
			},
			expectErrorContains: []string{
				"slices differ",
				"index 2",
				"got: x",
				"want: c",
			},
			shouldPass: false,
		},
		{
			name: "float64 slices - different lengths should fail",
			setupAndAssert: func(assert *Assert) {
				got := []float64{1.1, 2.2}
				want := []float64{1.1, 2.2, 3.3}
				assert.SliceDiffGeneric(got, want)
			},
			expectErrorContains: []string{
				"slices differ",
				"length",
				"got: 2",
				"want: 3",
			},
			shouldPass: false,
		},
		{
			name: "boolean slices - different elements should fail",
			setupAndAssert: func(assert *Assert) {
				got := []bool{true, false, true}
				want := []bool{true, false, false}
				assert.SliceDiffGeneric(got, want)
			},
			expectErrorContains: []string{
				"slices differ",
				"index 2",
				"got: true",
				"want: false",
			},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a dummy testing.T that captures failures
			dummyT := &capturingT{}
			assert := New(dummyT)

			// Execute the generic slice diff assertion
			tt.setupAndAssert(assert)

			errorMsg := assert.Error()

			if tt.shouldPass {
				if errorMsg != "" {
					t.Errorf("Expected assertion to pass but got error: %s", errorMsg)
				}
			} else {
				if errorMsg == "" {
					t.Fatalf("Expected error message but got none")
				}

				for _, expected := range tt.expectErrorContains {
					if !strings.Contains(errorMsg, expected) {
						t.Errorf("Error message missing expected content %q\nFull error message:\n%s", expected, errorMsg)
					}
				}
			}
		})
	}
}