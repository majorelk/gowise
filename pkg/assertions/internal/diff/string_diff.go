package diff

import "fmt"

// DiffResult represents the result of comparing two strings.
type DiffResult struct {
	HasDiff    bool   // Whether the strings differ
	Summary    string // Human-readable summary of the difference
	Position   *int   // Position where strings first differ (nil if no difference)
	Context    string // Context window around the difference
	LineNumber *int   // Line number where multi-line strings first differ (nil if single-line or no difference)
}

// StringDiff compares two strings and returns a DiffResult indicating
// where they differ, if at all.
func StringDiff(got, want string) DiffResult {
	if got == want {
		return DiffResult{
			HasDiff:    false,
			Summary:    "",
			Position:   nil,
			Context:    "",
			LineNumber: nil,
		}
	}

	// Find the first position where strings differ
	minLen := len(got)
	if len(want) < minLen {
		minLen = len(want)
	}

	pos := 0
	for i := 0; i < minLen; i++ {
		if got[i] != want[i] {
			pos = i
			break
		}
	}

	// If all characters match up to minLen but lengths differ
	if pos == 0 && minLen > 0 && got[:minLen] == want[:minLen] {
		pos = minLen
	}

	return DiffResult{
		HasDiff:    true,
		Summary:    fmt.Sprintf("string values differ at position %d", pos),
		Position:   &pos,
		Context:    "",
		LineNumber: nil,
	}
}

// StringDiffWithContext compares two strings and returns a DiffResult with
// a context window around the difference.
func StringDiffWithContext(got, want string, contextSize int) DiffResult {
	if got == want {
		return DiffResult{
			HasDiff:    false,
			Summary:    "",
			Position:   nil,
			Context:    "",
			LineNumber: nil,
		}
	}

	// Find the first position where strings differ
	minLen := len(got)
	if len(want) < minLen {
		minLen = len(want)
	}

	pos := 0
	for i := 0; i < minLen; i++ {
		if got[i] != want[i] {
			pos = i
			break
		}
	}

	// If all characters match up to minLen but lengths differ
	if pos == 0 && minLen > 0 && got[:minLen] == want[:minLen] {
		pos = minLen
	}

	// Generate context window
	contextStr := generateContext(got, want, pos, contextSize)

	return DiffResult{
		HasDiff:    true,
		Summary:    fmt.Sprintf("string values differ at position %d", pos),
		Position:   &pos,
		Context:    contextStr,
		LineNumber: nil,
	}
}

// generateContext creates a context window around the difference position.
func generateContext(got, want string, pos, contextSize int) string {
	// For short strings, show the entire strings
	totalLen := contextSize*2 + 1
	if len(got) <= totalLen && len(want) <= totalLen {
		return fmt.Sprintf("%s vs %s", got, want)
	}

	// Calculate start and end positions for context window
	start := pos - contextSize
	if start < 0 {
		start = 0
	}

	gotEnd := pos + contextSize
	if gotEnd > len(got) {
		gotEnd = len(got)
	}

	wantEnd := pos + contextSize
	if wantEnd > len(want) {
		wantEnd = len(want)
	}

	// Extract context windows
	gotContext := got[start:gotEnd]
	wantContext := want[start:wantEnd]

	// Add ellipses if we're showing partial strings
	if start > 0 {
		gotContext = "..." + gotContext
		wantContext = "..." + wantContext
	}

	if gotEnd < len(got) {
		gotContext = gotContext + "..."
	}
	if wantEnd < len(want) {
		wantContext = wantContext + "..."
	}

	return fmt.Sprintf("%s vs %s", gotContext, wantContext)
}

// MultiLineStringDiff compares multi-line strings and returns a DiffResult
// with line number information for the first differing line.
func MultiLineStringDiff(got, want string) DiffResult {
	if got == want {
		return DiffResult{
			HasDiff:    false,
			Summary:    "",
			Position:   nil,
			Context:    "",
			LineNumber: nil,
		}
	}

	// Split into lines for comparison
	gotLines := splitLines(got)
	wantLines := splitLines(want)

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

	var summary string
	if len(gotLines) != len(wantLines) {
		summary = fmt.Sprintf("strings differ in length: got %d lines, want %d lines", len(gotLines), len(wantLines))
	} else {
		summary = fmt.Sprintf("strings differ at line %d", lineNum)
	}

	// Calculate character position for the differing line
	pos := 0
	for i := 0; i < lineNum-1 && i < len(gotLines); i++ {
		pos += len(gotLines[i]) + 1 // +1 for newline character
	}

	return DiffResult{
		HasDiff:    true,
		Summary:    summary,
		Position:   &pos,
		Context:    "",
		LineNumber: &lineNum,
	}
}

// splitLines splits a string into lines, preserving line endings.
func splitLines(s string) []string {
	if s == "" {
		return []string{}
	}

	var lines []string
	start := 0

	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i+1])
			start = i + 1
		}
	}

	// Add remaining content if it doesn't end with newline
	if start < len(s) {
		lines = append(lines, s[start:])
	}

	return lines
}

