// Package sliceassert provides slice comparison functions for testing.
//
// SliceAssertTest contains unit tests for the sliceassert package.

package sliceassert

import (
	"testing"
)

// TestSliceIsEqual contains the regular tests for the sliceassert package.
func TestSliceIsEqual(t *testing.T) {
	t.Run("Pass", func(t *testing.T) {
		expectedSlice := []int{1, 2, 3}
		actualSlice := []int{1, 2, 3}

		SliceIsEqual(t, actualSlice, expectedSlice)
	})
}

