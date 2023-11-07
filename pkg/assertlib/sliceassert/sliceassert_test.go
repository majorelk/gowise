package sliceassert

import (
	"testing"
)

func TestSliceIsEqual_Pass(t *testing.T) {
	expectedSlice := []int{1, 2, 3}
	actualSlice := []int{1, 2, 3}
	
	SliceIsEqual(t,actualSlice, expectedSlice)
}

func TestSliceIsEqual_FailDifferentLengths(t *testing.T) {
	expectedSlice := []int{1, 2, 3}
	actualSlice := []int{1, 2, 3, 4}
	
	SliceIsEqual(t,actualSlice, expectedSlice)
}

func TestSliceIsEqual_FailNotDeeplyEqual(t *testing.T) {
	expectedSlice := []int{1, 2, 3}
	actualSlice := []int{3, 2, 1}
	
	SliceIsEqual(t,actualSlice, expectedSlice)
}

func TestSliceIsEqual_FailDifferentTypes(t *testing.T) {
	expectedSlice := []int{1, 2, 3}
	actualSlice := []string{"1", "2", "3"}

	SliceIsEqual(t,actualSlice, expectedSlice)
}

func TestSliceIsEqual_FailNotSlice(t *testing.T) {
	expectedSlice := []int{1, 2, 3}
	actualNotSlice := 123 // Not a slice

	SliceIsEqual(t,actualNotSlice, expectedSlice) // Should fail with Undefined error
}

