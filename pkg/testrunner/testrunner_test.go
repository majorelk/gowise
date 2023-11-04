package testrunner

import "testing"

func TestCreateTestRunner(t *testing.T) {
	// Step 1: Create a TestRunner instance with a test name.
	testName:= "SampleTest"
	tr := NewTestRunner(testName)

	// Step 2: Check the properties of the TestRunner.
	if tr.testName != testName {
		t.Errorf("Expected test name to be %s, but got %s", testName, tr.testName)
	}

	// add more property checks here if required.
}
