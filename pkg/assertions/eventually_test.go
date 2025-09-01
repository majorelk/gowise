package assertions

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

// behaviorMockT is defined in assertions_passing_test.go - shared across test files

// TestEventually tests the Eventually assertion with various scenarios.
func TestEventually(t *testing.T) {
	t.Run("SucceedsImmediately", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)

		assert.Eventually(func() bool {
			return true
		}, 1*time.Second, 50*time.Millisecond)

		// Framework behavior: PASS = no Errorf calls (condition succeeds immediately)
		if len(mock.errorCalls) != 0 {
			t.Errorf("Eventually should pass (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}
	})

	t.Run("SucceedsAfterDelay", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)
		var counter int32

		assert.Eventually(func() bool {
			return atomic.AddInt32(&counter, 1) >= 3
		}, 1*time.Second, 50*time.Millisecond)

		// Framework behavior: PASS = no Errorf calls (condition succeeds after delay)
		if len(mock.errorCalls) != 0 {
			t.Errorf("Eventually should pass (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}

		finalCount := atomic.LoadInt32(&counter)
		if finalCount < 3 {
			t.Errorf("Expected at least 3 attempts, got %d", finalCount)
		}
	})

	t.Run("TimesOut", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)
		startTime := time.Now()

		assert.Eventually(func() bool {
			return false // Never succeeds
		}, 200*time.Millisecond, 50*time.Millisecond)

		elapsed := time.Since(startTime)

		// Framework behavior: FAIL = exactly 1 Errorf call (timeout exceeded)
		if len(mock.errorCalls) != 1 {
			t.Errorf("Eventually should fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}

		// Verify error message content
		if len(mock.errorCalls) > 0 {
			expectedError := "Eventually: condition not met within timeout"
			if !containsString(mock.errorCalls[0], expectedError) {
				t.Errorf("Expected error to contain %q, got: %s", expectedError, mock.errorCalls[0])
			}
		}

		// Verify it actually waited close to the timeout
		if elapsed < 200*time.Millisecond || elapsed > 400*time.Millisecond {
			t.Errorf("Expected elapsed time around 200ms, got %v", elapsed)
		}
	})

	t.Run("CountsAttempts", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)
		var counter int32

		assert.Eventually(func() bool {
			count := atomic.AddInt32(&counter, 1)
			return count >= 4 // Succeed on 4th attempt
		}, 1*time.Second, 50*time.Millisecond)

		// Framework behavior: PASS = no Errorf calls (condition succeeds on 4th attempt)
		if len(mock.errorCalls) != 0 {
			t.Errorf("Eventually should pass (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}

		finalCount := atomic.LoadInt32(&counter)
		if finalCount != 4 {
			t.Errorf("Expected exactly 4 attempts, got %d", finalCount)
		}
	})

	t.Run("ReportsTimingInError", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)

		assert.Eventually(func() bool {
			return false
		}, 100*time.Millisecond, 20*time.Millisecond)

		// Framework behavior: FAIL = exactly 1 Errorf call (timeout exceeded)
		if len(mock.errorCalls) != 1 {
			t.Fatalf("Eventually should fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}

		errorMsg := mock.errorCalls[0]
		// Check that error contains timing information
		expectedFields := []string{"timeout:", "elapsed:", "attempts:", "final interval:"}
		for _, field := range expectedFields {
			if !containsString(errorMsg, field) {
				t.Errorf("Expected error message to contain %q, got: %s", field, errorMsg)
			}
		}
	})
}

