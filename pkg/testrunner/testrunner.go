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
func (tr *TestRunner) RunTest() {
	fmt.Printf("Running test %s\n", tr.testName)
	// add execution logic here.
}

// Example function
func Example() {
	tr := NewTestRunner("SampleTest")
	tr.RunTest()
}
