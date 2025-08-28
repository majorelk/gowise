package assertions

import (
	"sync/atomic"
	"testing"
	"time"
)

// ContractEventually runs behavioural checks for any Eventually implementation.
// This follows GoWise contract testing patterns to ensure all implementations
// satisfy the same behavioural guarantees.
func ContractEventually(t *testing.T, name string, newAssert func(t *testing.T) *Assert) {
	t.Run(name+"/immediate_success", func(t *testing.T) {
		assert := newAssert(t)
		assert.Eventually(func() bool { return true }, 1*time.Second, 50*time.Millisecond)
		if assert.Error() != "" {
			t.Fatalf("behaviour: expected immediate success, got error: %s", assert.Error())
		}
	})

	t.Run(name+"/delayed_success", func(t *testing.T) {
		assert := newAssert(t)
		var ready int32

		go func() {
			time.Sleep(100 * time.Millisecond)
			atomic.StoreInt32(&ready, 1)
		}()

		assert.Eventually(func() bool {
			return atomic.LoadInt32(&ready) == 1
		}, 500*time.Millisecond, 25*time.Millisecond)

		if assert.Error() != "" {
			t.Fatalf("behaviour: expected delayed success, got error: %s", assert.Error())
		}
	})

	t.Run(name+"/timeout_failure", func(t *testing.T) {
		assert := newAssert(t)
		assert.Eventually(func() bool { return false }, 100*time.Millisecond, 25*time.Millisecond)
		if assert.Error() == "" {
			t.Fatalf("behaviour: expected timeout error, got success")
		}
		if !containsString(assert.Error(), "condition not met within timeout") {
			t.Fatalf("behaviour: expected timeout message, got: %s", assert.Error())
		}
	})

	t.Run(name+"/resource_cleanup", func(t *testing.T) {
		// Test multiple assertions to ensure proper cleanup
		for i := 0; i < 5; i++ {
			assert := newAssert(t)
			assert.Eventually(func() bool { return true }, 100*time.Millisecond, 10*time.Millisecond)
			if assert.Error() != "" {
				t.Fatalf("behaviour: cleanup test iteration %d failed: %s", i, assert.Error())
			}
		}
	})

	t.Run(name+"/condition_evaluation_order", func(t *testing.T) {
		assert := newAssert(t)
		var callCount int32

		assert.Eventually(func() bool {
			count := atomic.AddInt32(&callCount, 1)
			return count >= 3
		}, 500*time.Millisecond, 50*time.Millisecond)

		if assert.Error() != "" {
			t.Fatalf("behaviour: expected success after 3 calls, got error: %s", assert.Error())
		}

		finalCount := atomic.LoadInt32(&callCount)
		if finalCount != 3 {
			t.Fatalf("behaviour: expected exactly 3 calls, got %d", finalCount)
		}
	})
}

// ContractNever runs behavioural checks for any Never implementation.
func ContractNever(t *testing.T, name string, newAssert func(t *testing.T) *Assert) {
	t.Run(name+"/never_true_success", func(t *testing.T) {
		assert := newAssert(t)
		assert.Never(func() bool { return false }, 200*time.Millisecond, 40*time.Millisecond)
		if assert.Error() != "" {
			t.Fatalf("behaviour: expected success when condition never true, got: %s", assert.Error())
		}
	})

	t.Run(name+"/immediate_failure", func(t *testing.T) {
		assert := newAssert(t)
		assert.Never(func() bool { return true }, 200*time.Millisecond, 40*time.Millisecond)
		if assert.Error() == "" {
			t.Fatalf("behaviour: expected immediate failure when condition is true")
		}
		if !containsString(assert.Error(), "became true unexpectedly") {
			t.Fatalf("behaviour: expected 'became true unexpectedly' message, got: %s", assert.Error())
		}
	})

	t.Run(name+"/delayed_failure", func(t *testing.T) {
		assert := newAssert(t)
		var shouldFail int32

		go func() {
			time.Sleep(100 * time.Millisecond)
			atomic.StoreInt32(&shouldFail, 1)
		}()

		assert.Never(func() bool {
			return atomic.LoadInt32(&shouldFail) == 1
		}, 500*time.Millisecond, 30*time.Millisecond)

		if assert.Error() == "" {
			t.Fatalf("behaviour: expected failure when condition becomes true")
		}
	})

	t.Run(name+"/timing_accuracy", func(t *testing.T) {
		assert := newAssert(t)
		startTime := time.Now()

		assert.Never(func() bool { return false }, 150*time.Millisecond, 25*time.Millisecond)

		elapsed := time.Since(startTime)
		if assert.Error() != "" {
			t.Fatalf("behaviour: expected success, got error: %s", assert.Error())
		}

		// Should wait close to the full timeout
		if elapsed < 140*time.Millisecond || elapsed > 200*time.Millisecond {
			t.Fatalf("behaviour: expected ~150ms elapsed, got %v", elapsed)
		}
	})
}

