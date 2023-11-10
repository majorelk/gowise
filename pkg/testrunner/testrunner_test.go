// Package testrunner provides a basic test runner for running tests.
//
// TestRunnerTest contains unit tests for the testrunner package.
package testrunner

import (
	"testing"
	"gowise/pkg/assertions"    // Import the assertions package
	"gowise/pkg/interfaces/teststatus"    // Import the teststatus package
)

// TestCreateTestRunner is an example test case for creating a TestRunner and using assertions in tests.
func TestCreateTestRunner(t *testing.T) {
	tr := NewTestRunner(t)

	tr.RunTest("TestCreateTestRunner", func(assert *assertions.Assert) teststatus.TestStatus {
		// Use assertions.Assert type
		// Testing if 2 + 2 equals 4
		assert.True(2+2 == 4)

		// Check the assertion result
		if assert.Error() != "" {
			// Return the test result using teststatus.Failed if there is an error
			return teststatus.Result(teststatus.Failed)
		}

		// Return the test result using teststatus.Passed if the assertion passes
		return teststatus.Result(teststatus.Passed)
	})
}

