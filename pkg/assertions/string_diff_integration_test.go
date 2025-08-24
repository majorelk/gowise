package assertions

import (
	"strings"
	"testing"
)

// TestStringDiffIntegration tests that the diff infrastructure is properly integrated
// with the Equal assertion for enhanced string error messages.
func TestStringDiffIntegration(t *testing.T) {
	tests := []struct {
		name                string
		got, want           string
		expectErrorContains []string // Parts that should be in error message
	}{
		{
			name: "simple string difference",
			got:  "hello world",
			want: "hello World",
			expectErrorContains: []string{
				"string values differ at position 6",
				`got:  "hello world"`,
				`want: "hello World"`,
			},
		},
		{
			name: "Unicode string difference",
			got:  "Hello 🌍 World",
			want: "Hello 🌎 World",
			expectErrorContains: []string{
				"string values differ at rune position 6",
				`got:  "Hello 🌍 World"`,
				`want: "Hello 🌎 World"`,
			},
		},
		{
			name: "multi-line string difference",
			got: `line 1
line 2 modified
line 3`,
			want: `line 1
line 2
line 3`,
			expectErrorContains: []string{
				"difference at line 2", // Enhanced format now uses "difference at line"
				`got:  "line 1\nline 2 modified\nline 3"`,
				`want: "line 1\nline 2\nline 3"`,
				"context:", // Enhanced format includes context
			},
		},
		{
			name: "long string with context",
			got:  "this is a very long string that should trigger context window functionality",
			want: "this is a very long string that should trigger context windoW functionality",
			expectErrorContains: []string{
				"string values differ at position",
				"diff:",
				`got:  "this is a very long string that should trigger context window functionality"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a dummy testing.T that captures failures
			dummyT := &capturingT{}
			assert := New(dummyT)

			// This should fail and generate an error message
			assert.Equal(tt.got, tt.want)

			errorMsg := assert.Error()
			if errorMsg == "" {
				t.Fatalf("Expected error message but got none")
			}

			// Check that all expected parts are in the error message
			for _, expected := range tt.expectErrorContains {
				if !strings.Contains(errorMsg, expected) {
					t.Errorf("Error message missing expected content %q\nFull error message:\n%s", expected, errorMsg)
				}
			}
		})
	}
}

// TestStringDiffVsNonString ensures we don't break existing behaviour for non-string types
func TestStringDiffVsNonString(t *testing.T) {
	dummyT := &capturingT{}
	assert := New(dummyT)

	// Test non-string types still get the old error format
	assert.Equal(42, 24)

	errorMsg := assert.Error()
	if errorMsg == "" {
		t.Fatalf("Expected error message but got none")
	}

	// Should contain the traditional format
	expectedParts := []string{
		"values differ",
		"got:  42",
		"want: 24",
	}

	for _, expected := range expectedParts {
		if !strings.Contains(errorMsg, expected) {
			t.Errorf("Error message missing expected content %q\nFull error message:\n%s", expected, errorMsg)
		}
	}

	// Should NOT contain diff-specific content
	if strings.Contains(errorMsg, "position") || strings.Contains(errorMsg, "diff:") {
		t.Errorf("Non-string comparison should not use diff infrastructure, got: %s", errorMsg)
	}
}

// TestEnhancedMultiLineDiffIntegration tests enhanced multi-line diff through public API
func TestEnhancedMultiLineDiffIntegration(t *testing.T) {
	tests := []struct {
		name                string
		got, want           string
		expectErrorContains []string
	}{
		{
			name: "multi-line with context behaviour",
			got: `line 1
line 2 original
line 3`,
			want: `line 1
line 2 changed
line 3`,
			expectErrorContains: []string{
				"difference at line 2",
				`got:  "line 1\nline 2 original\nline 3"`,
				`want: "line 1\nline 2 changed\nline 3"`,
				"context:",
			},
		},
		{
			name: "multi-line different lengths behaviour",
			got: `line 1
line 2
line 3`,
			want: `line 1
line 2
line 3
line 4`,
			expectErrorContains: []string{
				"difference at line", // The algorithm detects line 3 as first different (content mismatch)
				"context:",
				`got:  "line 1\nline 2\nline 3"`,
				"line 4", // Should show the additional line in context
			},
		},
		{
			name: "empty vs multi-line behaviour",
			got:  "",
			want: `line 1
line 2`,
			expectErrorContains: []string{
				"difference at line 1",
				`got:  ""`,
				`want: "line 1\nline 2"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a dummy testing.T that captures failures
			dummyT := &capturingT{}
			assert := New(dummyT)

			// This should fail and generate enhanced error message
			assert.Equal(tt.got, tt.want)

			errorMsg := assert.Error()
			if errorMsg == "" {
				t.Fatalf("Expected error message but got none")
			}

			// Check that all expected parts are in the error message
			for _, expected := range tt.expectErrorContains {
				if !strings.Contains(errorMsg, expected) {
					t.Errorf("Error message missing expected content %q\nFull error message:\n%s", expected, errorMsg)
				}
			}
		})
	}
}

// capturingT is a minimal testing.T implementation for testing our assertions
type capturingT struct {
	failed bool
	logs   []string
}

func (t *capturingT) Helper() {}

func (t *capturingT) Errorf(format string, args ...interface{}) {
	t.failed = true
	// We don't actually log since our assertions use Error() method
}

func (t *capturingT) Fatalf(format string, args ...interface{}) {
	t.failed = true
	// We don't actually log since our assertions use Error() method
}