// ContractEventuallyWith runs behavioural checks for EventuallyWith with various configurations.
func ContractEventuallyWith(t *testing.T, name string, newAssert func(t *testing.T) *Assert) {
	t.Run(name+"/basic_config", func(t *testing.T) {
		assert := newAssert(t)
		config := EventuallyConfig{
			Timeout:  300 * time.Millisecond,
			Interval: 50 * time.Millisecond,
		}

		assert.EventuallyWith(func() bool { return true }, config)
		if assert.Error() != "" {
			t.Fatalf("behaviour: expected success with basic config, got: %s", assert.Error())
		}
	})

	t.Run(name+"/exponential_backoff_behaviour", func(t *testing.T) {
		assert := newAssert(t)
		var attempts []time.Time

		config := EventuallyConfig{
			Timeout:       400 * time.Millisecond,
			Interval:      40 * time.Millisecond,
			BackoffFactor: 2.0,
		}

		assert.EventuallyWith(func() bool {
			attempts = append(attempts, time.Now())
			return len(attempts) >= 4
		}, config)

		if assert.Error() != "" {
			t.Fatalf("behaviour: expected success with backoff, got: %s", assert.Error())
		}

		if len(attempts) < 4 {
			t.Fatalf("behaviour: expected at least 4 attempts, got %d", len(attempts))
		}

		// Verify exponential backoff (later intervals should be longer)
		if len(attempts) >= 3 {
			interval1 := attempts[1].Sub(attempts[0])
			interval2 := attempts[2].Sub(attempts[1])

			// Second interval should be roughly twice the first (within tolerance)
			if interval2 < interval1*3/2 {
				t.Fatalf("behaviour: expected exponential backoff, intervals: %v, %v", interval1, interval2)
			}
		}
	})

	t.Run(name+"/max_interval_enforcement", func(t *testing.T) {
		assert := newAssert(t)
		var attempts []time.Time

		config := EventuallyConfig{
			Timeout:       600 * time.Millisecond,
			Interval:      20 * time.Millisecond,
			BackoffFactor: 3.0,
			MaxInterval:   80 * time.Millisecond,
		}

		assert.EventuallyWith(func() bool {
			attempts = append(attempts, time.Now())
			return len(attempts) >= 6
		}, config)

		if assert.Error() != "" {
			t.Fatalf("behaviour: expected success with max interval, got: %s", assert.Error())
		}

		// Verify no interval exceeds MaxInterval (with tolerance)
		for i := 1; i < len(attempts); i++ {
			interval := attempts[i].Sub(attempts[i-1])
			if interval > 120*time.Millisecond { // MaxInterval + tolerance
				t.Fatalf("behaviour: interval %v exceeded MaxInterval constraint", interval)
			}
		}
	})

	t.Run(name+"/invalid_config_defaults", func(t *testing.T) {
		assert := newAssert(t)
		config := EventuallyConfig{
			Timeout:       -1 * time.Second,        // Invalid
			Interval:      -100 * time.Millisecond, // Invalid
			BackoffFactor: 0.5,                     // Invalid
		}

		assert.EventuallyWith(func() bool { return true }, config)
		if assert.Error() != "" {
			t.Fatalf("behaviour: expected success with corrected defaults, got: %s", assert.Error())
		}
	})
}

