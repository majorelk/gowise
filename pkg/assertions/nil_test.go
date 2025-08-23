package assertions

import (
	"fmt"
	"testing"
)

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
			assert := New(t)

			assert.Nil(tt.value)

			hasError := assert.Error() != ""
			if tt.shouldPass && hasError {
				t.Errorf("Expected pass but got error: %s", assert.Error())
			} else if !tt.shouldPass && !hasError {
				t.Errorf("Expected failure but assertion passed")
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
			assert := New(t)

			assert.NotNil(tt.value)

			hasError := assert.Error() != ""
			if tt.shouldPass && hasError {
				t.Errorf("Expected pass but got error: %s", assert.Error())
			} else if !tt.shouldPass && !hasError {
				t.Errorf("Expected failure but assertion passed")
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

		assert := New(t)

		// All should be detected as nil
		assert.Nil(untypedNil)
		if assert.Error() != "" {
			t.Errorf("Untyped nil should be nil: %s", assert.Error())
		}

		assert = New(t) // reset
		assert.Nil(typedNilPointer)
		if assert.Error() != "" {
			t.Errorf("Typed nil pointer should be nil: %s", assert.Error())
		}

		assert = New(t) // reset
		assert.Nil(typedNilSlice)
		if assert.Error() != "" {
			t.Errorf("Typed nil slice should be nil: %s", assert.Error())
		}

		assert = New(t) // reset
		assert.Nil(typedNilMap)
		if assert.Error() != "" {
			t.Errorf("Typed nil map should be nil: %s", assert.Error())
		}

		assert = New(t) // reset
		assert.Nil(typedNilChan)
		if assert.Error() != "" {
			t.Errorf("Typed nil channel should be nil: %s", assert.Error())
		}

		assert = New(t) // reset
		assert.Nil(typedNilFunc)
		if assert.Error() != "" {
			t.Errorf("Typed nil function should be nil: %s", assert.Error())
		}
	})

	t.Run("nil interface with concrete nil value", func(t *testing.T) {
		// Test the subtle case where an interface contains a typed nil
		var err error = (*TestError)(nil) // interface containing typed nil

		assert := New(t)
		assert.Nil(err)
		if assert.Error() != "" {
			t.Errorf("Interface containing typed nil should be nil: %s", assert.Error())
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
