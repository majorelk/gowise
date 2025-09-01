package assertions

import (
	"fmt"
	"testing"
)

// behaviorMockT is defined in assertions_passing_test.go

// TestBehavioralChaining tests chaining behavior through the TestingT interface (not internal state)
func TestBehavioralChaining(t *testing.T) {
	t.Run("SuccessfulChainCallsNoTestingTMethods", func(t *testing.T) {
		mock := &behaviorMockT{}
		assert := New(mock)

		// Execute a successful chain - this should NOT call any TestingT methods
		assert.Equal(42, 42).True(true).Contains("hello", "ell").Len("test", 4)

		// BEHAVIORAL TEST: Verify no TestingT methods were called
		if len(mock.errorCalls) != 0 {
			t.Errorf("Successful chain should not call Errorf, but got %d calls: %v", len(mock.errorCalls), mock.errorCalls)
		}
		if mock.failNowCalls != 0 {
			t.Errorf("Successful chain should not call FailNow, but got %d calls", mock.failNowCalls)
		}
	})

	t.Run("FailedAssertionCallsTestingTExactlyOnce", func(t *testing.T) {
		mock := &behaviorMockT{}
		assert := New(mock)

		// Execute a failing assertion - this SHOULD call TestingT.Errorf exactly once
		assert.Equal(1, 2) // This should fail

		// BEHAVIORAL TEST: Verify TestingT.Errorf was called exactly once
		if len(mock.errorCalls) != 1 {
			t.Errorf("Failed assertion should call Errorf exactly once, but got %d calls: %v", len(mock.errorCalls), mock.errorCalls)
		}
		if len(mock.errorCalls) > 0 && !containsSubstring(mock.errorCalls[0], "values differ") {
			t.Errorf("Error message should describe the difference, got: %q", mock.errorCalls[0])
		}
	})

	t.Run("ChainFailFastBehaviorCallsTestingTOnceOnly", func(t *testing.T) {
		mock := &behaviorMockT{}
		assert := New(mock)

		// Execute a chain where first assertion fails - only first should call TestingT.Errorf
		assert.Equal(1, 2).Equal(3, 4).Equal(5, 6)

		// BEHAVIORAL TEST: Verify TestingT.Errorf was called exactly once (fail-fast)
		if len(mock.errorCalls) != 1 {
			t.Errorf("Fail-fast chain should call Errorf exactly once, but got %d calls: %v", len(mock.errorCalls), mock.errorCalls)
		}

		// Verify the error is about the first failure, not subsequent ones
		if len(mock.errorCalls) > 0 {
			errorMsg := mock.errorCalls[0]
			if !containsSubstring(errorMsg, "1") || !containsSubstring(errorMsg, "2") {
				t.Errorf("Error should be about first failure (1 vs 2), got: %q", errorMsg)
			}
			if containsSubstring(errorMsg, "3") || containsSubstring(errorMsg, "4") {
				t.Errorf("Error should not contain subsequent failures (3,4,5,6), got: %q", errorMsg)
			}
		}
	})

	t.Run("DifferentAssertionTypesInChainCallTestingTOnFailure", func(t *testing.T) {
		mock := &behaviorMockT{}
		assert := New(mock)

		// Execute a mixed chain where Contains fails
		assert.Equal(42, 42).Contains("hello", "xyz").Len("test", 4) // Contains should fail

		// BEHAVIORAL TEST: Verify TestingT.Errorf was called exactly once
		if len(mock.errorCalls) != 1 {
			t.Errorf("Failed Contains in chain should call Errorf exactly once, but got %d calls: %v", len(mock.errorCalls), mock.errorCalls)
		}
		if len(mock.errorCalls) > 0 && !containsSubstring(mock.errorCalls[0], "hello") {
			t.Errorf("Contains error should mention the searched string, got: %q", mock.errorCalls[0])
		}
	})

	t.Run("ErrorHandlingChainsCallTestingTCorrectly", func(t *testing.T) {
		mock := &behaviorMockT{}
		assert := New(mock)

		// Execute error handling chain - HasError with nil should fail
		testErr := fmt.Errorf("test error")
		assert.HasError(testErr).HasError(nil) // Second should fail

		// BEHAVIORAL TEST: Should fail on second assertion
		if len(mock.errorCalls) != 1 {
			t.Errorf("HasError(nil) should call Errorf exactly once, but got %d calls: %v", len(mock.errorCalls), mock.errorCalls)
		}
	})

	t.Run("ToleranceMethodsInChainsCallTestingTOnFailure", func(t *testing.T) {
		mock := &behaviorMockT{}
		assert := New(mock)

		// Execute tolerance chain where WithinTolerance fails
		assert.WithinTolerance(1.0, 2.0, 0.5).True(true) // Should fail on tolerance

		// BEHAVIORAL TEST: Verify failure captured
		if len(mock.errorCalls) != 1 {
			t.Errorf("Failed WithinTolerance should call Errorf exactly once, but got %d calls: %v", len(mock.errorCalls), mock.errorCalls)
		}
		if len(mock.errorCalls) > 0 && !containsSubstring(mock.errorCalls[0], "tolerance") {
			t.Errorf("WithinTolerance error should mention tolerance, got: %q", mock.errorCalls[0])
		}
	})
}

