package diff

import "fmt"

// DiffResult represents the result of comparing two strings.
type DiffResult struct {
	HasDiff  bool   // Whether the strings differ
	Summary  string // Human-readable summary of the difference
	Position *int   // Position where strings first differ (nil if no difference)
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
	}
}