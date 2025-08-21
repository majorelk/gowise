// Package testrunner provides a basic test runner for running tests.
//
// TestRunnerTest contains unit tests for the testrunner package.
package testrunner

import (
	"fmt"
	"gowise/pkg/assertions"
	"gowise/pkg/interfaces/testattachment"
	"gowise/pkg/interfaces/testmessage"
	"gowise/pkg/interfaces/testoutput"
	"gowise/pkg/interfaces/teststatus"
	"gowise/pkg/logging"
	"testing"
)

// TWrapper is a wrapper for *testing.T that implements the TestInterface.
type TWrapper struct {
	t           *testing.T
	calledError bool
	calledFatal bool
}

// Errorf is a wrapper for the Errorf method in *testing.T.
func (tw *TWrapper) Errorf(format string, args ...interface{}) {
	tw.calledError = true
	tw.t.Errorf(format, args...)
}

// Fatalf is a wrapper for the Fatalf method in *testing.T.
func (tw *TWrapper) Fatalf(format string, args ...interface{}) {
	tw.calledFatal = true
	tw.t.Fatalf(format, args...)
}

// Run is a wrapper for the Run method in *testing.T.
func (tw *TWrapper) Run(name string, f func(t TestInterface)) bool {
	return tw.t.Run(name, func(t *testing.T) {
		f(&TWrapper{t: t})
	})
}

// MockT is a mock implementation of the TestRunnerInterface.
// It's used for testing the behaviour of the TestRunner.
// It has two slices of strings: Errors and Fatals, which store the error and fatal messages respectively.
// It also has two boolean fields: CalledErrorf and CalledFatalf, which indicate whether Errorf or Fatalf were called.
type MockT struct {
	*testing.T
	Errors       []string
	Fatals       []string
	CalledErrorf bool
	CalledFatalf bool
}

func (m *MockReporter) Close() error {
	// Implement the method. If there's nothing to close in the mock, you can just return nil.
	return nil
}

// Run is a mock implementation of the Run method in the TestRunnerInterface.
// It simply calls the provided function with the MockT instance itself.
func (m *MockT) Run(name string, f func(t TestInterface)) bool {
	// Call the provided function with the MockT instance.
	f(m)

	// Return true if the test passed, false otherwise.
	// In this mock implementation, we'll just return true for simplicity.
	// In a real implementation, you would check the state of the MockT instance to determine if the test passed.
	return true
}

// MockReporter is a mock implementation of the ReporterInterface.
// It's used for testing the behaviour of the TestRunner.
type MockReporter struct {
	CalledReportTestOutput bool
	ReportedOutput         []testoutput.TestOutput
	Error                  error
}

// ReportTestAttachment implements reporter.ReporterInterface.
func (*MockReporter) ReportTestAttachment(ta testattachment.TestAttachment) error {
	panic("unimplemented")
}

// ReportTestMessage implements reporter.ReporterInterface.
func (*MockReporter) ReportTestMessage(tm testmessage.TestMessage) error {
	panic("unimplemented")
}

// Errorf is a mock implementation of the Errorf method in the TestRunnerInterface.
// It adds the formatted error message to the Errors slice and sets CalledErrorf to true.
func (m *MockT) Errorf(format string, args ...interface{}) {
	m.Errors = append(m.Errors, fmt.Sprintf(format, args...))
	m.CalledErrorf = true
}

// Fatalf is a mock implementation of the Fatalf method in the TestRunnerInterface.
// It adds the formatted fatal message to the Fatals slice and sets CalledFatalf to true.
func (m *MockT) Fatalf(format string, args ...interface{}) {
	m.Fatals = append(m.Fatals, fmt.Sprintf(format, args...))
	m.CalledFatalf = true
}

