// Package testrunner provides a basic test runner for running tests.
//
// TestRunner is a utility package that offers a basic structure for running
// tests and reporting results.
package testrunner

import (
	"fmt"
	"gowise/pkg/assertions"            // Import assertions package
	"gowise/pkg/interfaces/teststatus" // Import test status package
	"gowise/pkg/logging"               // Import logging package"
	"testing"
)

// TestRunnerInterface is an interface that includes the methods needed for running tests.
type TestRunnerInterface interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Run(name string, f func(t *testing.T)) bool
}

// TestRunner represents a basic test runner component.
type TestRunner struct {
	t              TestRunnerInterface
	logger         logging.LoggerInterface
	continueOnFail bool
}

// NewTestRunner creates a new TestRunner with the given TestRunnerInterface, LoggerInterface, and a boolean indicating whether to continue on fail.
// If continueOnFail is true, the TestRunner will continue executing the remaining tests even if a test fails.
// If continueOnFail is false, the TestRunner will stop executing the remaining tests as soon as a test fails.
func NewTestRunner(t TestRunnerInterface, logger logging.LoggerInterface, continueOnFail bool) *TestRunner {
	return &TestRunner{
		t:              t,
		logger:         logger,
		continueOnFail: continueOnFail,
	}
}

// RunTest executes a test with the specified test name and test function.
// The test function is a function that takes an assertions.Assert and returns a teststatus.TestStatus.
// RunTest logs a message indicating whether the test passed or failed.
// If a test fails and continueOnFail is false, RunTest stops the test immediately.
// If a test fails and continueOnFail is true, RunTest logs the failure and continues with the next test.
// It returns the test result.
func (tr *TestRunner) RunTest(testName string, testFunc func(assert *assertions.Assert) teststatus.TestStatus) teststatus.TestStatus {
	var result teststatus.TestStatus

	tr.t.Run(testName, func(t *testing.T) {
		assert := assertions.New(t)
		result = testFunc(assert)

		resultString := result.GetResult()

		if resultString != teststatus.Passed.GetResult() {
			tr.logger.LogError(fmt.Errorf("test %s failed", testName))
			if tr.continueOnFail {
				t.Errorf("Test %s failed", testName)
			} else {
				t.Fatalf("Test %s failed", testName)
			}
		} else {
			tr.logger.LogInfo(fmt.Sprintf("Test %s passed", testName))
			t.Logf("Test %s passed", testName)
		}
	})

	return result
}