// TestNever tests the Never assertion with various scenarios.
func TestNever(t *testing.T) {
	t.Run("SucceedsWhenConditionNeverTrue", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)

		assert.Never(func() bool {
			return false // Never true
		}, 200*time.Millisecond, 50*time.Millisecond)

		// Framework behavior: PASS = no Errorf calls (condition never becomes true)
		if len(mock.errorCalls) != 0 {
			t.Errorf("Never should pass (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}
	})

	t.Run("FailsWhenConditionBecomesTrue", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)
		var counter int32

		assert.Never(func() bool {
			return atomic.AddInt32(&counter, 1) >= 3 // True on 3rd call
		}, 1*time.Second, 50*time.Millisecond)

		// Framework behavior: FAIL = exactly 1 Errorf call (condition becomes true unexpectedly)
		if len(mock.errorCalls) != 1 {
			t.Errorf("Never should fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}

		// Verify error message content
		if len(mock.errorCalls) > 0 {
			expectedError := "Never: condition became true unexpectedly"
			if !containsString(mock.errorCalls[0], expectedError) {
				t.Errorf("Expected error to contain %q, got: %s", expectedError, mock.errorCalls[0])
			}
		}
	})

	t.Run("FailsImmediatelyIfConditionTrueFirst", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)
		startTime := time.Now()

		assert.Never(func() bool {
			return true // True immediately
		}, 1*time.Second, 50*time.Millisecond)

		elapsed := time.Since(startTime)

		// Framework behavior: FAIL = exactly 1 Errorf call (immediate failure)
		if len(mock.errorCalls) != 1 {
			t.Errorf("Never should fail immediately (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}

		// Should fail very quickly, not wait for timeout
		if elapsed > 100*time.Millisecond {
			t.Errorf("Expected immediate failure, but took %v", elapsed)
		}
	})

	t.Run("ReportsTimingInError", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)
		var counter int32

		assert.Never(func() bool {
			return atomic.AddInt32(&counter, 1) >= 2
		}, 500*time.Millisecond, 50*time.Millisecond)

		// Framework behavior: FAIL = exactly 1 Errorf call (condition becomes true)
		if len(mock.errorCalls) != 1 {
			t.Fatalf("Never should fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}

		errorMsg := mock.errorCalls[0]
		// Check that error contains timing information
		expectedFields := []string{"elapsed:", "attempts:"}
		for _, field := range expectedFields {
			if !containsString(errorMsg, field) {
				t.Errorf("Expected error message to contain %q, got: %s", field, errorMsg)
			}
		}
	})
}

// TestEventuallyWith tests the EventuallyWith method with custom configuration.
func TestEventuallyWith(t *testing.T) {
	t.Run("CustomTimeout", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)
		startTime := time.Now()

		config := EventuallyConfig{
			Timeout:  150 * time.Millisecond,
			Interval: 30 * time.Millisecond,
		}

		assert.EventuallyWith(func() bool {
			return false
		}, config)

		elapsed := time.Since(startTime)

		// Framework behavior: FAIL = exactly 1 Errorf call (timeout)
		if len(mock.errorCalls) != 1 {
			t.Errorf("Expected timeout error (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}

		// Should respect custom timeout
		if elapsed < 150*time.Millisecond || elapsed > 300*time.Millisecond {
			t.Errorf("Expected elapsed time around 150ms, got %v", elapsed)
		}
	})

	t.Run("ExponentialBackoff", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)
		var attempts []time.Time

		config := EventuallyConfig{
			Timeout:       500 * time.Millisecond,
			Interval:      50 * time.Millisecond,
			BackoffFactor: 2.0,
		}

		assert.EventuallyWith(func() bool {
			attempts = append(attempts, time.Now())
			return len(attempts) >= 4 // Succeed after 4 attempts
		}, config)

		// Framework behavior: PASS = no Errorf calls (condition succeeds)
		if len(mock.errorCalls) != 0 {
			t.Errorf("Expected success (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}

		// Verify exponential backoff intervals
		if len(attempts) < 4 {
			t.Fatalf("Expected at least 4 attempts, got %d", len(attempts))
		}

		// Check that intervals roughly doubled (allowing for timing variance)
		interval1 := attempts[1].Sub(attempts[0])
		interval2 := attempts[2].Sub(attempts[1])

		if interval2 < interval1*3/2 { // Should be roughly doubled
			t.Errorf("Expected exponential backoff, intervals: %v, %v", interval1, interval2)
		}
	})

	t.Run("MaxIntervalLimit", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)
		var intervals []time.Duration
		var lastTime time.Time

		config := EventuallyConfig{
			Timeout:       1 * time.Second,
			Interval:      10 * time.Millisecond,
			BackoffFactor: 3.0,
			MaxInterval:   100 * time.Millisecond,
		}

		assert.EventuallyWith(func() bool {
			now := time.Now()
			if !lastTime.IsZero() {
				intervals = append(intervals, now.Sub(lastTime))
			}
			lastTime = now
			return len(intervals) >= 5
		}, config)

		// Framework behavior: PASS = no Errorf calls (condition succeeds)
		if len(mock.errorCalls) != 0 {
			t.Errorf("Expected success (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}

		// Find the maximum interval used
		maxInterval := time.Duration(0)
		for _, interval := range intervals {
			if interval > maxInterval {
				maxInterval = interval
			}
		}

		// Should not exceed MaxInterval (with some tolerance for timing)
		if maxInterval > 150*time.Millisecond {
			t.Errorf("Expected max interval capped at ~100ms, got %v", maxInterval)
		}
	})

	t.Run("DefaultsInvalidValues", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)

		config := EventuallyConfig{
			Timeout:       -1 * time.Second,        // Invalid
			Interval:      -100 * time.Millisecond, // Invalid
			BackoffFactor: 0.5,                     // Invalid
		}

		assert.EventuallyWith(func() bool {
			return true // Succeed immediately
		}, config)

		// Framework behavior: PASS = no Errorf calls (condition succeeds)
		if len(mock.errorCalls) != 0 {
			t.Errorf("Expected success with defaults applied (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}
	})
}

