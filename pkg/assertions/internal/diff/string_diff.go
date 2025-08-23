package diff

import "fmt"

// DiffResult represents the result of comparing two strings.
type DiffResult struct {
	HasDiff  bool   // Whether the strings differ
	Summary  string // Human-readable summary of the difference
	Position *int   // Position where strings first differ (nil if no difference)
	Context  string // Context window around the difference
}

// StringDiff compares two strings and returns a DiffResult indicating
// where they differ, if at all.
func StringDiff(got, want string) DiffResult {
	if got == want {
		return DiffResult{
			HasDiff:  false,
			Summary:  "",
			Position: nil,
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
		HasDiff:  true,
		Summary:  fmt.Sprintf("string values differ at position %d", pos),
		Position: &pos,
		Context:  "",
	}
}

// StringDiffWithContext compares two strings and returns a DiffResult with
// a context window around the difference.
func StringDiffWithContext(got, want string, contextSize int) DiffResult {
	if got == want {
		return DiffResult{
			HasDiff:  false,
			Summary:  "",
			Position: nil,
			Context:  "",
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
		HasDiff:  true,
		Summary:  fmt.Sprintf("string values differ at position %d", pos),
		Position: &pos,
		Context:  contextStr,
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