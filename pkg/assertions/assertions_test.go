package assertions

import (
	"testing"
)

func TestAssertions(t *testing.T) {
	// Test Equal
	assert := New(t)
	assert.Equal(42,42)

	// Test NotEqual
	assert = New(t)
	assert.NotEqual(42,23)

	// Test True
	assert = New(t)
	assert.True(true)

	// Test False
	assert = New(t)
	assert.False(false)

	// Test failures
	assert = New(t)
	assert.Equal(55,22)
	assert.NotEqual(23,23)
	assert.True(false) // Expecting true to fail
	assert.False(true) // Expecting false to fail

	// Check for assertion failures
	if assert.Error() != "" {
		t.Error(assert.Error())
	}
}

