package assertions

import (
	"fmt"
	"testing"
)

// behaviorMockT is defined in assertions_passing_test.go - shared across test files

// TestNilAssertions tests the Nil assertion across different types.
func TestNilAssertions(t *testing.T) {
	tests := []struct {
		name       string
		value      interface{}
		shouldPass bool
	}{
		// Untyped nil
		{"untyped nil", nil, true},

		// Pointer types
		{"nil pointer", (*int)(nil), true},
		{"non-nil pointer", func() *int { x := 42; return &x }(), false},

		// Interface types
		{"nil interface", (error)(nil), true},
		{"non-nil interface", func() error { return &TestError{} }(), false},

		// Slice types
		{"nil slice", ([]int)(nil), true},
		{"empty slice", []int{}, false},
		{"non-empty slice", []int{1, 2, 3}, false},

		// Map types
		{"nil map", (map[string]int)(nil), true},
		{"empty map", map[string]int{}, false},
		{"non-empty map", map[string]int{"a": 1}, false},

		// Channel types
		{"nil channel", (chan int)(nil), true},
		{"non-nil channel", make(chan int), false},

		// Function types
		{"nil function", (func())(nil), true},
		{"non-nil function", func() {}, false},

		// Non-nillable types
		{"int value", 42, false},
		{"string value", "hello", false},
		{"bool value", true, false},
		{"struct value", struct{ X int }{42}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test GoWise framework behavioral contract
			mock := &behaviorMockT{}
			assert := New(mock)

			assert.Nil(tt.value)

			// Framework behavior: PASS = no Errorf calls, FAIL = exactly 1 Errorf call
			if tt.shouldPass && len(mock.errorCalls) != 0 {
				t.Errorf("Nil should pass (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
			} else if !tt.shouldPass && len(mock.errorCalls) != 1 {
				t.Errorf("Nil should fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
			}
		})
	}
}

// TestNotNilAssertions tests the NotNil assertion across different types.
func TestNotNilAssertions(t *testing.T) {
	tests := []struct {
		name       string
		value      interface{}
		shouldPass bool
	}{
		// Untyped nil
		{"untyped nil", nil, false},

		// Pointer types
		{"nil pointer", (*int)(nil), false},
		{"non-nil pointer", func() *int { x := 42; return &x }(), true},

		// Interface types
		{"nil interface", (error)(nil), false},
		{"non-nil interface", func() error { return &TestError{} }(), true},

		// Slice types
		{"nil slice", ([]int)(nil), false},
		{"empty slice", []int{}, true},
		{"non-empty slice", []int{1, 2, 3}, true},

		// Map types
		{"nil map", (map[string]int)(nil), false},
		{"empty map", map[string]int{}, true},
		{"non-empty map", map[string]int{"a": 1}, true},

		// Channel types
		{"nil channel", (chan int)(nil), false},
		{"non-nil channel", make(chan int), true},

		// Function types
		{"nil function", (func())(nil), false},
		{"non-nil function", func() {}, true},

		// Non-nillable types (should always pass NotNil)
		{"int value", 42, true},
		{"string value", "hello", true},
		{"bool value", true, true},
		{"struct value", struct{ X int }{42}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test GoWise framework behavioral contract
			mock := &behaviorMockT{}
			assert := New(mock)

			assert.NotNil(tt.value)

			// Framework behavior: PASS = no Errorf calls, FAIL = exactly 1 Errorf call
			if tt.shouldPass && len(mock.errorCalls) != 0 {
				t.Errorf("NotNil should pass (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
			} else if !tt.shouldPass && len(mock.errorCalls) != 1 {
				t.Errorf("NotNil should fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
			}
		})
	}
}

// TestTypedVsUntypedNil tests the distinction between typed and untyped nil.
func TestTypedVsUntypedNil(t *testing.T) {
	t.Run("typed nil vs untyped nil", func(t *testing.T) {
		var untypedNil interface{} = nil
		var typedNilPointer *int = nil
		var typedNilSlice []int = nil
		var typedNilMap map[string]int = nil
		var typedNilChan chan int = nil
		var typedNilFunc func() = nil

		tests := []struct {
			name  string
			value interface{}
		}{
			{"untyped nil", untypedNil},
			{"typed nil pointer", typedNilPointer},
			{"typed nil slice", typedNilSlice},
			{"typed nil map", typedNilMap},
			{"typed nil channel", typedNilChan},
			{"typed nil function", typedNilFunc},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Test GoWise framework behavioral contract - all should be detected as nil
				mock := &behaviorMockT{}
				assert := New(mock)

				assert.Nil(tt.value)

				// Framework behavior: PASS = no Errorf calls (all should be nil)
				if len(mock.errorCalls) != 0 {
					t.Errorf("%s should be nil (no Errorf calls), got %d: %v", tt.name, len(mock.errorCalls), mock.errorCalls)
				}
			})
		}
	})

	t.Run("nil interface with concrete nil value", func(t *testing.T) {
		// Test the subtle case where an interface contains a typed nil
		var err error = (*TestError)(nil) // interface containing typed nil

		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)

		assert.Nil(err)

		// Framework behavior: PASS = no Errorf calls (should be detected as nil)
		if len(mock.errorCalls) != 0 {
			t.Errorf("Interface containing typed nil should be nil (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}
	})
}

// TestError is a helper type for testing.
type TestError struct {
	message string
}

func (e *TestError) Error() string {
	return e.message
}

// Test examples for nil assertions.
func ExampleAssert_Nil() {
	assert := New(&testing.T{})

	// Test various nil types
	assert.Nil((*int)(nil))           // nil pointer
	assert.Nil(([]int)(nil))          // nil slice
	assert.Nil((map[string]int)(nil)) // nil map
	assert.Nil((chan int)(nil))       // nil channel
	assert.Nil((func())(nil))         // nil function
	assert.Nil((error)(nil))          // nil interface

	fmt.Println("No error:", assert.Error() == "")
	// Output: No error: true
}

func ExampleAssert_NotNil() {
	assert := New(&testing.T{})

	x := 42
	assert.NotNil(&x)               // non-nil pointer
	assert.NotNil([]int{})          // empty but not nil slice
	assert.NotNil(map[string]int{}) // empty but not nil map
	assert.NotNil(make(chan int))   // non-nil channel
	assert.NotNil(func() {})        // non-nil function

	fmt.Println("No error:", assert.Error() == "")
	// Output: No error: true
}
