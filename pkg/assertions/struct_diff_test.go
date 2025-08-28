package assertions

import (
	"fmt"
	"strings"
	"testing"
)

// Test structs for struct diff testing
type Person struct {
	Name string
	Age  int
	City string
}

type Product struct {
	ID    int
	Name  string
	Price float64
	Tags  []string
}

// TestStructDiff tests struct diff functionality with behaviour-focused testing
func TestStructDiff(t *testing.T) {
	tests := []struct {
		name                string
		setupAndAssert      func(assert *Assert)
		expectErrorContains []string
		shouldPass          bool
	}{
		{
			name: "identical structs should pass",
			setupAndAssert: func(assert *Assert) {
				got := Person{Name: "Alice", Age: 30, City: "London"}
				want := Person{Name: "Alice", Age: 30, City: "London"}
				assert.StructDiff(got, want)
			},
			shouldPass: true,
		},
		{
			name: "different field values should fail with field details",
			setupAndAssert: func(assert *Assert) {
				got := Person{Name: "Alice", Age: 25, City: "London"}
				want := Person{Name: "Alice", Age: 30, City: "London"}
				assert.StructDiff(got, want)
			},
			expectErrorContains: []string{
				"structs differ",
				"field \"Age\"",
				"got: 25",
				"want: 30",
			},
			shouldPass: false,
		},
		{
			name: "multiple field differences should show first difference",
			setupAndAssert: func(assert *Assert) {
				got := Person{Name: "Bob", Age: 25, City: "Manchester"}
				want := Person{Name: "Alice", Age: 30, City: "London"}
				assert.StructDiff(got, want)
			},
			expectErrorContains: []string{
				"structs differ",
				"field \"Name\"",
				"got: Bob",
				"want: Alice",
			},
			shouldPass: false,
		},
		{
			name: "complex struct with slice field differences",
			setupAndAssert: func(assert *Assert) {
				got := Product{ID: 1, Name: "Widget", Price: 9.99, Tags: []string{"small", "red"}}
				want := Product{ID: 1, Name: "Widget", Price: 9.99, Tags: []string{"small", "blue"}}
				assert.StructDiff(got, want)
			},
			expectErrorContains: []string{
				"structs differ",
				"field \"Tags\"",
			},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a dummy testing.T that captures failures
			dummyT := &capturingT{}
			assert := New(dummyT)

			// Execute the struct diff assertion
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

// ExampleAssert_StructDiff demonstrates proper usage of struct diff assertion
func ExampleAssert_StructDiff() {
	assert := New(&testing.T{})

	// Test that two structs are identical
	type User struct {
		Name string
		Age  int
	}

	got := User{Name: "Alice", Age: 30}
	want := User{Name: "Alice", Age: 30}
	assert.StructDiff(got, want)

	fmt.Println("No error:", assert.Error() == "")
	// Output: No error: true
}
