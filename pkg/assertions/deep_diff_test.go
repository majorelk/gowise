package assertions

import (
	"fmt"
	"strings"
	"testing"
)

// behaviorMockT is defined in assertions_passing_test.go - shared across test files

// Test types for deep diff testing
type Contact struct {
	Name  string
	Email string
	Tags  []string
}

type Company struct {
	Name      string
	Employees []Contact
	Locations map[string]string
}

// TestDeepDiff tests universal deep diff functionality with behaviour-focused testing
func TestDeepDiff(t *testing.T) {
	tests := []struct {
		name                string
		setupAndAssert      func(assert *Assert)
		expectErrorContains []string
		shouldPass          bool
	}{
		{
			name: "identical primitives should pass",
			setupAndAssert: func(assert *Assert) {
				got := "hello"
				want := "hello"
				assert.DeepDiff(got, want)
			},
			shouldPass: true,
		},
		{
			name: "different primitives should fail",
			setupAndAssert: func(assert *Assert) {
				got := 42
				want := 24
				assert.DeepDiff(got, want)
			},
			expectErrorContains: []string{
				"values differ",
				"got: 42",
				"want: 24",
			},
			shouldPass: false,
		},
		{
			name: "slice differences should use SliceDiffGeneric",
			setupAndAssert: func(assert *Assert) {
				got := []string{"a", "b", "x"}
				want := []string{"a", "b", "c"}
				assert.DeepDiff(got, want)
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
			name: "map differences should use MapDiff",
			setupAndAssert: func(assert *Assert) {
				got := map[string]int{"a": 1, "b": 5}
				want := map[string]int{"a": 1, "b": 2}
				assert.DeepDiff(got, want)
			},
			expectErrorContains: []string{
				"maps differ",
				"key \"b\"",
				"got: 5",
				"want: 2",
			},
			shouldPass: false,
		},
		{
			name: "struct differences should use StructDiff",
			setupAndAssert: func(assert *Assert) {
				got := Contact{Name: "Alice", Email: "alice@old.com", Tags: []string{"dev"}}
				want := Contact{Name: "Alice", Email: "alice@new.com", Tags: []string{"dev"}}
				assert.DeepDiff(got, want)
			},
			expectErrorContains: []string{
				"structs differ",
				"field \"Email\"",
				"got: alice@old.com",
				"want: alice@new.com",
			},
			shouldPass: false,
		},
		{
			name: "complex nested structures should provide detailed diff",
			setupAndAssert: func(assert *Assert) {
				got := Company{
					Name: "TechCorp",
					Employees: []Contact{
						{Name: "Alice", Email: "alice@tech.com", Tags: []string{"dev"}},
					},
					Locations: map[string]string{"HQ": "London"},
				}
				want := Company{
					Name: "TechCorp",
					Employees: []Contact{
						{Name: "Bob", Email: "alice@tech.com", Tags: []string{"dev"}},
					},
					Locations: map[string]string{"HQ": "London"},
				}
				assert.DeepDiff(got, want)
			},
			expectErrorContains: []string{
				"structs differ",
				"field \"Employees\"",
			},
			shouldPass: false,
		},
		{
			name: "different types should fail clearly",
			setupAndAssert: func(assert *Assert) {
				got := []string{"a", "b"}
				want := map[string]int{"a": 1}
				assert.DeepDiff(got, want)
			},
			expectErrorContains: []string{
				"types differ",
				"[]string",
				"map[string]int",
			},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test behavior using proper mock
			mock := &behaviorMockT{}
			assert := New(mock)

			// Execute the deep diff assertion
			tt.setupAndAssert(assert)

			// Get error message from mock's captured behavior
			var errorMsg string
			if len(mock.errorCalls) > 0 {
				errorMsg = mock.errorCalls[0]
			}

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

// silentT is defined in assertions_passing_test.go - shared across test files

// ExampleAssert_DeepDiff demonstrates proper usage of universal deep diff assertion
func ExampleAssert_DeepDiff() {
	assert := New(&silentT{})

	// Test that two complex structures are identical
	type Person struct {
		Name string
		Tags []string
	}

	got := Person{Name: "Alice", Tags: []string{"dev", "lead"}}
	want := Person{Name: "Alice", Tags: []string{"dev", "lead"}}
	assert.DeepDiff(got, want)

	fmt.Println("No error:", true)
	// Output: No error: true
}