// ContractErrorReporting runs behavioural checks for error message quality.
func ContractErrorReporting(t *testing.T, name string, newAssert func(t *testing.T) *Assert) {
	t.Run(name+"/timeout_error_content", func(t *testing.T) {
		assert := newAssert(t)
		assert.Eventually(func() bool { return false }, 100*time.Millisecond, 20*time.Millisecond)

		errorMsg := assert.Error()
		if errorMsg == "" {
			t.Fatalf("behaviour: expected error message for timeout")
		}

		// Error should contain timing information
		requiredFields := []string{"timeout:", "elapsed:", "attempts:"}
		for _, field := range requiredFields {
			if !containsString(errorMsg, field) {
				t.Fatalf("behaviour: error message missing '%s': %s", field, errorMsg)
			}
		}
	})

	t.Run(name+"/never_error_content", func(t *testing.T) {
		assert := newAssert(t)
		assert.Never(func() bool { return true }, 100*time.Millisecond, 20*time.Millisecond)

		errorMsg := assert.Error()
		if errorMsg == "" {
			t.Fatalf("behaviour: expected error message for Never failure")
		}

		requiredFields := []string{"became true unexpectedly", "elapsed:", "attempts:"}
		for _, field := range requiredFields {
			if !containsString(errorMsg, field) {
				t.Fatalf("behaviour: error message missing '%s': %s", field, errorMsg)
			}
		}
	})

	t.Run(name+"/clear_success_state", func(t *testing.T) {
		assert := newAssert(t)
		assert.Eventually(func() bool { return true }, 100*time.Millisecond, 20*time.Millisecond)

		if assert.Error() != "" {
			t.Fatalf("behaviour: expected empty error on success, got: %s", assert.Error())
		}
	})
}

// TestAssertContractCompliance runs contract tests against the actual implementation.
func TestAssertContractCompliance(t *testing.T) {
	// Test our implementation against the contracts
	ContractEventually(t, "Assert", func(t *testing.T) *Assert {
		return New(&mockT{})
	})

	ContractNever(t, "Assert", func(t *testing.T) *Assert {
		return New(&mockT{})
	})

	ContractEventuallyWith(t, "Assert", func(t *testing.T) *Assert {
		return New(&mockT{})
	})

	ContractErrorReporting(t, "Assert", func(t *testing.T) *Assert {
		return New(&mockT{})
	})
}

// TestEventuallyConfigContract validates EventuallyConfig behaviour invariants.
func TestEventuallyConfigContract(t *testing.T) {
	t.Run("DefaultsValidation", func(t *testing.T) {
		defaults := defaultEventuallyConfig()

		if defaults.Timeout <= 0 {
			t.Errorf("Default timeout should be positive, got %v", defaults.Timeout)
		}

		if defaults.Interval <= 0 {
			t.Errorf("Default interval should be positive, got %v", defaults.Interval)
		}

		if defaults.BackoffFactor < 1.0 {
			t.Errorf("Default backoff factor should be >= 1.0, got %v", defaults.BackoffFactor)
		}
	})

	t.Run("ConfigurationImmutability", func(t *testing.T) {
		original := EventuallyConfig{
			Timeout:       1 * time.Second,
			Interval:      100 * time.Millisecond,
			BackoffFactor: 1.5,
			MaxInterval:   500 * time.Millisecond,
		}

		// Copy for comparison
		copy := original

		// Use config in assertion
		assert := New(&mockT{})
		assert.EventuallyWith(func() bool { return true }, original)

		// Verify original wasn't modified
		if original.Timeout != copy.Timeout ||
			original.Interval != copy.Interval ||
			original.BackoffFactor != copy.BackoffFactor ||
			original.MaxInterval != copy.MaxInterval {
			t.Error("EventuallyConfig should not be modified during use")
		}
	})

	t.Run("ReasonablePerformanceConstraints", func(t *testing.T) {
		assert := New(&mockT{})

		// Test that reasonable configurations don't cause excessive overhead
		config := EventuallyConfig{
			Timeout:  200 * time.Millisecond,
			Interval: 10 * time.Millisecond,
		}

		start := time.Now()
		assert.EventuallyWith(func() bool { return true }, config)
		duration := time.Since(start)

		// Should complete almost immediately for successful condition
		if duration > 50*time.Millisecond {
			t.Errorf("Immediate success should be fast, took %v", duration)
		}
	})
}

// TestConcurrentContractCompliance verifies thread-safety contracts.
func TestConcurrentContractCompliance(t *testing.T) {
	t.Run("ThreadSafetyContract", func(t *testing.T) {
		const numGoroutines = 10
		const assertionsPerGoroutine = 5

		done := make(chan bool, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer func() { done <- true }()

				for j := 0; j < assertionsPerGoroutine; j++ {
					assert := New(&mockT{})
					var ready int32

					go func() {
						time.Sleep(time.Duration(id*10) * time.Millisecond)
						atomic.StoreInt32(&ready, 1)
					}()

					assert.Eventually(func() bool {
						return atomic.LoadInt32(&ready) == 1
					}, 500*time.Millisecond, 20*time.Millisecond)

					if assert.Error() != "" {
						t.Errorf("Goroutine %d assertion %d failed: %s", id, j, assert.Error())
					}
				}
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < numGoroutines; i++ {
			<-done
		}
	})
}
