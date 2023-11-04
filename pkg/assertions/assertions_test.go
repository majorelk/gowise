package assertions

import (
	"fmt"
	"testing"
)

func TestEqual(t *testing.T) {
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
			assert := New(t)
			assert.Equal(tc.expected, tc.actual)

			if tc.pass && assert.Error() != "" {
				t.Errorf("Test case %d failed, expected no error but got: %s", i+1, assert.Error())
			} else if !tc.pass && assert.Error() == "" {
				t.Errorf("Test case %d failed, expected an error but got none", i+1)
			}
		})
	}
}

func TestNotEqual(t *testing.T) {
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
			assert := New(t)
			assert.NotEqual(tc.expected, tc.actual)
			if tc.pass && assert.Error() != "" {
				t.Errorf("Test case %d failed, expected no error but got: %s", i+1, assert.Error())
			} else if !tc.pass && assert.Error() == "" {
				t.Errorf("Test case %d failed, expected an error but got none", i+1)
			}
		})
	}
}

func TestTrue(t *testing.T) {
	// Test cases for the True assertion
	// Add similar test cases as the above functions.
}

func TestFalse(t *testing.T) {
	// Test cases for the False assertion
	// Add similar test cases as the above functions.
}