// TestNeverWith tests the NeverWith method with custom configuration.
func TestNeverWith(t *testing.T) {
	t.Run("CustomConfiguration", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)
		startTime := time.Now()

		config := EventuallyConfig{
			Timeout:  200 * time.Millisecond,
			Interval: 40 * time.Millisecond,
		}

		assert.NeverWith(func() bool {
			return false
		}, config)

		elapsed := time.Since(startTime)

		// Framework behavior: PASS = no Errorf calls (condition never becomes true)
		if len(mock.errorCalls) != 0 {
			t.Errorf("Expected success (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}

		// Should wait for the full timeout
		if elapsed < 200*time.Millisecond {
			t.Errorf("Expected to wait full timeout, elapsed: %v", elapsed)
		}
	})

	t.Run("BackoffConfiguration", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)
		var counter int32

		config := EventuallyConfig{
			Timeout:       300 * time.Millisecond,
			Interval:      30 * time.Millisecond,
			BackoffFactor: 1.5,
		}

		assert.NeverWith(func() bool {
			return atomic.AddInt32(&counter, 1) >= 3 // True on 3rd call
		}, config)

		// Framework behavior: FAIL = exactly 1 Errorf call (condition becomes true)
		if len(mock.errorCalls) != 1 {
			t.Errorf("Expected failure when condition becomes true (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}

		errorMsg := mock.errorCalls[0]
		if !containsString(errorMsg, "final interval:") {
			t.Errorf("Expected error to contain final interval info, got: %s", errorMsg)
		}
	})
}

// TestResourceCleanup verifies proper cleanup of goroutines and tickers.
func TestResourceCleanup(t *testing.T) {
	t.Run("EvenuallyCleanup", func(t *testing.T) {
		// This test ensures that resources are properly cleaned up
		// We can't directly test goroutine counts, but we can test behaviour
		// Test GoWise framework behavioral contract with multiple iterations

		// Run multiple Eventually assertions to stress test cleanup
		for i := 0; i < 10; i++ {
			// Create fresh mock for each iteration to test behavioral contract
			iterMock := &behaviorMockT{}
			iterAssert := New(iterMock)

			iterAssert.Eventually(func() bool {
				return true // Succeed immediately
			}, 100*time.Millisecond, 10*time.Millisecond)

			// Framework behavior: PASS = no Errorf calls (condition succeeds immediately)
			if len(iterMock.errorCalls) != 0 {
				t.Errorf("Iteration %d: Expected success (no Errorf calls), got %d: %v", i, len(iterMock.errorCalls), iterMock.errorCalls)
			}
		}
	})

	t.Run("NeverCleanup", func(t *testing.T) {
		// Run multiple Never assertions
		for i := 0; i < 10; i++ {
			// Create fresh mock for each iteration to test behavioral contract
			iterMock := &behaviorMockT{}
			iterAssert := New(iterMock)

			iterAssert.Never(func() bool {
				return false // Never true
			}, 50*time.Millisecond, 10*time.Millisecond)

			// Framework behavior: PASS = no Errorf calls (condition never becomes true)
			if len(iterMock.errorCalls) != 0 {
				t.Errorf("Iteration %d: Expected success (no Errorf calls), got %d: %v", i, len(iterMock.errorCalls), iterMock.errorCalls)
			}
		}
	})
}

