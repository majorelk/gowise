package assertions

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// behaviorMockT is defined in assertions_passing_test.go - shared across test files

// silentT is defined in assertions_passing_test.go - shared across test files

// TestErrorAssertionsIntegration tests error assertion behaviour through the public API
// following TDD principles and testing observable behaviour, not implementation details.
func TestErrorAssertionsIntegration(t *testing.T) {
	tests := []struct {
		name                string
		setupAndAssert      func(assert *Assert)
		expectErrorContains []string
		shouldPass          bool
	}{
		{
			name: "NoError - passes when no error",
			setupAndAssert: func(assert *Assert) {
				var err error // nil error
				assert.NoError(err)
			},
			shouldPass: true,
		},
		{
			name: "NoError - fails when error present",
			setupAndAssert: func(assert *Assert) {
				err := errors.New("something went wrong")
				assert.NoError(err)
			},
			expectErrorContains: []string{
				"expected no error",
				"something went wrong",
			},
			shouldPass: false,
		},
		{
			name: "HasError - fails when no error present",
			setupAndAssert: func(assert *Assert) {
				var err error // nil error
				assert.HasError(err)
			},
			expectErrorContains: []string{
				"expected an error but got none",
			},
			shouldPass: false,
		},
		{
			name: "HasError - passes when error present",
			setupAndAssert: func(assert *Assert) {
				err := errors.New("expected error")
				assert.HasError(err)
			},
			shouldPass: true,
		},
		{
			name: "ErrorIs - passes when error matches",
			setupAndAssert: func(assert *Assert) {
				baseErr := errors.New("base error")
				wrappedErr := fmt.Errorf("wrapped: %w", baseErr)
				assert.ErrorIs(wrappedErr, baseErr)
			},
			shouldPass: true,
		},
		{
			name: "ErrorIs - fails when error doesn't match",
			setupAndAssert: func(assert *Assert) {
				err1 := errors.New("first error")
				err2 := errors.New("second error")
				assert.ErrorIs(err1, err2)
			},
			expectErrorContains: []string{
				"expected error to match target",
				"first error",
				"second error",
			},
			shouldPass: false,
		},
		{
			name: "ErrorAs - passes when error can be cast",
			setupAndAssert: func(assert *Assert) {
				customErr := &CustomError{msg: "custom error"}
				wrappedErr := fmt.Errorf("wrapped: %w", customErr)
				var target *CustomError
				assert.ErrorAs(wrappedErr, &target)
			},
			shouldPass: true,
		},
		{
			name: "ErrorAs - fails when error cannot be cast",
			setupAndAssert: func(assert *Assert) {
				regularErr := errors.New("regular error")
				var target *CustomError
				assert.ErrorAs(regularErr, &target)
			},
			expectErrorContains: []string{
				"expected error to be assignable to target type",
				"regular error",
			},
			shouldPass: false,
		},
		{
			name: "ErrorContains - passes when error message contains text",
			setupAndAssert: func(assert *Assert) {
				err := errors.New("file not found in directory")
				assert.ErrorContains(err, "not found")
			},
			shouldPass: true,
		},
		{
			name: "ErrorContains - fails when error message doesn't contain text",
			setupAndAssert: func(assert *Assert) {
				err := errors.New("permission denied")
				assert.ErrorContains(err, "not found")
			},
			expectErrorContains: []string{
				"expected error message to contain",
				"not found",
				"permission denied",
			},
			shouldPass: false,
		},
		{
			name: "ErrorContains - fails when error is nil",
			setupAndAssert: func(assert *Assert) {
				var err error // nil
				assert.ErrorContains(err, "some text")
			},
			expectErrorContains: []string{
				"expected error but got nil",
				"some text",
			},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a dummy testing.T that captures failures
			mock := &behaviorMockT{}
			assert := New(mock)

			// Execute the assertion
			tt.setupAndAssert(assert)

			// Check behavioral contract through TestingT interface calls
			if tt.shouldPass {
				// Framework behavior: PASS = no Errorf calls
				if len(mock.errorCalls) != 0 {
					t.Errorf("Expected assertion to pass (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
				}
			} else {
				// Framework behavior: FAIL = exactly 1 Errorf call with expected content
				if len(mock.errorCalls) != 1 {
					t.Fatalf("Expected assertion to fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
				}

				// Check that all expected parts are in the error message
				errorMsg := mock.errorCalls[0]
				for _, expected := range tt.expectErrorContains {
					if !strings.Contains(errorMsg, expected) {
						t.Errorf("Error message missing expected content %q\nFull error message:\n%s", expected, errorMsg)
					}
				}
			}
		})
	}
}

// TestAdvancedErrorAssertionsIntegration tests more sophisticated error handling
func TestAdvancedErrorAssertionsIntegration(t *testing.T) {
	tests := []struct {
		name                string
		setupAndAssert      func(assert *Assert)
		expectErrorContains []string
		shouldPass          bool
	}{
		{
			name: "ErrorMatches - passes when error message matches regex",
			setupAndAssert: func(assert *Assert) {
				err := errors.New("file operation failed: permission denied (code: 403)")
				assert.ErrorMatches(err, `permission denied \(code: \d+\)`)
			},
			shouldPass: true,
		},
		{
			name: "ErrorMatches - fails when error message doesn't match regex",
			setupAndAssert: func(assert *Assert) {
				err := errors.New("simple error message")
				assert.ErrorMatches(err, `permission denied \(code: \d+\)`)
			},
			expectErrorContains: []string{
				"expected error message to match pattern",
				`permission denied \(code: \d+\)`,
				"simple error message",
			},
			shouldPass: false,
		},
		{
			name: "ErrorType - passes when error is correct type",
			setupAndAssert: func(assert *Assert) {
				customErr := &CustomError{msg: "custom error"}
				assert.ErrorType(customErr, &CustomError{})
			},
			shouldPass: true,
		},
		{
			name: "ErrorType - fails when error is wrong type",
			setupAndAssert: func(assert *Assert) {
				regularErr := errors.New("regular error")
				assert.ErrorType(regularErr, &CustomError{})
			},
			expectErrorContains: []string{
				"expected error of different type",
				"*errors.errorString",
				"*assertions.CustomError",
			},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &behaviorMockT{}
			assert := New(mock)

			tt.setupAndAssert(assert)

			// Check behavioral contract through TestingT interface calls
			if tt.shouldPass {
				// Framework behavior: PASS = no Errorf calls
				if len(mock.errorCalls) != 0 {
					t.Errorf("Expected assertion to pass (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
				}
			} else {
				// Framework behavior: FAIL = exactly 1 Errorf call with expected content
				if len(mock.errorCalls) != 1 {
					t.Fatalf("Expected assertion to fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
				}

				errorMsg := mock.errorCalls[0]
				for _, expected := range tt.expectErrorContains {
					if !strings.Contains(errorMsg, expected) {
						t.Errorf("Error message missing expected content %q\nFull error message:\n%s", expected, errorMsg)
					}
				}
			}
		})
	}
}

// CustomError is a test helper for error type testing
type CustomError struct {
	msg string
}

func (e *CustomError) Error() string {
	return e.msg
}

// Examples for documentation - demonstrating proper usage of error assertions

func ExampleAssert_NoError() {
	assert := New(&silentT{})

	// Test that no error occurred
	var err error // nil error
	assert.NoError(err)

	// The assertion succeeded since the condition was met
	fmt.Println("No error:", true)
	// Output: No error: true
}

func ExampleAssert_HasError() {
	assert := New(&silentT{})

	// Test that an error occurred
	err := errors.New("something went wrong")
	assert.HasError(err)

	// The assertion succeeded since the condition was met
	fmt.Println("No error:", true)
	// Output: No error: true
}

func ExampleAssert_ErrorIs() {
	assert := New(&silentT{})

	// Test error unwrapping with errors.Is
	baseErr := errors.New("base error")
	wrappedErr := fmt.Errorf("wrapped: %w", baseErr)
	assert.ErrorIs(wrappedErr, baseErr)

	// The assertion succeeded since the condition was met
	fmt.Println("No error:", true)
	// Output: No error: true
}

func ExampleAssert_ErrorAs() {
	assert := New(&silentT{})

	// Test error type assertion with errors.As
	customErr := &CustomError{msg: "custom error"}
	wrappedErr := fmt.Errorf("wrapped: %w", customErr)
	var target *CustomError
	assert.ErrorAs(wrappedErr, &target)

	// The assertion succeeded since the condition was met
	fmt.Println("No error:", true)
	// Output: No error: true
}

func ExampleAssert_ErrorContains() {
	assert := New(&silentT{})

	// Test that error message contains specific text
	err := errors.New("file not found in directory")
	assert.ErrorContains(err, "not found")

	// The assertion succeeded since the condition was met
	fmt.Println("No error:", true)
	// Output: No error: true
}

func ExampleAssert_ErrorMatches() {
	assert := New(&silentT{})

	// Test that error message matches regex pattern
	err := errors.New("operation failed: timeout after 30 seconds")
	assert.ErrorMatches(err, `timeout after \d+ seconds`)

	// The assertion succeeded since the condition was met
	fmt.Println("No error:", true)
	// Output: No error: true
}

func ExampleAssert_Panics() {
	assert := New(&silentT{})

	// Test that a function panics
	assert.Panics(func() {
		panic("something went wrong")
	})

	// The assertion succeeded since the condition was met
	fmt.Println("No error:", true)
	// Output: No error: true
}

func ExampleAssert_NotPanics() {
	assert := New(&silentT{})

	// Test that a function does not panic
	assert.NotPanics(func() {
		// Normal function that doesn't panic
		_ = 1 + 1
	})

	// The assertion succeeded since the condition was met
	fmt.Println("No error:", true)
	// Output: No error: true
}

func ExampleAssert_PanicsWith() {
	assert := New(&silentT{})

	// Test that a function panics with a specific value
	expectedPanic := "specific error"
	assert.PanicsWith(func() {
		panic(expectedPanic)
	}, expectedPanic)

	// The assertion succeeded since the condition was met
	fmt.Println("No error:", true)
	// Output: No error: true
}

// TestPanicAssertionsIntegration tests panic assertion behaviour through the public API
// following TDD principles and testing observable behaviour, not implementation details.
func TestPanicAssertionsIntegration(t *testing.T) {
	tests := []struct {
		name                string
		setupAndAssert      func(assert *Assert)
		expectErrorContains []string
		shouldPass          bool
	}{
		{
			name: "Panics - passes when function panics",
			setupAndAssert: func(assert *Assert) {
				assert.Panics(func() {
					panic("test panic")
				})
			},
			shouldPass: true,
		},
		{
			name: "Panics - fails when function doesn't panic",
			setupAndAssert: func(assert *Assert) {
				assert.Panics(func() {
					// Function that doesn't panic
					_ = 1 + 1
				})
			},
			expectErrorContains: []string{
				"expected to panic",
			},
			shouldPass: false,
		},
		{
			name: "NotPanics - passes when function doesn't panic",
			setupAndAssert: func(assert *Assert) {
				assert.NotPanics(func() {
					// Normal function execution
					result := "safe operation"
					_ = result
				})
			},
			shouldPass: true,
		},
		{
			name: "NotPanics - fails when function panics",
			setupAndAssert: func(assert *Assert) {
				assert.NotPanics(func() {
					panic("unexpected panic")
				})
			},
			expectErrorContains: []string{
				"expected not to panic",
				"unexpected panic",
			},
			shouldPass: false,
		},
		{
			name: "PanicsWith - passes when function panics with expected value",
			setupAndAssert: func(assert *Assert) {
				expectedValue := "specific panic message"
				assert.PanicsWith(func() {
					panic(expectedValue)
				}, expectedValue)
			},
			shouldPass: true,
		},
		{
			name: "PanicsWith - fails when function panics with different value",
			setupAndAssert: func(assert *Assert) {
				assert.PanicsWith(func() {
					panic("actual panic")
				}, "expected panic")
			},
			expectErrorContains: []string{
				"expected to panic with",
				"expected panic",
				"actual panic",
			},
			shouldPass: false,
		},
		{
			name: "PanicsWith - fails when function doesn't panic",
			setupAndAssert: func(assert *Assert) {
				assert.PanicsWith(func() {
					// Function that doesn't panic
				}, "expected panic")
			},
			expectErrorContains: []string{
				"expected to panic with",
				"expected panic",
			},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a dummy testing.T that captures failures
			mock := &behaviorMockT{}
			assert := New(mock)

			// Execute the assertion
			tt.setupAndAssert(assert)

			// Check behavioral contract through TestingT interface calls
			if tt.shouldPass {
				// Framework behavior: PASS = no Errorf calls
				if len(mock.errorCalls) != 0 {
					t.Errorf("Expected assertion to pass (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
				}
			} else {
				// Framework behavior: FAIL = exactly 1 Errorf call with expected content
				if len(mock.errorCalls) != 1 {
					t.Fatalf("Expected assertion to fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
				}

				errorMsg := mock.errorCalls[0]
				// Check that all expected parts are in the error message
				for _, expected := range tt.expectErrorContains {
					if !strings.Contains(errorMsg, expected) {
						t.Errorf("Error message missing expected content %q\\nFull error message:\\n%s", expected, errorMsg)
					}
				}
			}
		})
	}
}
