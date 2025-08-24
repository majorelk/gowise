package assertions

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

// TestComprehensiveEnhancedDiffIntegration tests enhanced diff integration across assertion types
func TestComprehensiveEnhancedDiffIntegration(t *testing.T) {
	tests := []struct {
		name                string
		setupAndAssert      func(assert *Assert)
		expectErrorContains []string
	}{
		{
			name: "HasPrefix with multi-line string",
			setupAndAssert: func(assert *Assert) {
				got := "Hello\nWorld\nFoo"
				want := "Hello\nworld\nFoo" // Different case on line 2
				assert.HasPrefix(got, want)
			},
			expectErrorContains: []string{
				"expected to have prefix",
				"difference at line 2", // Enhanced diff should show line difference
			},
		},
		{
			name: "HasSuffix with multi-line string",
			setupAndAssert: func(assert *Assert) {
				got := "Foo\nBar\nBaz"
				want := "Bar\nBaz\nExtra" // Different suffix
				assert.HasSuffix(got, want)
			},
			expectErrorContains: []string{
				"expected to have suffix",
				"unified diff:", // Multi-line diff will use unified format
			},
		},
		{
			name: "BodyContains with large response body",
			setupAndAssert: func(assert *Assert) {
				responseBody := `{
  "status": "success",
  "message": "Operation completed",
  "data": {
    "id": 123,
    "name": "test"
  }
}`
				response := &http.Response{
					Body: io.NopCloser(bytes.NewBufferString(responseBody)),
				}
				expectedContent := `{
  "status": "failure",
  "message": "Operation failed"
}`
				assert.BodyContains(response, expectedContent)
			},
			expectErrorContains: []string{
				"expected body to contain",
				"unified diff:", // Complex diff should use unified format
			},
		},
		{
			name: "HeaderEqual with string difference",
			setupAndAssert: func(assert *Assert) {
				response := &http.Response{
					Header: http.Header{
						"Content-Type": []string{"application/json; charset=utf-8"},
					},
				}
				assert.HeaderEqual(response, "Content-Type", "application/json; charset=iso-8859-1")
			},
			expectErrorContains: []string{
				"expected different header value",
				// HeaderEqual compares string to []string, so no enhanced diff
				"got:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dummyT := &capturingT{}
			assert := New(dummyT)

			// Execute the assertion (should fail)
			tt.setupAndAssert(assert)

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

// TestBroadEnhancedDiffCompatibility tests that enhanced diff works across many assertion types
func TestBroadEnhancedDiffCompatibility(t *testing.T) {
	tests := []struct {
		name           string
		setupAndAssert func(assert *Assert)
		shouldUseeDiff bool // Whether this should trigger enhanced diff
	}{
		{
			name: "Equal with strings",
			setupAndAssert: func(assert *Assert) {
				assert.Equal("hello\nworld", "hello\nWorld")
			},
			shouldUseeDiff: true,
		},
		{
			name: "NotEqual with identical strings (should not fail)",
			setupAndAssert: func(assert *Assert) {
				// This should pass, so we won't get an error message to test
				assert.NotEqual("hello", "world")
			},
			shouldUseeDiff: false,
		},
		{
			name: "DeepEqual with strings",
			setupAndAssert: func(assert *Assert) {
				assert.DeepEqual("line1\nline2\nline3", "line1\nLINE2\nline3")
			},
			shouldUseeDiff: true,
		},
		{
			name: "Same with different string pointers",
			setupAndAssert: func(assert *Assert) {
				s1 := "hello"
				s2 := "hello"
				// Force different pointers
				s2 = strings.Replace(s2, "h", "h", 1)
				assert.Same(&s1, &s2)
			},
			shouldUseeDiff: false, // Same compares pointers, not values
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dummyT := &capturingT{}
			assert := New(dummyT)

			// Execute the assertion
			tt.setupAndAssert(assert)

			errorMsg := assert.Error()

			if tt.shouldUseeDiff {
				// Should have enhanced diff features
				if errorMsg == "" {
					t.Fatalf("Expected error message for enhanced diff test")
				}
				if !strings.Contains(errorMsg, "difference at line") && !strings.Contains(errorMsg, "string values differ at position") {
					t.Errorf("Expected enhanced diff markers in error message, got: %s", errorMsg)
				}
			}
			// For non-enhanced diff tests, we just verify they don't crash
		})
	}
}