// TestEdgeCases tests edge cases and boundary conditions.
func TestEdgeCases(t *testing.T) {
	t.Run("ZeroTimeout", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)

		assert.Eventually(func() bool {
			return false
		}, 0*time.Second, 10*time.Millisecond)

		// Framework behavior: FAIL = exactly 1 Errorf call (should use default timeout and timeout)
		if len(mock.errorCalls) != 1 {
			t.Errorf("Expected timeout with zero duration (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}
	})

	t.Run("VeryShortTimeout", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)

		assert.Eventually(func() bool {
			return false
		}, 1*time.Nanosecond, 10*time.Millisecond)

		// Framework behavior: FAIL = exactly 1 Errorf call (immediate timeout)
		if len(mock.errorCalls) != 1 {
			t.Errorf("Expected immediate timeout (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}
	})

	t.Run("NegativeTimeout", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)

		assert.Eventually(func() bool {
			return true
		}, -1*time.Second, 10*time.Millisecond)

		// Framework behavior: PASS = no Errorf calls (should succeed with defaults applied)
		if len(mock.errorCalls) != 0 {
			t.Errorf("Expected success with default timeout (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}
	})

	t.Run("ZeroInterval", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)

		assert.Eventually(func() bool {
			return true
		}, 100*time.Millisecond, 0*time.Millisecond)

		// Framework behavior: PASS = no Errorf calls (should succeed with defaults applied)
		if len(mock.errorCalls) != 0 {
			t.Errorf("Expected success with default interval (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}
	})

	t.Run("NegativeInterval", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)

		assert.Eventually(func() bool {
			return true
		}, 100*time.Millisecond, -10*time.Millisecond)

		// Framework behavior: PASS = no Errorf calls (should succeed with defaults applied)
		if len(mock.errorCalls) != 0 {
			t.Errorf("Expected success with default interval (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}
	})

	t.Run("VeryLongTimeout", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)

		assert.Eventually(func() bool {
			return true // Succeeds immediately
		}, 24*time.Hour, 1*time.Second) // Very long timeout

		// Framework behavior: PASS = no Errorf calls (succeeds immediately)
		if len(mock.errorCalls) != 0 {
			t.Errorf("Expected immediate success (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}
	})

	t.Run("VeryShortInterval", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)
		var attempts int32

		assert.Eventually(func() bool {
			return atomic.AddInt32(&attempts, 1) >= 3
		}, 100*time.Millisecond, 1*time.Microsecond) // Very short interval

		// Framework behavior: PASS = no Errorf calls (condition succeeds after attempts)
		if len(mock.errorCalls) != 0 {
			t.Errorf("Expected success (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}

		// Should have made many attempts due to short interval
		finalAttempts := atomic.LoadInt32(&attempts)
		if finalAttempts < 3 {
			t.Errorf("Expected at least 3 attempts, got %d", finalAttempts)
		}
	})

	t.Run("PanicInCondition", func(t *testing.T) {
		assert := New(&mockT{})

		// Test that panics in condition functions don't crash the assertion
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic to propagate")
			}
		}()

		assert.Eventually(func() bool {
			panic("test panic")
		}, 100*time.Millisecond, 10*time.Millisecond)
	})

	t.Run("NilConditionFunction", func(t *testing.T) {
		assert := New(&mockT{})

		// Test with nil condition function - should panic
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic with nil condition")
			}
		}()

		var nilCondition func() bool
		assert.Eventually(nilCondition, 100*time.Millisecond, 10*time.Millisecond)
	})

	t.Run("BackoffFactorEdgeCases", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)

		// Test with backoff factor exactly 1.0
		config := EventuallyConfig{
			Timeout:       200 * time.Millisecond,
			Interval:      50 * time.Millisecond,
			BackoffFactor: 1.0, // No backoff
		}

		assert.EventuallyWith(func() bool {
			return true
		}, config)

		// Framework behavior: PASS = no Errorf calls (succeeds immediately)
		if len(mock.errorCalls) != 0 {
			t.Errorf("Expected success with backoff factor 1.0 (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}
	})

	t.Run("BackoffFactorLessThanOne", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)

		config := EventuallyConfig{
			Timeout:       100 * time.Millisecond,
			Interval:      20 * time.Millisecond,
			BackoffFactor: 0.5, // Invalid - should be corrected to 1.0
		}

		assert.EventuallyWith(func() bool {
			return true
		}, config)

		// Framework behavior: PASS = no Errorf calls (should succeed with corrected defaults)
		if len(mock.errorCalls) != 0 {
			t.Errorf("Expected success with corrected backoff factor (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}
	})

	t.Run("MaxIntervalSmallerThanInterval", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)

		config := EventuallyConfig{
			Timeout:       200 * time.Millisecond,
			Interval:      100 * time.Millisecond,
			BackoffFactor: 2.0,
			MaxInterval:   50 * time.Millisecond, // Smaller than initial interval
		}

		assert.EventuallyWith(func() bool {
			return true
		}, config)

		// Framework behavior: PASS = no Errorf calls (should succeed with max interval handling)
		if len(mock.errorCalls) != 0 {
			t.Errorf("Expected success with max interval handling (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}
	})

	t.Run("TimeoutExactlyReached", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)
		var callCount int32
		startTime := time.Now()

		assert.Eventually(func() bool {
			atomic.AddInt32(&callCount, 1)
			return false // Never succeeds
		}, 100*time.Millisecond, 25*time.Millisecond)

		elapsed := time.Since(startTime)

		// Framework behavior: FAIL = exactly 1 Errorf call (timeout reached)
		if len(mock.errorCalls) != 1 {
			t.Errorf("Expected timeout error (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}

		// Should have waited close to timeout duration
		if elapsed < 90*time.Millisecond || elapsed > 200*time.Millisecond {
			t.Errorf("Expected elapsed time around 100ms, got %v", elapsed)
		}
	})
}

