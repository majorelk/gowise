package testrunner

import (
	"testing"
	"gowise/pkg/assertions" // Import the assertions package
)

func TestCreateTestRunner(t *testing.T) {
	tr := NewTestRunner(t)

	tr.RunTest("TestCreateTestRunner", func(assert *assertions.Assert) error { 
		// Use assertions.Assert type
		// Testing if 2 + 2 equals 4
		assert.True(2+2 == 4)
		return nil
	})
}

