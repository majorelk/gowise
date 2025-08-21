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

	"crypto/rand"
	"encoding/hex"
)

// TestInterface is an interface that includes the methods needed for running tests.
// Errorf is used to indicate a non-fatal error that allows the test to continue.
// Fatalf is used to indicate a fatal error that stops the test immediately.
// Run executes a subtest with the given name and function.
type TestInterface interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Run(name string, f func(t TestInterface)) bool
}

// TestRunner represents a basic test runner component.
// t is the interface for running tests.
// logger is used for logging test results.
// continueOnFail determines whether the TestRunner should continue executing the remaining tests if a test fails.
// results stores the results of all executed tests.
// reporter is used for reporting test results.
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
// testName is the name of the test.
// testFunc is a function that takes an assertions.Assert and returns a teststatus.TestStatus. It contains the logic of the test.
// The method logs a message indicating whether the test passed or failed.
// If a test fails and continueOnFail is false, the method stops the test immediately.
// If a test fails and continueOnFail is true, the method logs the failure and continues with the next test.
// The method returns the result of the test.
func (tr *TestRunner) RunTest(testName string, testFunc func(assert *assertions.Assert) teststatus.TestStatus) teststatus.TestStatus {
	var resultOutside teststatus.TestStatus

	tr.t.Run(testName, func(t TestInterface) {
		startTime := time.Now() // Get the current time

		assert := assertions.New(t)
		resultInside := testFunc(assert) // Declare a new result variable here

		endTime := time.Now()              // Get the current time
		duration := endTime.Sub(startTime) // Calculate the duration of the test

		resultString := resultInside.GetResult()
		testID := generateTestID() // Generate a unique ID for the test

		if resultInside != teststatus.Passed {
			tr.logger.LogError(fmt.Errorf("test %s failed", testName))
			var err error
			if tr.continueOnFail {
				t.Errorf("Test %s failed", testName)
			} else {
				t.Fatalf("Test %s failed", testName)
			}
			if err != nil {
				tr.logger.LogError(fmt.Errorf("failed to log error: %v", err))
			}
		} else {
			tr.logger.LogInfo(fmt.Sprintf("Test %s passed", testName))
		}

		// Report the test output
		output := testoutput.NewTestOutput(duration.String(), resultInside.GetResult(), testID, testName, resultString)
		if err := tr.reporter.ReportTestOutput(output); err != nil {
			tr.logger.LogError(fmt.Errorf("failed to report test output: %v", err))
		}

		resultOutside = resultInside // Assign the result from inside the goroutine to the outside variable
	})

	tr.results = append(tr.results, resultOutside)

	return resultOutside
}

// GenerateReport generates a test report.
// The report contains the results of all executed tests.
func (tr *TestRunner) GenerateReport() *reporter.TestReport { // Fix the undeclared name error by using the imported package
	report := reporter.NewTestReport() // Use the imported package to create a new TestReport

	for _, result := range tr.results {
		report.AddResult(result)
	}

	return report
}

// generateTestID creates a unique test ID using crypto/rand (stdlib only)
func generateTestID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
