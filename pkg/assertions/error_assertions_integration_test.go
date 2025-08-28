package assertions

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

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
			dummyT := &capturingT{}
			assert := New(dummyT)

			// Execute the assertion
			tt.setupAndAssert(assert)

			errorMsg := assert.Error()

			if tt.shouldPass {
				// Should pass - no error message
				if errorMsg != "" {
					t.Errorf("Expected assertion to pass but got error: %s", errorMsg)
				}
			} else {
				// Should fail - check error message contains expected content
				if errorMsg == "" {
					t.Fatalf("Expected error message but got none")
				}

				// Check that all expected parts are in the error message
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
			dummyT := &capturingT{}
			assert := New(dummyT)

			tt.setupAndAssert(assert)

			errorMsg := assert.Error()

			if tt.shouldPass {
				if errorMsg != "" {
					t.Errorf("Expected assertion to pass but got error: %s", errorMsg)
				}
			} else {
				if errorMsg == "" {
					t.Fatalf("Expected error message but got none")
				}

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
	assert := New(&testing.T{})

	// Test that no error occurred
	var err error // nil error
	assert.NoError(err)

	fmt.Println("No error:", assert.Error() == "")
	// Output: No error: true
}

func ExampleAssert_HasError() {
	assert := New(&testing.T{})

	// Test that an error occurred
	err := errors.New("something went wrong")
	assert.HasError(err)

	fmt.Println("No error:", assert.Error() == "")
	// Output: No error: true
}

func ExampleAssert_ErrorIs() {
	assert := New(&testing.T{})

	// Test error unwrapping with errors.Is
	baseErr := errors.New("base error")
	wrappedErr := fmt.Errorf("wrapped: %w", baseErr)
	assert.ErrorIs(wrappedErr, baseErr)

	fmt.Println("No error:", assert.Error() == "")
	// Output: No error: true
}

func ExampleAssert_ErrorAs() {
	assert := New(&testing.T{})

	// Test error type assertion with errors.As
	customErr := &CustomError{msg: "custom error"}
	wrappedErr := fmt.Errorf("wrapped: %w", customErr)
	var target *CustomError
	assert.ErrorAs(wrappedErr, &target)

	fmt.Println("No error:", assert.Error() == "")
	// Output: No error: true
}

func ExampleAssert_ErrorContains() {
	assert := New(&testing.T{})

	// Test that error message contains specific text
	err := errors.New("file not found in directory")
	assert.ErrorContains(err, "not found")

	fmt.Println("No error:", assert.Error() == "")
	// Output: No error: true
}

func ExampleAssert_ErrorMatches() {
	assert := New(&testing.T{})

	// Test that error message matches regex pattern
	err := errors.New("operation failed: timeout after 30 seconds")
	assert.ErrorMatches(err, `timeout after \d+ seconds`)

	fmt.Println("No error:", assert.Error() == "")
	// Output: No error: true
}
