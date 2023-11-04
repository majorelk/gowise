package testrunner

import "testing"

func TestCreateTestRunner(t *testing.T) {
	// Create a TestRunner instance with a test name.
	tr := NewTestRunner("SampleTest")

	// Run the test using the TestRunner and get the test results.
	result := tr.RunTest()

	// Define the expected behaviour of the test.
	expectedResult := true

	// check the test result using the testing framework.
	if result == expectedResult {
		t.Logf("Test %s: Passed", tr.testName)
	} else {
		t.Errorf("Test %s: Expected result was %v, but got %v", tr.testName, expectedResult, result)
	}
}
