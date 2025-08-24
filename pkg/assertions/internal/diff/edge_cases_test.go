package diff

import (
	"strings"
	"testing"
)

// TestStringDiffEdgeCases tests comprehensive edge cases to ensure robust behaviour
func TestStringDiffEdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		got, want      string
		expectDiff     bool
		checkBehaviour func(t *testing.T, result DiffResult)
	}{
		{
			name:       "both empty strings",
			got:        "",
			want:       "",
			expectDiff: false,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.Position != nil {
					t.Error("Expected nil position for identical empty strings")
				}
			},
		},
		{
			name:       "empty vs non-empty",
			got:        "",
			want:       "text",
			expectDiff: true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.Position == nil || *result.Position != 0 {
					t.Errorf("Expected position 0 for empty vs non-empty, got %v", result.Position)
				}
			},
		},
		{
			name:       "non-empty vs empty",
			got:        "text",
			want:       "",
			expectDiff: true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.Position == nil || *result.Position != 0 {
					t.Errorf("Expected position 0 for non-empty vs empty, got %v", result.Position)
				}
			},
		},
		{
			name:       "whitespace only differences",
			got:        "hello world",
			want:       "hello  world", // double space
			expectDiff: true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.Position == nil || *result.Position != 6 {
					t.Errorf("Expected position 6 for whitespace diff, got %v", result.Position)
				}
			},
		},
		{
			name:       "trailing whitespace",
			got:        "hello",
			want:       "hello ",
			expectDiff: true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.Position == nil || *result.Position != 5 {
					t.Errorf("Expected position 5 for trailing whitespace, got %v", result.Position)
				}
			},
		},
		{
			name:       "leading whitespace",
			got:        " hello",
			want:       "hello",
			expectDiff: true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.Position == nil || *result.Position != 0 {
					t.Errorf("Expected position 0 for leading whitespace, got %v", result.Position)
				}
			},
		},
		{
			name:       "tab vs spaces",
			got:        "hello\tworld",
			want:       "hello world",
			expectDiff: true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.Position == nil || *result.Position != 5 {
					t.Errorf("Expected position 5 for tab vs space, got %v", result.Position)
				}
			},
		},
		{
			name:       "null byte in string",
			got:        "hello\x00world",
			want:       "helloworld",
			expectDiff: true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.Position == nil || *result.Position != 5 {
					t.Errorf("Expected position 5 for null byte, got %v", result.Position)
				}
			},
		},
		{
			name:       "very long identical strings",
			got:        strings.Repeat("a", 10000),
			want:       strings.Repeat("a", 10000),
			expectDiff: false,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.HasDiff {
					t.Error("Expected no diff for identical long strings")
				}
			},
		},
		{
			name:       "very long strings with difference at end",
			got:        strings.Repeat("a", 9999) + "b",
			want:       strings.Repeat("a", 10000),
			expectDiff: true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.Position == nil || *result.Position != 9999 {
					t.Errorf("Expected position 9999 for diff at end, got %v", result.Position)
				}
			},
		},
		{
			name:       "single character strings",
			got:        "a",
			want:       "b",
			expectDiff: true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.Position == nil || *result.Position != 0 {
					t.Errorf("Expected position 0 for single char diff, got %v", result.Position)
				}
			},
		},
		{
			name:       "identical single characters",
			got:        "x",
			want:       "x",
			expectDiff: false,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.HasDiff {
					t.Error("Expected no diff for identical single chars")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringDiff(tt.got, tt.want)

			if result.HasDiff != tt.expectDiff {
				t.Errorf("HasDiff = %v, want %v", result.HasDiff, tt.expectDiff)
			}

			// Run custom behaviour checks
			if tt.checkBehaviour != nil {
				tt.checkBehaviour(t, result)
			}
		})
	}
}

