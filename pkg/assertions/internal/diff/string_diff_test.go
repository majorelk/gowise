package diff

import (
	"testing"
)

// TestDiffResult tests the basic DiffResult struct functionality.
func TestDiffResult(t *testing.T) {
	t.Run("no difference", func(t *testing.T) {
		result := StringDiff("hello", "hello")

		if result.HasDiff {
			t.Errorf("Expected no difference for identical strings")
		}

		if result.Summary != "" {
			t.Errorf("Expected empty summary for identical strings, got: %q", result.Summary)
		}

		if result.Position != nil {
			t.Errorf("Expected nil position for identical strings, got: %d", *result.Position)
		}
	})

	t.Run("simple difference", func(t *testing.T) {
		result := StringDiff("hello", "world")

		if !result.HasDiff {
			t.Errorf("Expected difference for different strings")
		}

		if result.Summary == "" {
			t.Errorf("Expected non-empty summary for different strings")
		}

		if result.Position == nil {
			t.Errorf("Expected position for different strings")
		} else if *result.Position != 0 {
			t.Errorf("Expected position 0 for completely different strings, got: %d", *result.Position)
		}
	})
}

// TestStringDiffBasic tests the basic string diff functionality.
func TestStringDiffBasic(t *testing.T) {
	tests := []struct {
		name            string
		got, want       string
		expectDiff      bool
		expectedPos     *int // nil if no diff expected
		expectedSummary string
	}{
		{
			name:            "identical strings",
			got:             "hello",
			want:            "hello",
			expectDiff:      false,
			expectedPos:     nil,
			expectedSummary: "",
		},
		{
			name:            "different at start",
			got:             "Hello",
			want:            "hello",
			expectDiff:      true,
			expectedPos:     intPtr(0),
			expectedSummary: "string values differ at position 0",
		},
		{
			name:            "different in middle",
			got:             "hello world",
			want:            "hello World",
			expectDiff:      true,
			expectedPos:     intPtr(6),
			expectedSummary: "string values differ at position 6",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringDiff(tt.got, tt.want)

			if result.HasDiff != tt.expectDiff {
				t.Errorf("HasDiff = %v, want %v", result.HasDiff, tt.expectDiff)
			}

			if tt.expectedPos == nil && result.Position != nil {
				t.Errorf("Expected nil position, got %d", *result.Position)
			} else if tt.expectedPos != nil {
				if result.Position == nil {
					t.Errorf("Expected position %d, got nil", *tt.expectedPos)
				} else if *result.Position != *tt.expectedPos {
					t.Errorf("Position = %d, want %d", *result.Position, *tt.expectedPos)
				}
			}

			if result.Summary != tt.expectedSummary {
				t.Errorf("Summary = %q, want %q", result.Summary, tt.expectedSummary)
			}
		})
	}
}

// TestStringDiffWithContext tests context window functionality around differences.
func TestStringDiffWithContext(t *testing.T) {
	tests := []struct {
		name        string
		got, want   string
		contextSize int
		expectDiff  bool
	}{
		{
			name:        "short strings no context needed",
			got:         "abc",
			want:        "abd",
			contextSize: 5,
			expectDiff:  true,
		},
		{
			name:        "long strings with context",
			got:         "the quick brown fox jumps",
			want:        "the quick black fox jumps",
			contextSize: 3,
			expectDiff:  true,
		},
		{
			name:        "difference at start",
			got:         "Hello world this is a test",
			want:        "hello world this is a test",
			contextSize: 5,
			expectDiff:  true,
		},
		{
			name:        "difference at end",
			got:         "this is a test case",
			want:        "this is a test cast",
			contextSize: 4,
			expectDiff:  true,
		},
		{
			name:        "identical strings",
			got:         "same string",
			want:        "same string",
			contextSize: 3,
			expectDiff:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringDiffWithContext(tt.got, tt.want, tt.contextSize)

			if result.HasDiff != tt.expectDiff {
				t.Errorf("HasDiff = %v, want %v", result.HasDiff, tt.expectDiff)
			}

			if tt.expectDiff {
				if result.Position == nil {
					t.Errorf("Expected position for different strings")
				}
				if result.Context == "" {
					t.Errorf("Expected non-empty context for different strings")
				}
				if result.Summary == "" {
					t.Errorf("Expected non-empty summary for different strings")
				}
			} else {
				if result.Position != nil {
					t.Errorf("Expected nil position for identical strings")
				}
				if result.Context != "" {
					t.Errorf("Expected empty context for identical strings, got: %q", result.Context)
				}
			}
		})
	}
}