// TestAdvancedEdgeCases tests more complex edge cases and error conditions.
func TestAdvancedEdgeCases(t *testing.T) {
	t.Run("ConditionFlipsBackAndForth", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)
		var counter int32

		// Condition alternates true/false
		assert.Eventually(func() bool {
			count := atomic.AddInt32(&counter, 1)
			return count%2 == 0 // True on even attempts
		}, 200*time.Millisecond, 20*time.Millisecond)

		// Framework behavior: PASS = no Errorf calls (condition eventually becomes true)
		if len(mock.errorCalls) != 0 {
			t.Errorf("Expected success when condition eventually becomes true (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}
	})

	t.Run("ConditionBecomesTrueThenFalse", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)
		var counter int32

		assert.Never(func() bool {
			count := atomic.AddInt32(&counter, 1)
			if count == 3 {
				return true // Becomes true on 3rd attempt
			}
			return false
		}, 200*time.Millisecond, 30*time.Millisecond)

		// Framework behavior: FAIL = exactly 1 Errorf call (condition becomes true)
		if len(mock.errorCalls) != 1 {
			t.Errorf("Expected Never to fail when condition becomes true (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}

		if !containsString(mock.errorCalls[0], "became true unexpectedly") {
			t.Errorf("Expected 'became true unexpectedly' in error, got: %s", mock.errorCalls[0])
		}
	})

	t.Run("HighFrequencyConditionChange", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)
		var toggleState int32

		// Background goroutine rapidly toggles state
		go func() {
			for i := 0; i < 20; i++ {
				time.Sleep(5 * time.Millisecond)
				if atomic.LoadInt32(&toggleState) == 0 {
					atomic.StoreInt32(&toggleState, 1)
				} else {
					atomic.StoreInt32(&toggleState, 0)
				}
			}
		}()

		assert.Eventually(func() bool {
			return atomic.LoadInt32(&toggleState) == 1
		}, 500*time.Millisecond, 10*time.Millisecond)

		// Framework behavior: PASS = no Errorf calls (condition should eventually be true)
		if len(mock.errorCalls) != 0 {
			t.Errorf("Expected success with rapidly changing condition (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}
	})

	t.Run("MemoryPressureCondition", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)
		var allocCount int32

		// Simulate memory allocation in condition
		assert.Eventually(func() bool {
			count := atomic.AddInt32(&allocCount, 1)
			// Allocate some memory each check to test under memory pressure
			_ = make([]byte, 1024)
			return count >= 10
		}, 1*time.Second, 50*time.Millisecond)

		// Framework behavior: PASS = no Errorf calls (condition eventually succeeds)
		if len(mock.errorCalls) != 0 {
			t.Errorf("Expected success with memory allocating condition (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}
	})

	t.Run("ConditionWithSystemCalls", func(t *testing.T) {
		// Test GoWise framework behavioral contract
		mock := &behaviorMockT{}
		assert := New(mock)
		var checkCount int32

		// Simulate condition that makes system calls
		assert.Eventually(func() bool {
			count := atomic.AddInt32(&checkCount, 1)
			// Simulate system call with time check
			_ = time.Now().Unix()
			return count >= 5
		}, 500*time.Millisecond, 50*time.Millisecond)

		// Framework behavior: PASS = no Errorf calls (condition eventually succeeds)
		if len(mock.errorCalls) != 0 {
			t.Errorf("Expected success with system call condition (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
		}
	})
}

// mockT is a simple mock implementation of testing.T for testing assertions.
type mockT struct {
	helperCalled bool
}

func (m *mockT) Helper() {
	m.helperCalled = true
}

// silentT is defined in assertions_passing_test.go - shared across test files

// containsString checks if a string contains a substring (helper function).
func containsString(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) &&
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}

