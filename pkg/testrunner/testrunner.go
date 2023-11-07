package testrunner

import (
	"fmt"
	"testing"
	"gowise/pkg/assertions"
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
func (tr *TestRunner) RunTest(testName string, testFunc func(assert *assertions.Assert) error) {
	assert := assertions.New(tr.t)
	if err := testFunc(assert); err != nil {
		tr.t.Errorf("Test %s failed: %v", testName, err)
	} else {
		fmt.Printf("Test %s passed\n", testName)
	}
}

func (tr *TestRunner) getStatus(result bool) string {
	if result {
		return "Passed"
	}
	return "Failed"
}

