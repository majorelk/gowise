// Package testrunner provides a basic test runner for running tests.
//
// TestRunner is a utility package that offers a basic structure for running
// tests and reporting results.
package testrunner

import (
	"fmt"
	"testing"
	"gowise/pkg/assertions" // Import assertions package
	"gowise/pkg/interfaces/teststatus" // Import assertions package
)

// TestRunner represents a basic test runner component.
type TestRunner struct {
	t *testing.T
}

// NewTestRunner creates a new TestRunner with the given testing.T.
func NewTestRunner(t *testing.T) *TestRunner {
	return &TestRunner{
		t:	t,
	}
}

// RunTest executes a test with the specified test name.
func (tr *TestRunner) RunTest(testName string, testFunc func(assert *assertions.Assert) teststatus.TestStatus) {
	assert := assertions.New(tr.t)
	result :- testFunc(assert)

	// Use the GetResult method to get the result as a string
	resultString := result.GetResult()

	if resultString == teststatus.Passed.GetResult() {
		fmt.Printf("Test %s passed\n", testName)
	} else {
		tr.t.Errorf("Test %s failed: %v", testName, resultString)
	}
}

