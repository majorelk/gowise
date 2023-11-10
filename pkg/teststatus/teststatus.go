// Package teststatus provides an enumeration for the result of running a test
package teststatus

// TestStatus represents the result of running a test
type TestStatus interface {
	GetResult() string
}


// Result represents a test result
type Result string

const (
	// Inconclusive indicates that the test was inconclusive
	Inconclusive Result = "Inconclusive"

	// Skipped indicates that the test has been skipped
	Skipped Result = "Skipped"

	// Passed indicates that the test succeeded
	Passed Result = "Passed"

	// Warning indicates that there was a warning
	Warning Result = "Warning"

	// Failed indicates that the test failed
	Failed Result = "Failed"
)

