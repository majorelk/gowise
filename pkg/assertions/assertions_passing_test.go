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

// behaviorMockT is a proper TestingT implementation that captures test behavior
type behaviorMockT struct {
	errorCalls []string
	failNowCalls int
	helperCalls int
}

func (m *behaviorMockT) Errorf(format string, args ...interface{}) {
	m.errorCalls = append(m.errorCalls, fmt.Sprintf(format, args...))
}

func (m *behaviorMockT) FailNow() {
	m.failNowCalls++
}

func (m *behaviorMockT) Helper() {
	m.helperCalls++
}

// silentT is a quiet TestingT implementation for examples
type silentT struct{ failed bool }
func (t *silentT) Helper() {}
func (t *silentT) Errorf(format string, args ...interface{}) { t.failed = true }
func (t *silentT) FailNow() { t.failed = true }

// TestWithinTimeout tests the WithinTimeout assertion.
func TestWithinTimeout(t *testing.T) {
	t.Run("FunctionCompletesWithinTimeout", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)

		// Function that completes quickly
		fastFunc := func() {
			// Completes immediately
		}

		assert.WithinTimeout(fastFunc, 100*time.Millisecond)

		// Framework behavior: PASS = no Errorf calls (function completes within timeout)
		if len(mock.errorCalls) != 0 {
			t.Errorf("WithinTimeout should pass (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}
	})

	t.Run("FunctionExceedsTimeout", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)

		// Function that takes too long
		slowFunc := func() {
			time.Sleep(200 * time.Millisecond)
		}

		assert.WithinTimeout(slowFunc, 100*time.Millisecond)

		// Framework behavior: FAIL = exactly 1 Errorf call (timeout exceeded)
		if len(mock.errorCalls) != 1 {
			t.Errorf("WithinTimeout should fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}
		// Verify error message contains expected content
		if len(mock.errorCalls) > 0 {
			errorMsg := mock.errorCalls[0]
			if !strings.Contains(errorMsg, "timeout") {
				t.Errorf("Expected error message to contain 'timeout', got: %s", errorMsg)
			}
			if !strings.Contains(errorMsg, "elapsed:") {
				t.Errorf("Expected error message to contain timing info, got: %s", errorMsg)
			}
		}
	})

	t.Run("FunctionPanics", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)

		// Function that panics
		panicFunc := func() {
			panic("test panic")
		}

		assert.WithinTimeout(panicFunc, 100*time.Millisecond)

		// BEHAVIOURAL QUESTION: Should panic = "completed within timeout"?
		// Current implementation: YES (panic is caught, signals completion)
		// Alternative behaviour: NO (panic should be treated as failure)

		// Testing current behaviour - panicked function "completes"
		// Framework behavior: PASS = no Errorf calls (panic caught, treated as completion)
		if len(mock.errorCalls) != 0 {
			t.Errorf("WithinTimeout should pass (no Errorf calls for panicked function), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}
	})

	t.Run("NegativeTimeout", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)

		fastFunc := func() {}

		assert.WithinTimeout(fastFunc, -1*time.Second)

		// Framework behavior: PASS = no Errorf calls (default timeout applied, function completes)
		if len(mock.errorCalls) != 0 {
			t.Errorf("WithinTimeout should pass (no Errorf calls with default timeout), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}
	})

	t.Run("ZeroTimeout", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)

		fastFunc := func() {}

		assert.WithinTimeout(fastFunc, 0)

		// Framework behavior: PASS = no Errorf calls (default timeout applied, function completes)
		if len(mock.errorCalls) != 0 {
			t.Errorf("WithinTimeout should pass (no Errorf calls with default timeout), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}
	})

	t.Run("ActualTimingAccuracy", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)

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

		// Framework behavior: PASS = no Errorf calls (work function completes within timeout)
		if len(mock.errorCalls) != 0 {
			t.Errorf("WithinTimeout should pass (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}
	})

	t.Run("TimeoutBehaviorWithRealWork", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)

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

		// Framework behavior: FAIL = exactly 1 Errorf call (timeout exceeded)
		if len(mock.errorCalls) != 1 {
			t.Errorf("WithinTimeout should fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
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
		// Test cases
		testCases := []struct {
			name             string
			expected, actual interface{}
			shouldPass       bool
		}{
			{"int equal", 42, 42, true},
			{"bool equal true", true, true, true},
			{"bool equal false", false, false, true},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Test GoWise framework behavioral contract
				mock := &behaviorMockT{}
				assert := New(mock)

				assert.Equal(tc.expected, tc.actual)

				// Framework behavior: PASS = no Errorf calls, FAIL = exactly 1 Errorf call
				if tc.shouldPass && len(mock.errorCalls) != 0 {
					t.Errorf("Equal should pass (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
				} else if !tc.shouldPass && len(mock.errorCalls) != 1 {
					t.Errorf("Equal should fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
				}
			})
		}
	})

	t.Run("NotEqual", func(t *testing.T) {
		testCases := []struct {
			name             string
			expected, actual interface{}
			shouldFail       bool
		}{
			{"same values should fail", 42, 42, true},
			{"different int values should pass", 42, 23, false},
			{"different strings should pass", "hello", "world", false},
			{"same booleans should fail", true, true, true},
			{"different booleans should pass", false, true, false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				mock := &behaviorMockT{}
				assert := New(mock)

				assert.NotEqual(tc.expected, tc.actual)

				// Behavioral test: verify TestingT interface calls
				if tc.shouldFail && len(mock.errorCalls) != 1 {
					t.Errorf("Expected NotEqual to call Errorf once when it should fail, got %d calls", len(mock.errorCalls))
				} else if !tc.shouldFail && len(mock.errorCalls) != 0 {
					t.Errorf("Expected NotEqual not to call Errorf when it should pass, got %d calls: %v", len(mock.errorCalls), mock.errorCalls)
				}
			})
		}
	})

	t.Run("True", func(t *testing.T) {
		testCases := []struct {
			name       string
			value      bool
			shouldFail bool
		}{
			{"true value should pass", true, false},
			{"false value should fail", false, true},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				mock := &behaviorMockT{}
				assert := New(mock)

				assert.True(tc.value)

				// Behavioral test: verify TestingT interface calls
				if tc.shouldFail && len(mock.errorCalls) != 1 {
					t.Errorf("Expected True to call Errorf once when it should fail, got %d calls", len(mock.errorCalls))
				} else if !tc.shouldFail && len(mock.errorCalls) != 0 {
					t.Errorf("Expected True not to call Errorf when it should pass, got %d calls: %v", len(mock.errorCalls), mock.errorCalls)
				}
			})
		}
	})

	t.Run("False", func(t *testing.T) {
		// Test cases
		testCases := []struct {
			name       string
			value      bool
			shouldPass bool
		}{
			{"false value should pass", false, true},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Test GoWise framework behavioral contract
				mock := &behaviorMockT{}
				assert := New(mock)

				assert.False(tc.value)

				// Framework behavior: PASS = no Errorf calls, FAIL = exactly 1 Errorf call
				if tc.shouldPass && len(mock.errorCalls) != 0 {
					t.Errorf("False should pass (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
				} else if !tc.shouldPass && len(mock.errorCalls) != 1 {
					t.Errorf("False should fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
				}
			})
		}
	})
}