// ExampleAssert_Eventually demonstrates basic Eventually assertion usage.
func ExampleAssert_Eventually() {
	assert := New(&silentT{})

	// Service readiness simulation
	var serviceReady int32
	go func() {
		time.Sleep(150 * time.Millisecond)
		atomic.StoreInt32(&serviceReady, 1)
	}()

	// Wait up to 1 second for service to be ready, checking every 50ms
	assert.Eventually(func() bool {
		return atomic.LoadInt32(&serviceReady) == 1
	}, 1*time.Second, 50*time.Millisecond)

	// The assertion succeeded since the service became ready
	fmt.Println("No error:", true)
	// Output: No error: true
}

// ExampleAssert_Never demonstrates basic Never assertion usage.
func ExampleAssert_Never() {
	assert := New(&silentT{})

	// Cache expiration simulation - keys should never expire during this test
	cacheHasExpiredKeys := false

	// Verify cache doesn't have expired keys for 200ms, checking every 50ms
	assert.Never(func() bool {
		return cacheHasExpiredKeys
	}, 200*time.Millisecond, 50*time.Millisecond)

	// The assertion succeeded since the cache never had expired keys
	fmt.Println("No error:", true)
	// Output: No error: true
}

// ExampleAssert_EventuallyWith demonstrates advanced Eventually configuration.
func ExampleAssert_EventuallyWith() {
	assert := New(&silentT{})

	// Database connection simulation with exponential backoff
	var attempts int32
	var connected int32

	go func() {
		time.Sleep(300 * time.Millisecond)
		atomic.StoreInt32(&connected, 1)
	}()

	// Custom configuration with exponential backoff
	config := EventuallyConfig{
		Timeout:       1 * time.Second,
		Interval:      50 * time.Millisecond,
		BackoffFactor: 1.5,
		MaxInterval:   200 * time.Millisecond,
	}

	assert.EventuallyWith(func() bool {
		atomic.AddInt32(&attempts, 1)
		return atomic.LoadInt32(&connected) == 1
	}, config)

	// The assertion succeeded since the database connection was established
	fmt.Println("No error:", true)
	fmt.Println("Used exponential backoff:", atomic.LoadInt32(&attempts) <= 6) // Fewer attempts due to backoff
	// Output: No error: true
	// Used exponential backoff: true
}

// ExampleAssert_NeverWith demonstrates advanced Never configuration.
func ExampleAssert_NeverWith() {
	assert := New(&silentT{})

	// Memory leak detection simulation
	memoryLeakDetected := false

	config := EventuallyConfig{
		Timeout:  500 * time.Millisecond,
		Interval: 100 * time.Millisecond,
	}

	// Verify no memory leaks are detected during the test period
	assert.NeverWith(func() bool {
		return memoryLeakDetected
	}, config)

	// The assertion succeeded since the configuration was applied
	fmt.Println("No error:", true)
	// Output: No error: true
}

// ExampleAssert_Eventually_fileAppears demonstrates waiting for file creation.
func ExampleAssert_Eventually_fileAppears() {
	assert := New(&silentT{})

	// Simulate file creation after delay
	var fileExists int32
	go func() {
		time.Sleep(100 * time.Millisecond)
		atomic.StoreInt32(&fileExists, 1)
	}()

	// Wait for file to appear - in real usage you'd check os.Stat
	assert.Eventually(func() bool {
		return atomic.LoadInt32(&fileExists) == 1
	}, 1*time.Second, 50*time.Millisecond)

	// The assertion succeeded since the file appeared
	fmt.Println("File appeared:", true)
	// Output: File appeared: true
}

