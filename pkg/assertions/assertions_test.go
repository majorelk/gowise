package assertions

import (
	"testing"
)

func TestAssertions(t *testing.T) {
	// Define test cases for the Equal function.
	t.Run("Equal", func(t *testing.T) {
		assert := New(t)

		// Test cases
		testCases := []struct {
			expected, actual interface{}
			pass             bool
		}{
			{42, 42, true},
			{42, 23, false},
			{"hello", "world", false},
			{true, true, true},
			{false, true, false},
		}

		for i, tc := range testCases {
			t.Run("Test case "+string(i+1), func(t *testing.T) {
				assert.Equal(tc.expected, tc.actual)

				if tc.pass && assert.Error() != "" {
					t.Errorf("Test case %d failed, expected no error but got: %s", i+1, assert.Error())
				} else if !tc.pass && assert.Error() == "" {
					t.Errorf("Test case %d failed, expected an error but got none", i+1)
				}
			})
		}
	})

	// Define test cases for the NotEqual function.
	t.Run("NotEqual", func(t *testing.T) {
		assert := New(t)

		// Test cases
		testCases := []struct {
			expected, actual interface{}
			pass             bool
		}{
			{42, 42, false},
			{42, 23, true},
			{"hello", "world", true},
			{true, true, false},
			{false, true, true},
		}

		for i, tc := range testCases {
			t.Run("Test case "+string(i+1), func(t *testing.T) {
				assert.NotEqual(tc.expected, tc.actual)

				if tc.pass && assert.Error() != "" {
					t.Errorf("Test case %d failed, expected no error but got: %s", i+1, assert.Error())
				} else if !tc.pass && assert.Error() == "" {
					t.Errorf("Test case %d failed, expected an error but got none", i+1)
				}
			})
		}
	})
}

