package assertions

import (
	"fmt"
	"testing"
)

// TestContainsAssertion tests the Contains assertion with various container types.
func TestContainsAssertion(t *testing.T) {
	t.Run("string contains substring", func(t *testing.T) {
		tests := []struct {
			name       string
			container  string
			item       string
			shouldPass bool
		}{
			{"contains substring", "hello world", "world", true},
			{"does not contain", "hello world", "foo", false},
			{"empty string contains empty", "", "", true},
			{"empty string does not contain", "", "foo", false},
			{"contains at start", "hello world", "hello", true},
			{"contains at end", "hello world", "world", true},
			{"case sensitive", "Hello World", "hello", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert := New(t)

				assert.Contains(tt.container, tt.item)

				hasError := assert.Error() != ""
				if tt.shouldPass && hasError {
					t.Errorf("Expected pass but got error: %s", assert.Error())
				} else if !tt.shouldPass && !hasError {
					t.Errorf("Expected failure but assertion passed")
				}
			})
		}
	})

	t.Run("slice contains element", func(t *testing.T) {
		tests := []struct {
			name       string
			container  interface{}
			item       interface{}
			shouldPass bool
		}{
			{"int slice contains", []int{1, 2, 3, 4, 5}, 3, true},
			{"int slice does not contain", []int{1, 2, 3, 4, 5}, 6, false},
			{"string slice contains", []string{"a", "b", "c"}, "b", true},
			{"string slice does not contain", []string{"a", "b", "c"}, "d", false},
			{"empty slice", []int{}, 1, false},
			{"struct slice contains", []struct{ X int }{{1}, {2}, {3}}, struct{ X int }{2}, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert := New(t)

				assert.Contains(tt.container, tt.item)

				hasError := assert.Error() != ""
				if tt.shouldPass && hasError {
					t.Errorf("Expected pass but got error: %s", assert.Error())
				} else if !tt.shouldPass && !hasError {
					t.Errorf("Expected failure but assertion passed")
				}
			})
		}
	})

	t.Run("array contains element", func(t *testing.T) {
		tests := []struct {
			name       string
			container  interface{}
			item       interface{}
			shouldPass bool
		}{
			{"int array contains", [5]int{1, 2, 3, 4, 5}, 3, true},
			{"int array does not contain", [5]int{1, 2, 3, 4, 5}, 6, false},
			{"string array contains", [3]string{"a", "b", "c"}, "b", true},
			{"string array does not contain", [3]string{"a", "b", "c"}, "d", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert := New(t)

				assert.Contains(tt.container, tt.item)

				hasError := assert.Error() != ""
				if tt.shouldPass && hasError {
					t.Errorf("Expected pass but got error: %s", assert.Error())
				} else if !tt.shouldPass && !hasError {
					t.Errorf("Expected failure but assertion passed")
				}
			})
		}
	})

	t.Run("map contains key", func(t *testing.T) {
		tests := []struct {
			name       string
			container  interface{}
			item       interface{}
			shouldPass bool
		}{
			{"string map contains key", map[string]int{"a": 1, "b": 2, "c": 3}, "b", true},
			{"string map does not contain key", map[string]int{"a": 1, "b": 2, "c": 3}, "d", false},
			{"int map contains key", map[int]string{1: "one", 2: "two", 3: "three"}, 2, true},
			{"int map does not contain key", map[int]string{1: "one", 2: "two", 3: "three"}, 4, false},
			{"empty map", map[string]int{}, "a", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert := New(t)

				assert.Contains(tt.container, tt.item)

				hasError := assert.Error() != ""
				if tt.shouldPass && hasError {
					t.Errorf("Expected pass but got error: %s", assert.Error())
				} else if !tt.shouldPass && !hasError {
					t.Errorf("Expected failure but assertion passed")
				}
			})
		}
	})

	t.Run("edge cases", func(t *testing.T) {
		assert := New(t)

		// Nil container should fail
		assert.Contains(nil, "test")
		if assert.Error() == "" {
			t.Errorf("Expected error for nil container")
		}

		// Wrong item type for string
		assert = New(t) // reset
		assert.Contains("hello", 123)
		if assert.Error() == "" {
			t.Errorf("Expected error for wrong item type in string")
		}

		// Wrong key type for map
		assert = New(t) // reset
		assert.Contains(map[string]int{"a": 1}, 123)
		if assert.Error() == "" {
			t.Errorf("Expected error for wrong key type in map")
		}

		// Unsupported container type
		assert = New(t) // reset
		assert.Contains(123, 1)
		if assert.Error() == "" {
			t.Errorf("Expected error for unsupported container type")
		}
	})
}

