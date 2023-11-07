package assertions

import (
	"fmt"
	"reflect"
)

// Assert is a struct that holds the testing context and error message.
type Assert struct {
	t	interface{}
	errorMsg string
}

// New creates a new Assert instance with the given testing context.
func New(t interface{}) *Assert {
	return &Assert{t: t}
}

// Equal asserts that two values are equal.
func (a *Assert) Equal(expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		a.reportError(expected, actual, "expected to be equal")
	}
}

// NotEqual asserts that two values are not equal.
func (a *Assert) NotEqual(expected, actual interface{}) {
	if reflect.DeepEqual(expected, actual) {
		a.reportError(expected, actual, "expected to be not equal")
	}
}

// True asserts that a value is true.
func (a *Assert) True(value bool) {
	if !value {
		a.reportError(true, value, "expected to be true")
	}
}

// False asserts that a value is false.
func (a *Assert) False(value bool) {
	if value {
		a.reportError(false, value, "expected to be false")
	}
}

// reportError is a helper function to report test failures.
func (a *Assert) reportError(expected, actual interface{}, message string) {
	a.errorMsg = fmt.Sprintf("%s - Expected %v, Actual: %v", message, expected, actual)
}

// Error returns the error message if the assertion failed.
func (a *Assert) Error() string {
	return a.errorMsg
}