// TestMultiLineStringDiff tests diff functionality for multi-line strings.
func TestMultiLineStringDiff(t *testing.T) {
	tests := []struct {
		name       string
		got, want  string
		expectDiff bool
	}{
		{
			name: "identical multi-line strings",
			got: `line 1
line 2
line 3`,
			want: `line 1
line 2
line 3`,
			expectDiff: false,
		},
		{
			name: "different line content",
			got: `line 1
line 2 modified
line 3`,
			want: `line 1
line 2
line 3`,
			expectDiff: true,
		},
		{
			name: "extra line at end",
			got: `line 1
line 2
line 3
line 4`,
			want: `line 1
line 2
line 3`,
			expectDiff: true,
		},
		{
			name: "missing line",
			got: `line 1
line 3`,
			want: `line 1
line 2
line 3`,
			expectDiff: true,
		},
		{
			name:       "different line endings",
			got:        "line 1\nline 2\n",
			want:       "line 1\r\nline 2\r\n",
			expectDiff: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MultiLineStringDiff(tt.got, tt.want)

			if result.HasDiff != tt.expectDiff {
				t.Errorf("HasDiff = %v, want %v", result.HasDiff, tt.expectDiff)
			}

			if tt.expectDiff {
				if result.LineNumber == nil {
					t.Errorf("Expected line number for different multi-line strings")
				}
				if result.Summary == "" {
					t.Errorf("Expected non-empty summary for different multi-line strings")
				}
			} else {
				if result.LineNumber != nil {
					t.Errorf("Expected nil line number for identical strings")
				}
			}
		})
	}
}

// TestUnicodeStringDiff tests diff functionality with Unicode characters.
func TestUnicodeStringDiff(t *testing.T) {
	tests := []struct {
		name       string
		got, want  string
		expectDiff bool
		expectPos  *int // Expected rune position, not byte position
	}{
		{
			name:       "identical Unicode strings",
			got:        "Hello 🌍 World",
			want:       "Hello 🌍 World",
			expectDiff: false,
			expectPos:  nil,
		},
		{
			name:       "different Unicode characters",
			got:        "Hello 🌍 World",
			want:       "Hello 🌎 World",
			expectDiff: true,
			expectPos:  intPtr(6), // Position of the emoji (6th rune)
		},
		{
			name:       "ASCII vs Unicode",
			got:        "café",
			want:       "cafe",
			expectDiff: true,
			expectPos:  intPtr(3), // Position of é (4th rune, 0-indexed = 3)
		},
		{
			name:       "multi-byte character at start",
			got:        "🚀 launch",
			want:       "🛸 launch",
			expectDiff: true,
			expectPos:  intPtr(0), // First character differs
		},
		{
			name:       "combining characters",
			got:        "é",          // single character
			want:       "e\u0301",   // e + combining acute accent  
			expectDiff: true,
			expectPos:  intPtr(1), // They differ at the second position (combining accent)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UnicodeStringDiff(tt.got, tt.want)
			
			if result.HasDiff != tt.expectDiff {
				t.Errorf("HasDiff = %v, want %v", result.HasDiff, tt.expectDiff)
			}
			
			if tt.expectDiff {
				if result.Position == nil {
					t.Errorf("Expected position for different strings")
				} else if tt.expectPos != nil && *result.Position != *tt.expectPos {
					t.Errorf("Position = %d, want %d", *result.Position, *tt.expectPos)
				}
				if result.Summary == "" {
					t.Errorf("Expected non-empty summary for different strings")
				}
			} else {
				if result.Position != nil {
					t.Errorf("Expected nil position for identical strings")
				}
			}
		})
	}
}

// Helper function for test readability
func intPtr(i int) *int {
	return &i
}