// TestLenAssertion tests the Len assertion with various container types.
func TestLenAssertion(t *testing.T) {
	tests := []struct {
		name        string
		container   interface{}
		expectedLen int
		shouldPass  bool
	}{
		// String
		{"empty string", "", 0, true},
		{"non-empty string", "hello", 5, true},
		{"wrong string length", "hello", 3, false},

		// Slice
		{"empty slice", []int{}, 0, true},
		{"slice with elements", []int{1, 2, 3}, 3, true},
		{"wrong slice length", []int{1, 2, 3}, 2, false},

		// Array
		{"array with elements", [3]int{1, 2, 3}, 3, true},
		{"wrong array length", [3]int{1, 2, 3}, 2, false},

		// Map
		{"empty map", map[string]int{}, 0, true},
		{"map with elements", map[string]int{"a": 1, "b": 2}, 2, true},
		{"wrong map length", map[string]int{"a": 1, "b": 2}, 3, false},

		// Channel
		{"buffered channel", make(chan int, 5), 0, true}, // Empty buffered channel has length 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := New(t)

			assert.Len(tt.container, tt.expectedLen)

			hasError := assert.Error() != ""
			if tt.shouldPass && hasError {
				t.Errorf("Expected pass but got error: %s", assert.Error())
			} else if !tt.shouldPass && !hasError {
				t.Errorf("Expected failure but assertion passed")
			}
		})
	}

	t.Run("edge cases", func(t *testing.T) {
		assert := New(t)

		// Nil container should fail
		assert.Len(nil, 0)
		if assert.Error() == "" {
			t.Errorf("Expected error for nil container")
		}

		// Unsupported type should fail
		assert = New(t) // reset
		assert.Len(123, 3)
		if assert.Error() == "" {
			t.Errorf("Expected error for unsupported type")
		}
	})

	t.Run("channel with elements", func(t *testing.T) {
		// Test channel with elements
		ch := make(chan int, 3)
		ch <- 1
		ch <- 2

		assert := New(t)
		assert.Len(ch, 2) // Should have 2 elements
		if assert.Error() != "" {
			t.Errorf("Expected channel to have length 2: %s", assert.Error())
		}

		// Close and drain for cleanup
		close(ch)
		for range ch {
			// drain
		}
	})
}

// Examples for documentation.

func ExampleAssert_Contains() {
	assert := New(&testing.T{})

	// String contains substring
	assert.Contains("hello world", "world")

	// Slice contains element
	assert.Contains([]int{1, 2, 3, 4, 5}, 3)

	// Array contains element
	assert.Contains([3]string{"a", "b", "c"}, "b")

	// Map contains key
	assert.Contains(map[string]int{"foo": 1, "bar": 2}, "foo")

	fmt.Println("No error:", assert.Error() == "")
	// Output: No error: true
}

func ExampleAssert_Len() {
	assert := New(&testing.T{})

	// String length
	assert.Len("hello", 5)

	// Slice length
	assert.Len([]int{1, 2, 3}, 3)

	// Array length
	assert.Len([3]string{"a", "b", "c"}, 3)

	// Map length
	assert.Len(map[string]int{"a": 1, "b": 2}, 2)

	// Channel length (buffered)
	ch := make(chan int, 5)
	assert.Len(ch, 0) // Empty channel

	fmt.Println("No error:", assert.Error() == "")
	// Output: No error: true
}
