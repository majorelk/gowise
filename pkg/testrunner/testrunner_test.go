// Package testrunner provides a basic test runner for running tests.
//
// TestRunnerTest contains unit tests for the testrunner package.
package testrunner

import (
	"testing"
	"gowise/pkg/assertions" // Import the assertions package
)

// TestCreateTestRunner is an example test case for creating a TestRunner and using assertions in tests.
func TestCreateTestRunner(t *testing.T) {
	tr := NewTestRunner(t)

	tr.RunTest("TestCreateTestRunner", func(assert *assertions.Assert) error { 
		// Use assertions.Assert type
		// Testing if 2 + 2 equals 4
		assert.True(2+2 == 4)
		return nil
	})
}