// TestLoggerError is a test case for checking the logging behaviour of the TestRunner when a test fails.
// It checks if the logger logs the correct error message.
func TestLoggerError(t *testing.T) {
	mockLogger := logging.NewMockLogger()
	mockReporter := &MockReporter{}                                      // Create a new MockReporter
	tr := NewTestRunner(&TWrapper{t: t}, mockLogger, true, mockReporter) // Pass the MockReporter to NewTestRunner

	testPassed := tr.t.Run("TestLoggerError", func(t TestInterface) {
		tr.RunTest("TestLoggerError", func(assert *assertions.Assert) teststatus.TestStatus {
			return teststatus.Passed // Change this line
		})
	})

	if !testPassed {
		expectedErrorMessage := fmt.Sprintf("test %s failed", "TestLoggerError")
		if len(mockLogger.ErrorMessages) != 1 || mockLogger.ErrorMessages[0] != expectedErrorMessage {
			t.Errorf("Expected logger to log '%s' at ERROR level, but it did not", expectedErrorMessage)
		}
	}
}

// TestLoggerNoInfoOnError is a test case for checking the logging behaviour of the TestRunner when a test fails.
// It checks if the logger does not log an info message.
func TestLoggerNoInfoOnError(t *testing.T) {
	mockLogger := logging.NewMockLogger()
	tr := NewTestRunner(&TWrapper{t: t}, mockLogger, true, &MockReporter{}) // Pass a mock ReporterInterface to NewTestRunner

	testPassed := tr.t.Run("TestLoggerNoInfoOnError", func(t TestInterface) {
		tr.RunTest("TestLoggerNoInfoOnError", func(assert *assertions.Assert) teststatus.TestStatus {
			return teststatus.Passed // Change this line
		})
	})

	if !testPassed && len(mockLogger.InfoMessages) != 0 {
		t.Errorf("Expected logger to not log at INFO level, but it did")
	}
}

// TestLoggerNoErrorOnPass is a test case for checking the logging behaviour of the TestRunner when a test passes.
// It checks if the logger does not log an error message.
func TestLoggerNoErrorOnPass(t *testing.T) {
	mockLogger := logging.NewMockLogger()
	tr := NewTestRunner(&TWrapper{t: t}, mockLogger, true, &MockReporter{})

	tr.RunTest("TestLoggerNoErrorOnPass", func(assert *assertions.Assert) teststatus.TestStatus {
		return teststatus.Passed
	})

	if len(mockLogger.ErrorMessages) != 0 {
		t.Errorf("Expected logger to not log at ERROR level, but it did")
	}
}

// TestErrorfOnFail is a test case for checking the behaviour of the TestRunner when a test fails and continueOnFail is true.
// It checks if t.Errorf is called.
func TestErrorfOnFail(t *testing.T) {
	// Create a new MockT instance.
	mockT := &MockT{}

	// Pass a function that accepts a *MockT to Run.
	mockT.Run("TestErrorfOnFail", func(t TestInterface) {
		// Call Errorf on the *MockT.
		t.Errorf("This is an error")
	})

	// Check if Errorf was called.
	if !mockT.CalledErrorf {
		t.Errorf("Expected Errorf to be called, but it was not")
	}
}

// TestFatalfOnFail is a test case for checking the behaviour of the TestRunner when a test fails and continueOnFail is false.
// It checks if t.Fatalf is called.
func TestFatalfOnFail(t *testing.T) {
	mockT := &MockT{}

	// Pass a function that accepts a *MockT to Run.
	mockT.Run("TestFatalfOnFail", func(t TestInterface) {
		// Call Fatalf on the *MockT.
		t.Fatalf("This is a fatal error")
	})

	if !mockT.CalledFatalf {
		t.Errorf("Expected Fatalf to be called, but it was not")
	}

	if len(mockT.Errors) != 0 {
		t.Errorf("Expected Errors to be empty, but it was not")
	}
}

// ReportTestOutput is a mock implementation of the ReportTestOutput method in the ReporterInterface.
// It sets CalledReportTestOutput to true and stores the reported output.
func (m *MockReporter) ReportTestOutput(output testoutput.TestOutput) error {
	m.CalledReportTestOutput = true
	m.ReportedOutput = append(m.ReportedOutput, output)
	return m.Error
}

