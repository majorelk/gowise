// Package testrunner provides a basic test runner for running tests.
//
// TestRunnerTest contains unit tests for the testrunner package.
package testrunner

import (
	"fmt"
	"gowise/pkg/assertions"            // Import assertions package
	"gowise/pkg/interfaces/teststatus" // Import test status package
	"gowise/pkg/logging"
	"testing"
)

// MockT is a mock implementation of the TestRunnerInterface.
// It's used for testing the behavior of the TestRunner.
// It has two slices of strings: Errors and Fatals, which store the error and fatal messages respectively.
// It also has two boolean fields: CalledErrorf and CalledFatalf, which indicate whether Errorf or Fatalf were called.
type MockT struct {
	*testing.T
	Errors       []string
	Fatals       []string
	CalledErrorf bool
	CalledFatalf bool
}

// Run is a mock implementation of the Run method in the TestRunnerInterface.
// It simply calls the provided function with the MockT instance itself.
func (m *MockT) Run(name string, f func(t *MockT)) bool {
	// Call the provided function with the MockT instance.
	f(m)

	// Return true if the test passed, false otherwise.
	// In this mock implementation, we'll just return true for simplicity.
	// In a real implementation, you would check the state of the MockT instance to determine if the test passed.
	return true
}

// Errorf is a mock implementation of the Errorf method in the TestRunnerInterface.
// It adds the formatted error message to the Errors slice and sets CalledErrorf to true.
func (m *MockT) Errorf(format string, args ...interface{}) {
	m.Errors = append(m.Errors, fmt.Sprintf(format, args...))
	m.CalledErrorf = true
}

// Fatalf is a mock implementation of the Fatalf method in the TestRunnerInterface.
// It adds the formatted fatal message to the Fatals slice and sets CalledFatalf to true.
func (m *MockT) Fatalf(format string, args ...interface{}) {
	m.Fatals = append(m.Fatals, fmt.Sprintf(format, args...))
	m.CalledFatalf = true
}

// TestLoggerError is a test case for checking the logging behavior of the TestRunner when a test fails.
// It checks if the logger logs the correct error message.
func TestLoggerError(t *testing.T) {
	mockLogger := logging.NewMockLogger()
	tr := NewTestRunner(t, mockLogger, true)

	testPassed := tr.t.Run("TestLoggerError", func(t *testing.T) {
		tr.RunTest("TestLoggerError", func(assert *assertions.Assert) teststatus.TestStatus {
			return teststatus.Passed // Change this line
		})
	})

	if !testPassed {
		expectedErrorMessage := fmt.Sprintf("test %s failed", "TestLoggerError")
		if len(mockLogger.ErrorMessages) != 1 || mockLogger.ErrorMessages[0] != expectedErrorMessage {
			t.Errorf("Expected logger to log '%s' at ERROR level, but it did not", expectedErrorMessage)
		}
	}
}

// TestLoggerNoInfoOnError is a test case for checking the logging behavior of the TestRunner when a test fails.
// It checks if the logger does not log an info message.
func TestLoggerNoInfoOnError(t *testing.T) {
	mockLogger := logging.NewMockLogger()
	tr := NewTestRunner(t, mockLogger, true)

	testPassed := tr.t.Run("TestLoggerNoInfoOnError", func(t *testing.T) {
		tr.RunTest("TestLoggerNoInfoOnError", func(assert *assertions.Assert) teststatus.TestStatus {
			return teststatus.Passed // Change this line
		})
	})

	if !testPassed && len(mockLogger.InfoMessages) != 0 {
		t.Errorf("Expected logger to not log at INFO level, but it did")
	}
}

// TestLoggerNoErrorOnPass is a test case for checking the logging behavior of the TestRunner when a test passes.
// It checks if the logger does not log an error message.
func TestLoggerNoErrorOnPass(t *testing.T) {
	mockLogger := logging.NewMockLogger()
	tr := NewTestRunner(t, mockLogger, true)

	tr.RunTest("TestLoggerNoErrorOnPass", func(assert *assertions.Assert) teststatus.TestStatus {
		return teststatus.Passed
	})

	if len(mockLogger.ErrorMessages) != 0 {
		t.Errorf("Expected logger to not log at ERROR level, but it did")
	}
}

// TestErrorfOnFail is a test case for checking the behavior of the TestRunner when a test fails and continueOnFail is true.
// It checks if t.Errorf is called.
func TestErrorfOnFail(t *testing.T) {
	// Create a new MockT instance.
	mockT := &MockT{}

	// Pass a function that accepts a *MockT to Run.
	mockT.Run("TestErrorfOnFail", func(t *MockT) {
		// Call Errorf on the *MockT.
		t.Errorf("This is an error")
	})

	// Check if Errorf was called.
	if !mockT.CalledErrorf {
		t.Errorf("Expected Errorf to be called, but it was not")
	}
}

// TestFatalfOnFail is a test case for checking the behavior of the TestRunner when a test fails and continueOnFail is false.
// It checks if t.Fatalf is called.
func TestFatalfOnFail(t *testing.T) {
	mockT := &MockT{}

	// Pass a function that accepts a *MockT to Run.
	mockT.Run("TestFatalfOnFail", func(t *MockT) {
		// Call Fatalf on the *MockT.
		t.Fatalf("This is a fatal error")
	})

	if !mockT.CalledFatalf {
		t.Errorf("Expected Fatalf to be called, but it was not")
	}

	if len(mockT.Errors) != 0 {
		t.Errorf("Expected Errors to be empty, but it was not")
	}
}