// TestMultiLineEdgeCases tests comprehensive edge cases for multi-line string diff
func TestMultiLineEdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		got, want      string
		expectDiff     bool
		checkBehaviour func(t *testing.T, result DiffResult)
	}{
		{
			name:       "mixed line endings - CRLF vs LF",
			got:        "line1\r\nline2\r\nline3",
			want:       "line1\nline2\nline3",
			expectDiff: true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.LineNumber == nil || *result.LineNumber != 1 {
					t.Errorf("Expected line number 1 for CRLF vs LF, got %v", result.LineNumber)
				}
			},
		},
		{
			name:       "mixed line endings - LF vs CR",
			got:        "line1\nline2\nline3",
			want:       "line1\rline2\rline3",
			expectDiff: true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.LineNumber == nil || *result.LineNumber != 1 {
					t.Errorf("Expected line number 1 for LF vs CR, got %v", result.LineNumber)
				}
			},
		},
		{
			name: "empty lines at start",
			got: `
line1
line2`,
			want: `line1
line2`,
			expectDiff: true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.LineNumber == nil || *result.LineNumber != 1 {
					t.Errorf("Expected line number 1 for empty line at start, got %v", result.LineNumber)
				}
			},
		},
		{
			name: "empty lines at end",
			got: `line1
line2
`,
			want: `line1
line2`,
			expectDiff: true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				// The algorithm correctly detects difference at line 2 (trailing newline)
				if result.LineNumber == nil || *result.LineNumber != 2 {
					t.Errorf("Expected line number 2 for trailing newline diff, got %v", result.LineNumber)
				}
				if !strings.Contains(result.Summary, "line 2") {
					t.Errorf("Expected summary to mention line 2, got: %s", result.Summary)
				}
			},
		},
		{
			name: "only empty lines",
			got: `


`,
			want: `

`,
			expectDiff: true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if !result.HasDiff {
					t.Error("Expected diff for different number of empty lines")
				}
			},
		},
		{
			name:       "single line with newline vs without",
			got:        "hello\n",
			want:       "hello",
			expectDiff: true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if !result.HasDiff {
					t.Error("Expected diff for line with vs without newline")
				}
			},
		},
		{
			name:       "very long single line",
			got:        strings.Repeat("word ", 1000) + "\n",
			want:       strings.Repeat("word ", 1000) + "\n",
			expectDiff: false,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.HasDiff {
					t.Error("Expected no diff for identical long lines")
				}
			},
		},
		{
			name:       "many short lines",
			got:        strings.Repeat("x\n", 1000),
			want:       strings.Repeat("x\n", 1000),
			expectDiff: false,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.HasDiff {
					t.Error("Expected no diff for identical many lines")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MultiLineStringDiff(tt.got, tt.want)

			if result.HasDiff != tt.expectDiff {
				t.Errorf("HasDiff = %v, want %v", result.HasDiff, tt.expectDiff)
			}

			// Run custom behaviour checks
			if tt.checkBehaviour != nil {
				tt.checkBehaviour(t, result)
			}
		})
	}
}