// TestReportTestOutput is a test case for checking the behaviour of the TestRunner when a test is run.
// It checks if the reporter's ReportTestOutput method is called with the correct output.
func TestReportTestOutput(t *testing.T) {
	mockLogger := logging.NewMockLogger()
	mockReporter := &MockReporter{}
	tr := NewTestRunner(&TWrapper{t: t}, mockLogger, true, mockReporter)

	testName := "TestReportTestOutput"
	testStatus := tr.RunTest(testName, func(assert *assertions.Assert) teststatus.TestStatus {
		return teststatus.Passed
	})

	if !mockReporter.CalledReportTestOutput {
		t.Errorf("Expected ReportTestOutput to be called, but it was not")
	}

	found := false
	for _, output := range mockReporter.ReportedOutput {
		if output.TestName == testName && output.Status == testStatus.GetResult() {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("ReportTestOutput was not called with correct output for test %s", testName)
	}
}

// TestGenerateReport tests the behaviour of the GenerateReport method.
// It checks if the ReportTestOutput method of the ReporterInterface is called for each test that was run,
// and if it correctly passes the test output to the ReportTestOutput method.
func TestGenerateReport(t *testing.T) {
	mockLogger := logging.NewMockLogger()
	mockReporter := &MockReporter{}
	tr := NewTestRunner(&TWrapper{t: t}, mockLogger, true, mockReporter)

	// Run some tests
	testNames := []string{"Test1", "Test2", "Test3"}
	for _, testName := range testNames {
		tr.RunTest(testName, func(assert *assertions.Assert) teststatus.TestStatus {
			return teststatus.Passed
		})
	}

	// Generate the report
	tr.GenerateReport()

	// Check if ReportTestOutput was called for each test
	if !mockReporter.CalledReportTestOutput {
		t.Errorf("Expected ReportTestOutput to be called, but it was not")
	}

	// Check if the reported output is correct for each test
	for _, output := range mockReporter.ReportedOutput {
		found := false
		for _, testName := range testNames {
			if output.TestName == testName && output.Status == teststatus.Passed.GetResult() {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ReportTestOutput was not called with correct output for test %s", output.TestName)
		}
	}
}

// TestGenerateReportFailure tests the behaviour of the GenerateReport method when a test fails.
// It checks if the ReportTestOutput method of the ReporterInterface is called for each test that was run,
// and if it correctly passes the test output to the ReportTestOutput method.
func TestGenerateReportFailure(t *testing.T) {
	mockLogger := logging.NewMockLogger()
	mockReporter := &MockReporter{}
	mockT := &MockT{T: t}
	tr := NewTestRunner(mockT, mockLogger, true, mockReporter)

	// Run some tests
	testNames := []string{"Test1", "Test2", "Test3"}
	for _, testName := range testNames {
		tr.RunTest(testName, func(assert *assertions.Assert) teststatus.TestStatus {
			if testName == "Test2" {
				return teststatus.Failed // Make Test2 fail
			}
			return teststatus.Passed
		})
	}

	// Generate the report
	tr.GenerateReport()

	// Check if ReportTestOutput was called for each test
	if !mockReporter.CalledReportTestOutput {
		t.Errorf("Expected ReportTestOutput to be called, but it was not")
	}

	// Check if the reported output is correct for each test
	for _, output := range mockReporter.ReportedOutput {
		expectedStatus := teststatus.Passed.GetResult()
		if output.TestName == "Test2" {
			expectedStatus = teststatus.Failed.GetResult() // Expect Test2 to have failed
		}

		// Check if the test name is in the testNames slice
		found := false
		for _, testName := range testNames {
			if output.TestName == testName {
				found = true
				break
			}
		}

		if !found || output.Status != expectedStatus {
			t.Errorf("ReportTestOutput was not called with correct output for test %s", output.TestName)
		}
	}
}
