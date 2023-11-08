// TestSliceIsEqual_Failing contains intentionally failing tests.

// +build failing_tests

package sliceassert

func TestSliceIsEqual_Failing(t *testing.T) {
	t.Run("FailDifferentLengths", func(t *testing.T) {
		expectedSlice := []int{1, 2, 3}
		actualSlice := []int{1, 2, 3, 4}

		SliceIsEqual(t, actualSlice, expectedSlice)
	})

	t.Run("FailNotDeeplyEqual", func(t *testing.T) {
		expectedSlice := []int{1, 2, 3}
		actualSlice := []int{3, 2, 1}

		SliceIsEqual(t, actualSlice, expectedSlice)
	})

	t.Run("FailDifferentTypes", func(t *testing.T) {
		expectedSlice := []int{1, 2, 3}
		actualSlice := []string{"1", "2", "3"}

		SliceIsEqual(t, actualSlice, expectedSlice)
	})

	t.Run("FailNotSlice", func(t *testing.T) {
		expectedSlice := []int{1, 2, 3}
		actualNotSlice := 123 // Not a slice

		SliceIsEqual(t, actualNotSlice, expectedSlice) // Should fail with an "Assertion failed" error
	})
}