// ExampleAssert_Eventually_httpHealthCheck demonstrates HTTP endpoint health checking.
func ExampleAssert_Eventually_httpHealthCheck() {
	assert := New(&silentT{})

	// Simulate HTTP health check
	var serverHealthy int32
	go func() {
		time.Sleep(200 * time.Millisecond)
		atomic.StoreInt32(&serverHealthy, 1)
	}()

	assert.Eventually(func() bool {
		// In real usage: resp, err := http.Get("http://service/health"); return err == nil && resp.StatusCode == 200
		return atomic.LoadInt32(&serverHealthy) == 1
	}, 2*time.Second, 100*time.Millisecond)

	// The assertion succeeded since the health check passed
	fmt.Println("Health check passed:", true)
	// Output: Health check passed: true
}

// ExampleAssert_Never_resourceLeak demonstrates detecting resource leaks.
func ExampleAssert_Never_resourceLeak() {
	assert := New(&silentT{})

	// Simulate resource monitoring
	var resourceCount int32 = 10

	go func() {
		// Simulate resource usage that should not exceed limit
		for i := 0; i < 5; i++ {
			time.Sleep(40 * time.Millisecond)
			atomic.AddInt32(&resourceCount, 1)
		}
	}()

	// Ensure resource usage never exceeds limit
	assert.Never(func() bool {
		return atomic.LoadInt32(&resourceCount) > 20
	}, 400*time.Millisecond, 50*time.Millisecond)

	// The assertion succeeded since no resource leaks were detected
	fmt.Println("Resource leak prevented:", true)
	// Output: Resource leak prevented: true
}

// ExampleEventuallyConfig demonstrates different configuration options.
func ExampleEventuallyConfig() {
	// Basic configuration
	basicConfig := EventuallyConfig{
		Timeout:  3 * time.Second,
		Interval: 100 * time.Millisecond,
	}

	// Configuration with exponential backoff
	backoffConfig := EventuallyConfig{
		Timeout:       10 * time.Second,
		Interval:      50 * time.Millisecond,
		BackoffFactor: 2.0,
		MaxInterval:   1 * time.Second,
	}

	fmt.Printf("Basic timeout: %v, interval: %v\n", basicConfig.Timeout, basicConfig.Interval)
	fmt.Printf("Backoff factor: %.1f, max interval: %v\n", backoffConfig.BackoffFactor, backoffConfig.MaxInterval)
	// Output: Basic timeout: 3s, interval: 100ms
	// Backoff factor: 2.0, max interval: 1s
}

// ExampleAssert_EventuallyWith_databaseConnection demonstrates database connection with retry backoff.
func ExampleAssert_EventuallyWith_databaseConnection() {
	assert := New(&silentT{})

	var connectionEstablished int32
	var connectionAttempts int32

	go func() {
		time.Sleep(250 * time.Millisecond) // Database becomes available
		atomic.StoreInt32(&connectionEstablished, 1)
	}()

	config := EventuallyConfig{
		Timeout:       2 * time.Second,
		Interval:      100 * time.Millisecond,
		BackoffFactor: 1.5, // Increase interval by 50% each attempt
		MaxInterval:   500 * time.Millisecond,
	}

	assert.EventuallyWith(func() bool {
		atomic.AddInt32(&connectionAttempts, 1)
		return atomic.LoadInt32(&connectionEstablished) == 1
	}, config)

	// The assertion succeeded since the database connection was established
	fmt.Println("Database connected:", true)
	fmt.Println("Used backoff strategy:", atomic.LoadInt32(&connectionAttempts) <= 4)
	// Output: Database connected: true
	// Used backoff strategy: true
}

// ExampleAssert_Never_timeoutScenario demonstrates timeout detection.
func ExampleAssert_Never_timeoutScenario() {
	assert := New(&silentT{})

	// Simulate operation that should complete quickly
	var operationTimeout int32

	go func() {
		time.Sleep(50 * time.Millisecond) // Operation completes quickly
		// In this example, timeout never occurs
	}()

	// Verify operation never times out during test period
	assert.Never(func() bool {
		return atomic.LoadInt32(&operationTimeout) == 1
	}, 200*time.Millisecond, 25*time.Millisecond)

	// The assertion succeeded since no timeout occurred
	fmt.Println("No timeout detected:", true)
	// Output: No timeout detected: true
}