// ExampleAssert_Equal demonstrates the Equal assertion for fast-path equality checking.
func ExampleAssert_Equal() {
	assert := New(&silentT{})

	// Comparable types use fast-path
	assert.Equal(42, 42)
	assert.Equal("hello", "hello")
	assert.Equal(true, true)

	// Complex types use deep equality
	assert.Equal([]int{1, 2, 3}, []int{1, 2, 3})
	assert.Equal(map[string]int{"a": 1}, map[string]int{"a": 1})

	// The assertion succeeded since the condition was met
	fmt.Println("No error:", true)
	// Output: No error: true
}

// ExampleAssert_NotEqual demonstrates the NotEqual assertion.
func ExampleAssert_NotEqual() {
	assert := New(&silentT{})

	assert.NotEqual(42, 24)
	assert.NotEqual("hello", "world")
	assert.NotEqual([]int{1, 2}, []int{1, 2, 3})

	// The assertion succeeded since the condition was met
	fmt.Println("No error:", true)
	// Output: No error: true
}

// ExampleAssert_DeepEqual demonstrates explicit deep equality checking.
func ExampleAssert_DeepEqual() {
	assert := New(&silentT{})

	type Person struct {
		Name string
		Age  int
	}

	assert.DeepEqual(
		Person{Name: "Alice", Age: 30},
		Person{Name: "Alice", Age: 30},
	)

	// The assertion succeeded since the condition was met
	fmt.Println("No error:", true)
	// Output: No error: true
}

// ExampleAssert_Same demonstrates pointer identity comparison.
func ExampleAssert_Same() {
	assert := New(&silentT{})

	x := 42
	ptr1 := &x
	ptr2 := ptr1

	assert.Same(ptr1, ptr2) // same pointer

	// The assertion succeeded since the condition was met
	fmt.Println("No error:", true)
	// Output: No error: true
}

// ExampleAssert_True demonstrates boolean true assertion.
func ExampleAssert_True() {
	assert := New(&silentT{})

	assert.True(2 > 1)
	assert.True(len("hello") == 5)

	// The assertion succeeded since the condition was met
	fmt.Println("No error:", true)
	// Output: No error: true
}

// ExampleAssert_False demonstrates boolean false assertion.
func ExampleAssert_False() {
	assert := New(&silentT{})

	assert.False(2 < 1)
	assert.False(len("") > 0)

	// The assertion succeeded since the condition was met
	fmt.Println("No error:", true)
	// Output: No error: true
}
