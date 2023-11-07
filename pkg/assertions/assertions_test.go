package assertions

import (
	"testing"
	"fmt"
)

func TestAssertions(t *testing.T) {
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
			t.Run(fmt.Sprintf("Test case %d", i+1), func(t *testing.T) {
				assert.Equal(tc.expected, tc.actual)

				if tc.pass && assert.Error() != "" {
					t.Errorf("Test case %d failed, expected no error but got: %s", i+1, assert.Error())
				} else if !tc.pass && assert.Error() == "" {
					t.Errorf("Test case %d failed, expected an error but got none", i+1)
				}
			})
		}
	})

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
			t.Run(fmt.Sprintf("Test case %d", i+1), func(t *testing.T) {
				assert.NotEqual(tc.expected, tc.actual)

				if tc.pass && assert.Error() != "" {
					t.Errorf("Test case %d failed, expected no error but got: %s", i+1, assert.Error())
				} else if !tc.pass && assert.Error() == "" {
					t.Errorf("Test case %d failed, expected an error but got none", i+1)
				}
			})
		}
	})

	t.Run("True", func(t *testing.T) {
		assert := New(t)

		// Test cases
		testCases := []struct {
			value bool
			pass  bool
		}{
			{true, true},
			{false, false},
			{true, false},
			{false, true},
		}

		for i, tc := range testCases {
			t.Run(fmt.Sprintf("Test case %d", i+1), func(t *testing.T) {
				assert.True(tc.value)

				if tc.pass && assert.Error() != "" {
					t.Errorf("Test case %d failed, expected no error but got: %s", i+1, assert.Error())
				} else if !tc.pass && assert.Error() == "" {
					t.Errorf("Test case %d failed, expected an error but got none", i+1)
				}
			})
		}
	})

	t.Run("False", func(t *testing.T) {
		assert := New(t)

		// Test cases
		testCases := []struct {
			value bool
			pass  bool
		}{
			{true, false},
			{false, true},
			{true, true},
			{false, false},
		}

		for i, tc := range testCases {
			t.Run(fmt.Sprintf("Test case %d", i+1), func(t *testing.T) {
				assert.False(tc.value)

				if tc.pass && assert.Error() != "" {
					t.Errorf("Test case %d failed, expected no error but got: %s", i+1, assert.Error())
				} else if !tc.pass && assert.Error() == "" {
					t.Errorf("Test case %d failed, expected an error but got none", i+1)
				}
			})
		}
	})
}
