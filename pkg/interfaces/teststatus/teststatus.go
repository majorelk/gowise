// Package teststatus provides enums and interfaces for test statuses.
package teststatus

// TestStatus is an interface for representing the result of running a test.
type TestStatus interface {
	GetResult() string
}

// Result is an enum-like type representing different test results.
type Result int

const (
	Passed Result = iota
	Failed
)

// GetResult returns the string representation of the test result.
func (r Result) GetResult() string {
	switch r {
	case Passed:
		return "Passed"
	case Failed:
		return "Failed"
	default:
		return "Unknown"
	}
}