// TestUnicodeEdgeCases tests comprehensive edge cases for Unicode string handling
func TestUnicodeEdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		got, want      string
		expectDiff     bool
		checkBehaviour func(t *testing.T, result DiffResult)
	}{
		{
			name:       "combining diacritical marks",
			got:        "café",       // single é character
			want:       "cafe\u0301", // e + combining acute accent
			expectDiff: true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.Position == nil || *result.Position != 3 {
					t.Errorf("Expected position 3 for combining marks, got %v", result.Position)
				}
			},
		},
		{
			name:       "zero-width characters",
			got:        "hello\u200bworld", // zero-width space
			want:       "helloworld",
			expectDiff: true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.Position == nil || *result.Position != 5 {
					t.Errorf("Expected position 5 for zero-width char, got %v", result.Position)
				}
			},
		},
		{
			name:       "right-to-left text",
			got:        "hello שלום world",
			want:       "hello שלום world",
			expectDiff: false,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.HasDiff {
					t.Error("Expected no diff for identical RTL text")
				}
			},
		},
		{
			name:       "emoji variations",
			got:        "👍🏻", // thumbs up with light skin tone
			want:       "👍",  // basic thumbs up
			expectDiff: true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.Position == nil || *result.Position != 1 {
					t.Errorf("Expected position 1 for emoji variation, got %v", result.Position)
				}
			},
		},
		{
			name:       "mixed scripts",
			got:        "Hello こんにちは мир 🌍",
			want:       "Hello こんにちは мир 🌎",
			expectDiff: true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				// Should find difference at emoji position
				if !result.HasDiff {
					t.Error("Expected diff for mixed scripts with different emoji")
				}
			},
		},
		{
			name:       "normalization forms",
			got:        "\u00e9",       // é as single character
			want:       "\u0065\u0301", // e + combining accent
			expectDiff: true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				// These are different at Unicode level even if visually identical
				if !result.HasDiff {
					t.Error("Expected diff for different normalization forms")
				}
			},
		},
		{
			name:       "surrogate pairs",
			got:        "𝕳𝖊𝖑𝖑𝖔", // mathematical script letters
			want:       "𝕳𝖊𝖑𝖑𝖔",
			expectDiff: false,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.HasDiff {
					t.Error("Expected no diff for identical surrogate pairs")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UnicodeStringDiff(tt.got, tt.want)

			if result.HasDiff != tt.expectDiff {
				t.Errorf("HasDiff = %v, want %v", result.HasDiff, tt.expectDiff)
			}

			// Run custom behaviour checks
			if tt.checkBehaviour != nil {
				tt.checkBehaviour(t, result)
			}
		})
	}
}

// TestContextWindowEdgeCases tests edge cases for context window functionality
func TestContextWindowEdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		got, want      string
		contextSize    int
		expectDiff     bool
		checkBehaviour func(t *testing.T, result DiffResult)
	}{
		{
			name:        "zero context size",
			got:         "hello world",
			want:        "hello World",
			contextSize: 0,
			expectDiff:  true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.Context == "" {
					t.Error("Expected some context even with size 0")
				}
			},
		},
		{
			name:        "negative context size",
			got:         "hello world",
			want:        "hello World",
			contextSize: -5,
			expectDiff:  true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				// Should handle gracefully, not panic
				if !result.HasDiff {
					t.Error("Expected diff despite negative context size")
				}
			},
		},
		{
			name:        "context size larger than string",
			got:         "hi",
			want:        "Hi",
			contextSize: 100,
			expectDiff:  true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if result.Context != "hi vs Hi" {
					t.Errorf("Expected full string context, got: %q", result.Context)
				}
			},
		},
		{
			name:        "difference at very end",
			got:         "this is a long string ending with X",
			want:        "this is a long string ending with Y",
			contextSize: 5,
			expectDiff:  true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if !strings.Contains(result.Context, "X") || !strings.Contains(result.Context, "Y") {
					t.Errorf("Expected context to show difference at end, got: %q", result.Context)
				}
			},
		},
		{
			name:        "difference at very start",
			got:         "X this is a long string",
			want:        "Y this is a long string",
			contextSize: 5,
			expectDiff:  true,
			checkBehaviour: func(t *testing.T, result DiffResult) {
				if !strings.Contains(result.Context, "X") || !strings.Contains(result.Context, "Y") {
					t.Errorf("Expected context to show difference at start, got: %q", result.Context)
				}
				// Should show some context around the difference
				if len(result.Context) < 10 {
					t.Errorf("Expected meaningful context, got: %q", result.Context)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringDiffWithContext(tt.got, tt.want, tt.contextSize)

			if result.HasDiff != tt.expectDiff {
				t.Errorf("HasDiff = %v, want %v", result.HasDiff, tt.expectDiff)
			}

			// Run custom behaviour checks
			if tt.checkBehaviour != nil {
				tt.checkBehaviour(t, result)
			}
		})
	}
}