// TestBehavioralIndividualAssertions tests that individual assertions work correctly through TestingT interface
func TestBehavioralIndividualAssertions(t *testing.T) {
	t.Run("EqualCallsTestingTOnMismatch", func(t *testing.T) {
		mock := &behaviorMockT{}
		assert := New(mock)

		assert.Equal("expected", "actual")

		if len(mock.errorCalls) != 1 {
			t.Errorf("Equal with different values should call Errorf once, got %d calls", len(mock.errorCalls))
		}
	})

	t.Run("TrueCallsTestingTOnFalse", func(t *testing.T) {
		mock := &behaviorMockT{}
		assert := New(mock)

		assert.True(false)

		if len(mock.errorCalls) != 1 {
			t.Errorf("True(false) should call Errorf once, got %d calls", len(mock.errorCalls))
		}
	})

	t.Run("ContainsCallsTestingTOnMissing", func(t *testing.T) {
		mock := &behaviorMockT{}
		assert := New(mock)

		assert.Contains("hello", "xyz")

		if len(mock.errorCalls) != 1 {
			t.Errorf("Contains with missing substring should call Errorf once, got %d calls", len(mock.errorCalls))
		}
	})

	t.Run("NoErrorCallsTestingTOnError", func(t *testing.T) {
		mock := &behaviorMockT{}
		assert := New(mock)

		assert.NoError(fmt.Errorf("test error"))

		if len(mock.errorCalls) != 1 {
			t.Errorf("NoError with error should call Errorf once, got %d calls", len(mock.errorCalls))
		}
	})
}

// TestBehavioralNonCircularVerification ensures our tests aren't circular
func TestBehavioralNonCircularVerification(t *testing.T) {
	t.Run("MockIndependentOfAssertionLibrary", func(t *testing.T) {
		// This test ensures our mock works independently of the assertion library
		mock := &behaviorMockT{}

		// Direct mock testing - this should work regardless of assertion implementation
		mock.Errorf("test message %d", 42)
		mock.Helper()
		mock.FailNow()

		// Verify mock behavior directly
		if len(mock.errorCalls) != 1 {
			t.Errorf("Mock should capture Errorf calls, got %d", len(mock.errorCalls))
		}
		if mock.errorCalls[0] != "test message 42" {
			t.Errorf("Mock should format messages correctly, got %q", mock.errorCalls[0])
		}
		if mock.helperCalls != 1 {
			t.Errorf("Mock should count Helper calls, got %d", mock.helperCalls)
		}
		if mock.failNowCalls != 1 {
			t.Errorf("Mock should count FailNow calls, got %d", mock.failNowCalls)
		}
	})

	t.Run("AssertionLibraryCallsTestingTInterface", func(t *testing.T) {
		// This test verifies the assertion library actually implements the TestingT contract
		mock := &behaviorMockT{}
		assert := New(mock)

		// Execute operations that should call TestingT methods
		assert.Equal(1, 2) // Should call Errorf

		// Verify TestingT interface was actually called
		if len(mock.errorCalls) == 0 {
			t.Fatal("Assertion library MUST call TestingT.Errorf on failures - this is a contract violation")
		}

		// Verify Helper was called for better stack traces
		if mock.helperCalls == 0 {
			t.Error("Assertion library should call TestingT.Helper for better stack traces")
		}
	})
}

// containsSubstring is a helper to check substring without using the assertion library being tested
func containsSubstring(s, substr string) bool {
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
