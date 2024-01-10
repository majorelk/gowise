// Package testrunner provides a basic test runner for running tests.
//
// TestRunner is a utility package that offers a basic structure for running
// tests and reporting results.
package testrunner

import (
	"fmt"
	"gowise/pkg/assertions"            // Import assertions package
	"gowise/pkg/interfaces/testoutput" // Import test output package
	"gowise/pkg/interfaces/teststatus" // Import test status package
	"gowise/pkg/logging"               // Import logging package"
	"gowise/pkg/reporter"              // Import reporter package
	"time"                             // Import time package

	"github.com/google/uuid" // Import uuid package
)

// TestInterface is an interface that includes the methods needed for running tests.
type TestInterface interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Run(name string, f func(t TestInterface)) bool
}

/// TestRunner represents a basic test runner component.
type TestRunner struct {
	t              TestInterface
	logger         logging.LoggerInterface
	continueOnFail bool
	results        []teststatus.TestStatus
	reporter       reporter.ReporterInterface // Fix the undeclared name error by using the imported package
}

// NewTestRunner creates a new TestRunner with the given TestInterface, LoggerInterface, a boolean indicating whether to continue on fail, and a ReporterInterface.
// If continueOnFail is true, the TestRunner will continue executing the remaining tests even if a test fails.
// If continueOnFail is false, the TestRunner will stop executing the remaining tests as soon as a test fails.
func NewTestRunner(t TestInterface, logger logging.LoggerInterface, continueOnFail bool, reporter reporter.ReporterInterface) *TestRunner {
	return &TestRunner{
		t:              t,
		logger:         logger,
		continueOnFail: continueOnFail,
		reporter:       reporter,
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

	tr.t.Run(testName, func(t TestInterface) {
		startTime := time.Now() // Get the current time

		assert := assertions.New(t)
		result = testFunc(assert)

		endTime := time.Now()              // Get the current time
		duration := endTime.Sub(startTime) // Calculate the duration of the test

		resultString := result.GetResult()
		testID := uuid.New().String() // Generate a unique ID for the test

		if resultString != teststatus.Passed.GetResult() {
			tr.logger.LogError(fmt.Errorf("test %s failed", testName))
			if tr.continueOnFail {
				t.Errorf("Test %s failed", testName)
			} else {
				t.Fatalf("Test %s failed", testName)
			}
		} else {
			tr.logger.LogInfo(fmt.Sprintf("Test %s passed", testName))
		}

		// Report the test output
		output := testoutput.NewTestOutput(duration.String(), result.GetResult(), testID, testName, resultString)
		if err := tr.reporter.ReportTestOutput(output); err != nil {
			tr.logger.LogError(fmt.Errorf("failed to report test output: %v", err))
		}
	})

	tr.results = append(tr.results, result)

	return result
}

// GenerateReport generates a test report.
func (tr *TestRunner) GenerateReport() *reporter.TestReport { // Fix the undeclared name error by using the imported package
	report := reporter.NewTestReport() // Use the imported package to create a new TestReport

	for _, result := range tr.results {
		report.AddResult(result)
	}

	return report
}
