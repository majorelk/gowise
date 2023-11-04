package testrunner

import (
	"fmt"
)


// TestRunner represents a basic test runner component.
type TestRunner struct {
	testName string
}

// NewTestRunner creates a new TestRunner with the given test name.
func NewTestRunner(testName string) *TestRunner {
	return &TestRunner{
		testName: testName,
	}
}

// RunTest executes a test with the specified test name.
func (tr *TestRunner) RunTest() bool {
	result := true

	// Print message to indicate test execution.
	fmt.Printf("Running test %s - Status: %s\n", tr.testName, tr.getStatus(result))

	return result
}

func (tr *TestRunner) getStatus(result bool) string {
	if result {
		return "Passed"
	}
	return "Failed"
}
