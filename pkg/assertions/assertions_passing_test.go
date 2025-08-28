// Package assertions provides assertion functions for testing.
//
// AssertionsTest contains unit tests for the assertions package.

package assertions

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestWithinTimeout tests the WithinTimeout assertion.
func TestWithinTimeout(t *testing.T) {
	t.Run("FunctionCompletesWithinTimeout", func(t *testing.T) {
		assert := New(&mockT{})

		// Function that completes quickly
		fastFunc := func() {
			// Completes immediately
		}

		assert.WithinTimeout(fastFunc, 100*time.Millisecond)

		if assert.Error() != "" {
			t.Errorf("Expected no error for fast function, got: %s", assert.Error())
		}
	})

	t.Run("FunctionExceedsTimeout", func(t *testing.T) {
		assert := New(&mockT{})

		// Function that takes too long
		slowFunc := func() {
			time.Sleep(200 * time.Millisecond)
		}

		assert.WithinTimeout(slowFunc, 100*time.Millisecond)

		if assert.Error() == "" {
			t.Error("Expected timeout error for slow function")
		}
		if !strings.Contains(assert.Error(), "timeout") {
			t.Errorf("Expected error message to contain 'timeout', got: %s", assert.Error())
		}
		if !strings.Contains(assert.Error(), "elapsed:") {
			t.Errorf("Expected error message to contain timing info, got: %s", assert.Error())
		}
	})

	t.Run("FunctionPanics", func(t *testing.T) {
		assert := New(&mockT{})

		// Function that panics
		panicFunc := func() {
			panic("test panic")
		}

		assert.WithinTimeout(panicFunc, 100*time.Millisecond)

		// BEHAVIOURAL QUESTION: Should panic = "completed within timeout"?
		// Current implementation: YES (panic is caught, signals completion)
		// Alternative behaviour: NO (panic should be treated as failure)

		// Testing current behaviour - panicked function "completes"
		if assert.Error() != "" {
			t.Errorf("Expected no error for panicked function (current behaviour), got: %s", assert.Error())
		}
	})

	t.Run("NegativeTimeout", func(t *testing.T) {
		assert := New(&mockT{})

		fastFunc := func() {}

		assert.WithinTimeout(fastFunc, -1*time.Second)

		// Should apply default timeout and pass for fast function
		if assert.Error() != "" {
			t.Errorf("Expected no error with default timeout applied, got: %s", assert.Error())
		}
	})

	t.Run("ZeroTimeout", func(t *testing.T) {
		assert := New(&mockT{})

		fastFunc := func() {}

		assert.WithinTimeout(fastFunc, 0)

		// Should apply default timeout and pass for fast function
		if assert.Error() != "" {
			t.Errorf("Expected no error with default timeout applied, got: %s", assert.Error())
		}
	})

	t.Run("ActualTimingAccuracy", func(t *testing.T) {
		assert := New(&mockT{})

		// Function that does actual work for a measurable time
		workFunc := func() {
			start := time.Now()
			for time.Since(start) < 50*time.Millisecond {
				// Busy work to consume actual time
				_ = make([]int, 1000)
			}
		}

		// Should complete within reasonable timeout
		assert.WithinTimeout(workFunc, 200*time.Millisecond)

		if assert.Error() != "" {
			t.Errorf("Expected work function to complete within timeout, got: %s", assert.Error())
		}
	})

	t.Run("TimeoutBehaviorWithRealWork", func(t *testing.T) {
		assert := New(&mockT{})

		// Function that definitely takes longer than timeout
		longWorkFunc := func() {
			start := time.Now()
			for time.Since(start) < 300*time.Millisecond {
				// Busy work - can't be optimised away
				_ = make([]int, 1000)
			}
		}

		// Should timeout with very short timeout
		startTime := time.Now()
		assert.WithinTimeout(longWorkFunc, 50*time.Millisecond)
		elapsed := time.Since(startTime)

		if assert.Error() == "" {
			t.Error("Expected timeout error for long-running function")
		}

		// Verify timing accuracy - should be close to timeout duration
		if elapsed < 40*time.Millisecond || elapsed > 100*time.Millisecond {
			t.Errorf("Expected timeout around 50ms, but took %v", elapsed)
		}
	})
}

// ExampleAssert_WithinTimeout demonstrates timeout assertion usage.
func ExampleAssert_WithinTimeout() {
	assert := New(nil) // mockT for example

	// Fast operation should complete within timeout
	fastOperation := func() {
		time.Sleep(10 * time.Millisecond)
	}

	assert.WithinTimeout(fastOperation, 100*time.Millisecond)
	// No error - completes quickly

	// Slow operation exceeds timeout
	slowOperation := func() {
		time.Sleep(200 * time.Millisecond)
	}

	assert.WithinTimeout(slowOperation, 100*time.Millisecond)
	// Error: "WithinTimeout: function did not complete within timeout\n  timeout: 100ms\n  elapsed: ~200ms"
}

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
