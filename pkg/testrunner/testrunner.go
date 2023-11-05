package testrunner

import (
	"fmt"
	"testing"
)


// TestRunner represents a basic test runner component.
type TestRunner struct {
	testName string
	t	*testing.T
}

// NewTestRunner creates a new TestRunner with the given test name.
func NewTestRunner(testName string, t *testing.T) *TestRunner {
	return &TestRunner{
		testName: testName,
		t:	t,
	}
}

// RunTest executes a test with the specified test name.
func (tr *TestRunner) RunTest(testFunc func(t *testing.T, assert Assertions)) bool {
	//Create an instance of Assertions.
	assert := New(tr.t)

	// Call the provided test function
	testFunc(tr.t, assert)

	// Determin the result based on the error message.
	result := assert.Error() == ""

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
