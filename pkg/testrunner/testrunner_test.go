// Package testrunner provides a basic test runner for running tests.
//
// TestRunnerTest contains unit tests for the testrunner package.
package testrunner

import (
	"gowise/pkg/assertions"
	"gowise/pkg/interfaces/teststatus"
	"testing"
)

// TestNewTestRunner is a test case for creating a TestRunner.
// It checks if the NewTestRunner function returns a non-nil value.
func TestNewTestRunner(t *testing.T) {
	tr := NewTestRunner(t)
	if tr == nil {
		t.Errorf("NewTestRunner was incorrect, got: nil")
	}
}

// TestRunTest is a test case for running a test with the TestRunner.
// It checks if the RunTest function correctly runs a test function and returns the correct test status.
// The test function used in TestRunTest is a simple function that checks if 2+2 equals 4,
// and returns teststatus.Failed if it does not, and teststatus.Passed if it does.
func TestRunTest(t *testing.T) {
	tr := NewTestRunner(t)

	result := tr.RunTest("TestRunTest", func(assert *assertions.Assert) teststatus.TestStatus {
		if 2+2 != 4 {
			assert.True(false) // This will record a failure
			return teststatus.Failed
		}

		return teststatus.Passed
	})

	if result != teststatus.Passed {
		t.Errorf("Expected test to pass, but it failed")
	}
}
