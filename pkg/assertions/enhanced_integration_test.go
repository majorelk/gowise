package assertions

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

// behaviorMockT is defined in assertions_passing_test.go - shared across test files

// TestEnhancedJsonEqualIntegration tests JsonEqual with enhanced diff
func TestEnhancedJsonEqualIntegration(t *testing.T) {
	tests := []struct {
		name                string
		actual, expected    string
		expectErrorContains []string
	}{
		{
			name: "JSON objects with different values",
			actual: `{
  "name": "John Doe",
  "age": 25,
  "city": "London"
}`,
			expected: `{
  "name": "Jane Smith", 
  "age": 30,
  "city": "Manchester"
}`,
			expectErrorContains: []string{
				"JSON objects differ",
				"difference at line 2",
				`"name": "John Doe"`,
				`"name": "Jane Smith"`,
				"unified diff:", // Multiple changes trigger unified format
			},
		},
		{
			name:   "JSON semantic differences",
			actual: `{"name":"John","age":25,"active":true}`,
			expected: `{
  "name": "John",
  "age": 25,
  "active": false
}`,
			expectErrorContains: []string{
				"JSON objects differ",
				"unified diff:",
			},
		},
		{
			name:     "Invalid JSON handling",
			actual:   `{"name":"John","age":25}`,
			expected: `{"name":"John","age":25,}`, // Trailing comma is invalid
			expectErrorContains: []string{
				"expected JSON is invalid",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &behaviorMockT{}
			assert := New(mock)

			assert.JsonEqual(tt.expected, tt.actual)

			// Framework behavior: FAIL = exactly 1 Errorf call with expected content
			if len(mock.errorCalls) != 1 {
				t.Fatalf("Expected assertion to fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
			}

			errorMsg := mock.errorCalls[0]

			for _, expected := range tt.expectErrorContains {
				if !strings.Contains(errorMsg, expected) {
					t.Errorf("Error message missing expected content %q\nFull error message:\n%s", expected, errorMsg)
				}
			}
		})
	}
}

// TestEnhancedBodyJsonEqualIntegration tests BodyJsonEqual with enhanced diff
func TestEnhancedBodyJsonEqualIntegration(t *testing.T) {
	tests := []struct {
		name                string
		responseBody        string
		expected            interface{}
		expectErrorContains []string
	}{
		{
			name: "JSON response differs from expected string",
			responseBody: `{
  "status": "error",
  "message": "Not found"
}`,
			expected: `{
  "status": "success",
  "message": "Found"
}`,
			expectErrorContains: []string{
				"response JSON differs from expected",
				"context:",
				`"status": "error"`,
				`"status": "success"`,
			},
		},
		{
			name:         "Invalid JSON response body",
			responseBody: `{"status":"error","message":"Not found",}`, // Invalid trailing comma
			expected:     `{"status":"success"}`,
			expectErrorContains: []string{
				"response body is not valid JSON",
				"string values differ at position",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &behaviorMockT{}
			assert := New(mock)

			// Create mock HTTP response
			response := &http.Response{
				Body: io.NopCloser(bytes.NewBufferString(tt.responseBody)),
			}

			assert.BodyJsonEqual(response, tt.expected)

			// Framework behavior: FAIL = exactly 1 Errorf call with expected content
			if len(mock.errorCalls) != 1 {
				t.Fatalf("Expected assertion to fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
			}

			errorMsg := mock.errorCalls[0]

			for _, expected := range tt.expectErrorContains {
				if !strings.Contains(errorMsg, expected) {
					t.Errorf("Error message missing expected content %q\nFull error message:\n%s", expected, errorMsg)
				}
			}
		})
	}
}

// TestEnhancedDeepEqualIntegration tests DeepEqual with various scenarios
func TestEnhancedDeepEqualIntegration(t *testing.T) {
	tests := []struct {
		name                string
		got, want           interface{}
		expectErrorContains []string
	}{
		{
			name: "string comparison through DeepEqual",
			got:  "hello\nworld",
			want: "hello\nWorld",
			expectErrorContains: []string{
				"values differ",
				"difference at line 2",
				"context:",
			},
		},
		{
			name: "struct comparison",
			got: struct {
				Name string
				Age  int
			}{"John", 25},
			want: struct {
				Name string
				Age  int
			}{"Jane", 25},
			expectErrorContains: []string{
				"values differ",
				// DeepEqual compares struct as a whole, so no enhanced diff
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &behaviorMockT{}
			assert := New(mock)

			assert.DeepEqual(tt.got, tt.want)

			// Framework behavior: FAIL = exactly 1 Errorf call with expected content
			if len(mock.errorCalls) != 1 {
				t.Fatalf("Expected assertion to fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
			}

			errorMsg := mock.errorCalls[0]

			for _, expected := range tt.expectErrorContains {
				if !strings.Contains(errorMsg, expected) {
					t.Errorf("Error message missing expected content %q\nFull error message:\n%s", expected, errorMsg)
				}
			}
		})
	}
}
