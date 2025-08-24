package diff

import (
	"strings"
	"testing"
)

// TestEnhancedMultiLineStringDiff tests enhanced multi-line diff with context windows
func TestEnhancedMultiLineStringDiff(t *testing.T) {
	tests := []struct {
		name         string
		got, want    string
		contextLines int
		expectDiff   bool
		checkResult  func(t *testing.T, result EnhancedDiffResult)
	}{
		{
			name: "simple multi-line with context",
			got: `line 1
line 2 original
line 3`,
			want: `line 1
line 2 changed
line 3`,
			contextLines: 1,
			expectDiff:   true,
			checkResult: func(t *testing.T, result EnhancedDiffResult) {
				if result.LineNumber == nil || *result.LineNumber != 2 {
					t.Errorf("Expected line number 2, got %v", result.LineNumber)
				}
				if !strings.Contains(result.ContextLines, "line 1") {
					t.Errorf("Expected context to include line 1, got: %q", result.ContextLines)
				}
				if !strings.Contains(result.ContextLines, "line 3") {
					t.Errorf("Expected context to include line 3, got: %q", result.ContextLines)
				}
			},
		},
		{
			name: "large file with limited context",
			got: `line 1
line 2
line 3
line 4 original
line 5
line 6
line 7`,
			want: `line 1
line 2
line 3
line 4 changed
line 5
line 6  
line 7`,
			contextLines: 2,
			expectDiff:   true,
			checkResult: func(t *testing.T, result EnhancedDiffResult) {
				if result.LineNumber == nil || *result.LineNumber != 4 {
					t.Errorf("Expected line number 4, got %v", result.LineNumber)
				}
				// Should show lines 2-6 (2 before, diff line, 2 after)
				lines := strings.Split(result.ContextLines, "\n")
				if len(lines) < 5 {
					t.Errorf("Expected at least 5 context lines, got %d", len(lines))
				}
			},
		},
		{
			name: "difference at start of file",
			got: `line 1 original
line 2
line 3`,
			want: `line 1 changed
line 2
line 3`,
			contextLines: 2,
			expectDiff:   true,
			checkResult: func(t *testing.T, result EnhancedDiffResult) {
				if result.LineNumber == nil || *result.LineNumber != 1 {
					t.Errorf("Expected line number 1, got %v", result.LineNumber)
				}
				// Should show lines from start since diff is at line 1
				if !strings.Contains(result.ContextLines, "line 2") {
					t.Errorf("Expected context to include line 2, got: %q", result.ContextLines)
				}
			},
		},
		{
			name: "difference at end of file",
			got: `line 1
line 2
line 3 original`,
			want: `line 1
line 2
line 3 changed`,
			contextLines: 2,
			expectDiff:   true,
			checkResult: func(t *testing.T, result EnhancedDiffResult) {
				if result.LineNumber == nil || *result.LineNumber != 3 {
					t.Errorf("Expected line number 3, got %v", result.LineNumber)
				}
				// Should show lines leading up to end since diff is at end
				if !strings.Contains(result.ContextLines, "line 2") {
					t.Errorf("Expected context to include line 2, got: %q", result.ContextLines)
				}
			},
		},
		{
			name: "identical multi-line strings",
			got: `line 1
line 2
line 3`,
			want: `line 1
line 2
line 3`,
			contextLines: 1,
			expectDiff:   false,
			checkResult: func(t *testing.T, result EnhancedDiffResult) {
				if result.HasDiff {
					t.Errorf("Expected no diff for identical strings")
				}
				if result.ContextLines != "" {
					t.Errorf("Expected empty context for identical strings, got: %q", result.ContextLines)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EnhancedMultiLineStringDiff(tt.got, tt.want, tt.contextLines)

			if result.HasDiff != tt.expectDiff {
				t.Errorf("HasDiff = %v, want %v", result.HasDiff, tt.expectDiff)
			}

			// Run custom result checks
			if tt.checkResult != nil {
				tt.checkResult(t, result)
			}
		})
	}
}

// TestUnifiedDiffFormat tests unified diff format output
func TestUnifiedDiffFormat(t *testing.T) {
	tests := []struct {
		name              string
		got, want         string
		expectUnifiedDiff func(t *testing.T, unifiedDiff string)
	}{
		{
			name: "basic unified diff",
			got: `line 1
line 2 original
line 3`,
			want: `line 1
line 2 changed
line 3`,
			expectUnifiedDiff: func(t *testing.T, unifiedDiff string) {
				if !strings.Contains(unifiedDiff, "@@") {
					t.Errorf("Expected unified diff to contain @@ markers, got: %q", unifiedDiff)
				}
				if !strings.Contains(unifiedDiff, "-line 2 original") {
					t.Errorf("Expected unified diff to show removed line, got: %q", unifiedDiff)
				}
				if !strings.Contains(unifiedDiff, "+line 2 changed") {
					t.Errorf("Expected unified diff to show added line, got: %q", unifiedDiff)
				}
			},
		},
		{
			name: "multiple changes unified diff",
			got: `line 1
line 2 original
line 3
line 4 original`,
			want: `line 1
line 2 changed
line 3
line 4 changed`,
			expectUnifiedDiff: func(t *testing.T, unifiedDiff string) {
				// Should show both changes in unified format
				removedCount := strings.Count(unifiedDiff, "-line")
				addedCount := strings.Count(unifiedDiff, "+line")
				if removedCount < 2 || addedCount < 2 {
					t.Errorf("Expected at least 2 removed and 2 added lines, got %d removed, %d added", removedCount, addedCount)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EnhancedMultiLineStringDiff(tt.got, tt.want, 3)

			if result.UnifiedDiff == "" {
				t.Errorf("Expected unified diff output, got empty string")
			}

			// Run custom unified diff checks
			if tt.expectUnifiedDiff != nil {
				tt.expectUnifiedDiff(t, result.UnifiedDiff)
			}
		})
	}
}

// TestSideBySideDiffBehaviour tests side-by-side diff behavioural output
func TestSideBySideDiffBehaviour(t *testing.T) {
	tests := []struct {
		name                 string
		got, want            string
		expectSideBySideDiff func(t *testing.T, sideBySideDiff string)
	}{
		{
			name: "basic side-by-side diff behaviour",
			got: `line 1
line 2 original
line 3`,
			want: `line 1
line 2 changed
line 3`,
			expectSideBySideDiff: func(t *testing.T, sideBySideDiff string) {
				// Behaviour: must contain column headers for user orientation
				if !strings.Contains(sideBySideDiff, "Got") || !strings.Contains(sideBySideDiff, "Want") {
					t.Errorf("Behaviour: expected side-by-side diff to contain Got/Want headers for user orientation, got: %q", sideBySideDiff)
				}
				// Behaviour: must show both versions of differing content
				if !strings.Contains(sideBySideDiff, "line 2 original") {
					t.Errorf("Behaviour: expected side-by-side diff to show original line content, got: %q", sideBySideDiff)
				}
				if !strings.Contains(sideBySideDiff, "line 2 changed") {
					t.Errorf("Behaviour: expected side-by-side diff to show changed line content, got: %q", sideBySideDiff)
				}
				// Behaviour: must show identical lines for context
				if !strings.Contains(sideBySideDiff, "line 1") || !strings.Contains(sideBySideDiff, "line 3") {
					t.Errorf("Behaviour: expected side-by-side diff to show identical lines for context, got: %q", sideBySideDiff)
				}
			},
		},
		{
			name: "different line counts behaviour",
			got: `line 1
line 2
line 3`,
			want: `line 1
line 2
line 3
line 4`,
			expectSideBySideDiff: func(t *testing.T, sideBySideDiff string) {
				// Behaviour: must show additional lines that exist in one version
				if !strings.Contains(sideBySideDiff, "line 4") {
					t.Errorf("Behaviour: expected side-by-side diff to show additional line from 'want', got: %q", sideBySideDiff)
				}
				// Behaviour: output must accommodate all lines from both versions
				lines := strings.Split(sideBySideDiff, "\n")
				if len(lines) < 5 { // Headers + separator + at least 4 content lines
					t.Errorf("Behaviour: expected side-by-side output to show all lines from both versions, got %d lines", len(lines))
				}
			},
		},
		{
			name: "empty strings behaviour",
			got:  "",
			want: "",
			expectSideBySideDiff: func(t *testing.T, sideBySideDiff string) {
				// Behaviour: should still show headers for consistency
				if !strings.Contains(sideBySideDiff, "Got") || !strings.Contains(sideBySideDiff, "Want") {
					t.Errorf("Behaviour: expected headers even for empty strings, got: %q", sideBySideDiff)
				}
			},
		},
		{
			name: "very long lines behaviour",
			got:  "this is a very long line that exceeds typical display width and should be truncated properly",
			want: "this is a very long line that exceeds typical display width and should be handled well",
			expectSideBySideDiff: func(t *testing.T, sideBySideDiff string) {
				// Behaviour: must handle long lines without breaking layout
				lines := strings.Split(sideBySideDiff, "\n")
				for _, line := range lines {
					if len(line) > 100 { // Reasonable terminal width
						t.Errorf("Behaviour: expected long lines to be handled without breaking layout, found line of length %d", len(line))
					}
				}
				// Behaviour: must indicate truncation occurred
				if !strings.Contains(sideBySideDiff, "...") {
					t.Errorf("Behaviour: expected truncation indicator for long lines, got: %q", sideBySideDiff)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EnhancedMultiLineStringDiff(tt.got, tt.want, 3)

			// Behaviour: side-by-side diff should be generated for all comparisons
			if result.SideBySideDiff == "" {
				t.Errorf("Behaviour: expected side-by-side diff output to be generated, got empty string")
			}

			// Run custom behavioural checks
			if tt.expectSideBySideDiff != nil {
				tt.expectSideBySideDiff(t, result.SideBySideDiff)
			}
		})
	}
}
