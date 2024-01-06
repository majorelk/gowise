// Package testrunner provides a basic test runner for running tests.
//
// TestRunner is a utility package that offers a basic structure for running
// tests and reporting results.
package testrunner

import (
	"gowise/pkg/assertions"            // Import assertions package
	"gowise/pkg/interfaces/teststatus" // Import test status package
	"testing"
)

// TestRunner represents a basic test runner component.
type TestRunner struct {
	t *testing.T
}

// NewTestRunner creates a new TestRunner with the given testing.T.
func NewTestRunner(t *testing.T) *TestRunner {
	return &TestRunner{
		t: t,
	}
}

// RunTest executes a test with the specified test name and test function.
// The test function is a function that takes an assertions.Assert and returns a teststatus.TestStatus.
// RunTest logs a message indicating whether the test passed or failed.
// It returns the test result.
func (tr *TestRunner) RunTest(testName string, testFunc func(assert *assertions.Assert) teststatus.TestStatus) teststatus.TestStatus {
	var result teststatus.TestStatus

	tr.t.Run(testName, func(t *testing.T) {
		assert := assertions.New(t)
		result = testFunc(assert)

		resultString := result.GetResult()

		if resultString != teststatus.Passed.GetResult() {
			t.Errorf("Test %s failed", testName)
		} else {
			t.Logf("Test %s passed", testName)
		}

		// Fail the testing.T if the test function's result is teststatus.Failed
		if result == teststatus.Failed {
			t.Fail()
		}
	})

	return result
}
