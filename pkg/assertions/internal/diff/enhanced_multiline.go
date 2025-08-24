package diff

import (
	"fmt"
	"strings"
)

// EnhancedDiffResult represents enhanced multi-line diff results with context and unified output
type EnhancedDiffResult struct {
	HasDiff        bool   // Whether the strings differ
	LineNumber     *int   // Line number where strings first differ (1-indexed, nil if no difference)
	ContextLines   string // Lines around the difference with context window
	UnifiedDiff    string // Unified diff format output
	SideBySideDiff string // Side-by-side diff format output
}

// EnhancedMultiLineStringDiff compares multi-line strings with enhanced context and formatting
func EnhancedMultiLineStringDiff(got, want string, contextLines int) EnhancedDiffResult {
	// Split into lines for comparison
	gotLines := splitLines(got)
	wantLines := splitLines(want)

	// Generate side-by-side diff format (always generated for consistency)
	sideBySideDiff := generateSideBySideDiff(gotLines, wantLines)

	if got == want {
		return EnhancedDiffResult{
			HasDiff:        false,
			LineNumber:     nil,
			ContextLines:   "",
			UnifiedDiff:    "",
			SideBySideDiff: sideBySideDiff,
		}
	}

	// Find first differing line
	minLines := len(gotLines)
	if len(wantLines) < minLines {
		minLines = len(wantLines)
	}

	lineNum := 0
	for i := 0; i < minLines; i++ {
		if gotLines[i] != wantLines[i] {
			lineNum = i + 1 // 1-indexed line numbers
			break
		}
	}

	// If all lines match but lengths differ
	if lineNum == 0 && len(gotLines) != len(wantLines) {
		lineNum = minLines + 1
	}

	// Generate context window around the difference
	contextStr := generateContextWindow(gotLines, wantLines, lineNum-1, contextLines)

	// Generate unified diff format
	unifiedDiff := generateUnifiedDiff(gotLines, wantLines)

	return EnhancedDiffResult{
		HasDiff:        true,
		LineNumber:     &lineNum,
		ContextLines:   contextStr,
		UnifiedDiff:    unifiedDiff,
		SideBySideDiff: sideBySideDiff,
	}
}

// generateContextWindow creates a context window around the differing line
func generateContextWindow(gotLines, wantLines []string, diffLineIdx, contextLines int) string {
	if len(gotLines) == 0 && len(wantLines) == 0 {
		return ""
	}

	// Calculate the range of lines to show
	start := diffLineIdx - contextLines
	if start < 0 {
		start = 0
	}

	maxLines := len(gotLines)
	if len(wantLines) > maxLines {
		maxLines = len(wantLines)
	}

	end := diffLineIdx + contextLines + 1
	if end > maxLines {
		end = maxLines
	}

	var contextBuilder strings.Builder

	// Show context lines
	for i := start; i < end; i++ {
		var gotLine, wantLine string

		if i < len(gotLines) {
			gotLine = strings.TrimSuffix(gotLines[i], "\n")
		}
		if i < len(wantLines) {
			wantLine = strings.TrimSuffix(wantLines[i], "\n")
		}

		// Add line to context
		if i == diffLineIdx {
			// This is the differing line - show both versions
			if gotLine != "" {
				contextBuilder.WriteString(fmt.Sprintf("- %s\n", gotLine))
			}
			if wantLine != "" {
				contextBuilder.WriteString(fmt.Sprintf("+ %s\n", wantLine))
			}
		} else if gotLine == wantLine {
			// Identical context line
			contextBuilder.WriteString(fmt.Sprintf("  %s\n", gotLine))
		} else {
			// Different context line
			if gotLine != "" {
				contextBuilder.WriteString(fmt.Sprintf("- %s\n", gotLine))
			}
			if wantLine != "" {
				contextBuilder.WriteString(fmt.Sprintf("+ %s\n", wantLine))
			}
		}
	}

	return strings.TrimSuffix(contextBuilder.String(), "\n")
}

// generateUnifiedDiff creates a unified diff format output
func generateUnifiedDiff(gotLines, wantLines []string) string {
	var result strings.Builder

	// Header
	result.WriteString("--- got\n")
	result.WriteString("+++ want\n")

	// Find continuous blocks of differences
	i, j := 0, 0
	for i < len(gotLines) || j < len(wantLines) {
		// Skip identical lines
		for i < len(gotLines) && j < len(wantLines) &&
			strings.TrimSuffix(gotLines[i], "\n") == strings.TrimSuffix(wantLines[j], "\n") {
			i++
			j++
		}

		if i >= len(gotLines) && j >= len(wantLines) {
			break
		}

		// Find end of difference block
		blockStartI, blockStartJ := i, j
		blockEndI, blockEndJ := i, j

		// Advance through different lines
		for blockEndI < len(gotLines) || blockEndJ < len(wantLines) {
			if blockEndI < len(gotLines) && blockEndJ < len(wantLines) {
				if strings.TrimSuffix(gotLines[blockEndI], "\n") == strings.TrimSuffix(wantLines[blockEndJ], "\n") {
					break
				}
				blockEndI++
				blockEndJ++
			} else if blockEndI < len(gotLines) {
				blockEndI++
			} else {
				blockEndJ++
			}
		}

		// Output hunk header
		gotCount := blockEndI - blockStartI
		wantCount := blockEndJ - blockStartJ
		result.WriteString(fmt.Sprintf("@@ -%d,%d +%d,%d @@\n", blockStartI+1, gotCount, blockStartJ+1, wantCount))

		// Output removed lines
		for k := blockStartI; k < blockEndI; k++ {
			line := strings.TrimSuffix(gotLines[k], "\n")
			result.WriteString(fmt.Sprintf("-%s\n", line))
		}

		// Output added lines
		for k := blockStartJ; k < blockEndJ; k++ {
			line := strings.TrimSuffix(wantLines[k], "\n")
			result.WriteString(fmt.Sprintf("+%s\n", line))
		}

		i = blockEndI
		j = blockEndJ
	}

	return strings.TrimSuffix(result.String(), "\n")
}

// generateSideBySideDiff creates a side-by-side diff format output
func generateSideBySideDiff(gotLines, wantLines []string) string {
	var result strings.Builder

	// Headers
	result.WriteString("Got                           | Want\n")
	result.WriteString("------------------------------|------------------------------\n")

	// Calculate max lines to process
	maxLines := len(gotLines)
	if len(wantLines) > maxLines {
		maxLines = len(wantLines)
	}

	// Process each line
	for i := 0; i < maxLines; i++ {
		var gotLine, wantLine string

		// Get lines, using empty string if beyond array bounds
		if i < len(gotLines) {
			gotLine = strings.TrimSuffix(gotLines[i], "\n")
		}
		if i < len(wantLines) {
			wantLine = strings.TrimSuffix(wantLines[i], "\n")
		}

		// Truncate lines if too long for display (keep first 29 chars)
		if len(gotLine) > 29 {
			gotLine = gotLine[:26] + "..."
		}
		if len(wantLine) > 29 {
			wantLine = wantLine[:26] + "..."
		}

		// Format the line pair with proper alignment
		result.WriteString(fmt.Sprintf("%-29s | %s\n", gotLine, wantLine))
	}

	return strings.TrimSuffix(result.String(), "\n")
}
