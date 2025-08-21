// Package sliceassert provides slice comparison functions for testing.
//
// SliceAssert is a utility package that offers a function to compare slices for
// equality, allowing deep comparison.
package sliceassert

import (
	"reflect"
	"testing"
)

// SliceIsEqual checks if two slices are deeply equal.
func SliceIsEqual(t *testing.T, actual, expected interface{}) {
	actualValue := reflect.ValueOf(actual)
	expectedValue := reflect.ValueOf(expected)

	if actualValue.Kind() != reflect.Slice || expectedValue.Kind() != reflect.Slice {
		t.Errorf("Assertion failed: expected a slice, but got %T and %T", actual, expected)
		return
	}

	if actualValue.Len() != expectedValue.Len() {
		t.Errorf("Assertion failed: slices have different lengths")
		return
	}

	for i := 0; i < actualValue.Len(); i++ {
		if !reflect.DeepEqual(actualValue.Index(i).Interface(), expectedValue.Index(i).Interface()) {
			t.Errorf("Assertion failed: elements at index %d are not deeply equal", i)
			return
		}
	}
}
