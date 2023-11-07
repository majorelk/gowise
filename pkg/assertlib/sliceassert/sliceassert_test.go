// Package sliceassert provides slice comparison functions for testing.
//
// SliceAssertTest contains unit tests for the sliceassert package.
package sliceassert

import (
	"testing"
)

// TestSliceIsEqual_Pass tests SliceIsEqual with equal slices.
func TestSliceIsEqual_Pass(t *testing.T) {
	expectedSlice := []int{1, 2, 3}
	actualSlice := []int{1, 2, 3}
	
	SliceIsEqual(t,actualSlice, expectedSlice)
}

// TestSliceIsEqual_FailDifferentLengths tests SliceIsEqual with different length slices.
func TestSliceIsEqual_FailDifferentLengths(t *testing.T) {
	expectedSlice := []int{1, 2, 3}
	actualSlice := []int{1, 2, 3, 4}
	
	SliceIsEqual(t,actualSlice, expectedSlice)
}

// TestSliceIsEqual_FailNotDeeplyEqual tests SliceIsEqual with not deeply equal slices.
func TestSliceIsEqual_FailNotDeeplyEqual(t *testing.T) {
	expectedSlice := []int{1, 2, 3}
	actualSlice := []int{3, 2, 1}
	
	SliceIsEqual(t,actualSlice, expectedSlice)
}

// TestSliceIsEqual_FailDifferentTypes tests SliceIsEqual with slices of different types.
func TestSliceIsEqual_FailDifferentTypes(t *testing.T) {
	expectedSlice := []int{1, 2, 3}
	actualSlice := []string{"1", "2", "3"}

	SliceIsEqual(t,actualSlice, expectedSlice)
}

// TestSliceIsEqual_FailNotSlice tests SliceIsEqual with a non-slice value.
func TestSliceIsEqual_FailNotSlice(t *testing.T) {
	expectedSlice := []int{1, 2, 3}
	actualNotSlice := 123 // Not a slice

	SliceIsEqual(t,actualNotSlice, expectedSlice) // Should fail with Undefined error
}

