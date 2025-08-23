// Package assertions provides assertion functions for testing.
//
// AssertionsTest contains unit tests for the assertions package.

package assertions

import (
	"fmt"
	"testing"
)

// TestAssertions contains the regular tests for the assertions package.
func TestAssertions(t *testing.T) {
	t.Run("Equal", func(t *testing.T) {
		assert := New(t)

		// Test cases
		testCases := []struct {
			expected, actual interface{}
			pass             bool
		}{
			{42, 42, true},
			{true, true, true},
			{false, false, true},
		}

		for i, tc := range testCases {
			t.Run(fmt.Sprintf("Test case %d", i+1), func(t *testing.T) {
				assert.Equal(tc.expected, tc.actual)

				if tc.pass && assert.Error() != "" {
					t.Errorf("Test case %d failed, expected no error but got: %s", i+1, assert.Error())
				} else if !tc.pass && assert.Error() == "" {
					t.Errorf("Test case %d failed, expected an error but got none", i+1)
				}
			})
		}
	})

	t.Run("NotEqual", func(t *testing.T) {
		assert := New(t)

		// Test cases
		testCases := []struct {
			expected, actual interface{}
			pass             bool
		}{
			{42, 42, false},
			{42, 23, false},
			{"hello", "world", false},
			{true, true, false},
			{false, true, false},
		}

		for i, tc := range testCases {
			t.Run(fmt.Sprintf("Test case %d", i+1), func(t *testing.T) {
				assert.NotEqual(tc.expected, tc.actual)

				if tc.pass && assert.Error() != "" {
					t.Errorf("Test case %d failed, expected no error but got: %s", i+1, assert.Error())
				} else if !tc.pass && assert.Error() == "" {
					t.Errorf("Test case %d failed, expected an error but got none", i+1)
				}
			})
		}
	})

	t.Run("True", func(t *testing.T) {
		assert := New(t)

		// Test cases
		testCases := []struct {
			value bool
			pass  bool
		}{
			{true, true},
			{false, false},
		}

		for i, tc := range testCases {
			t.Run(fmt.Sprintf("Test case %d", i+1), func(t *testing.T) {
				assert.True(tc.value)

				if tc.pass && assert.Error() != "" {
					t.Errorf("Test case %d failed, expected no error but got: %s", i+1, assert.Error())
				} else if !tc.pass && assert.Error() == "" {
					t.Errorf("Test case %d failed, expected an error but got none", i+1)
				}
			})
		}
	})

	t.Run("False", func(t *testing.T) {
		assert := New(t)

		// Test cases
		testCases := []struct {
			value bool
			pass  bool
		}{
			{false, true},
		}

		for i, tc := range testCases {
			t.Run(fmt.Sprintf("Test case %d", i+1), func(t *testing.T) {
				assert.False(tc.value)

				if tc.pass && assert.Error() != "" {
					t.Errorf("Test case %d failed, expected no error but got: %s", i+1, assert.Error())
				} else if !tc.pass && assert.Error() == "" {
					t.Errorf("Test case %d failed, expected an error but got none", i+1)
				}
			})
		}
	})
}

// ExampleAssert_Equal demonstrates the Equal assertion for fast-path equality checking.
func ExampleAssert_Equal() {
	assert := New(&testing.T{})

	// Comparable types use fast-path
	assert.Equal(42, 42)
	assert.Equal("hello", "hello")
	assert.Equal(true, true)

	// Complex types use deep equality
	assert.Equal([]int{1, 2, 3}, []int{1, 2, 3})
	assert.Equal(map[string]int{"a": 1}, map[string]int{"a": 1})

	fmt.Println("No error:", assert.Error() == "")
	// Output: No error: true
}

// ExampleAssert_NotEqual demonstrates the NotEqual assertion.
func ExampleAssert_NotEqual() {
	assert := New(&testing.T{})

	assert.NotEqual(42, 24)
	assert.NotEqual("hello", "world")
	assert.NotEqual([]int{1, 2}, []int{1, 2, 3})

	fmt.Println("No error:", assert.Error() == "")
	// Output: No error: true
}

// ExampleAssert_DeepEqual demonstrates explicit deep equality checking.
func ExampleAssert_DeepEqual() {
	assert := New(&testing.T{})

	type Person struct {
		Name string
		Age  int
	}

	assert.DeepEqual(
		Person{Name: "Alice", Age: 30},
		Person{Name: "Alice", Age: 30},
	)

	fmt.Println("No error:", assert.Error() == "")
	// Output: No error: true
}

// ExampleAssert_Same demonstrates pointer identity comparison.
func ExampleAssert_Same() {
	assert := New(&testing.T{})

	x := 42
	ptr1 := &x
	ptr2 := ptr1

	assert.Same(ptr1, ptr2) // same pointer

	fmt.Println("No error:", assert.Error() == "")
	// Output: No error: true
}

// ExampleAssert_True demonstrates boolean true assertion.
func ExampleAssert_True() {
	assert := New(&testing.T{})

	assert.True(2 > 1)
	assert.True(len("hello") == 5)

	fmt.Println("No error:", assert.Error() == "")
	// Output: No error: true
}

// ExampleAssert_False demonstrates boolean false assertion.
func ExampleAssert_False() {
	assert := New(&testing.T{})

	assert.False(2 < 1)
	assert.False(len("") > 0)

	fmt.Println("No error:", assert.Error() == "")
	// Output: No error: true
}
