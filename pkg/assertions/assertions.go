// Package assertions provides fast, ergonomic assertion functions for testing.
//
// GoWise assertions offer performance-optimised equality checking, comprehensive
// nil handling, and type-safe collection operations while maintaining zero
// external dependencies.
//
// Core Features:
//   - Fast-path optimisation for comparable types (avoids reflection when possible)
//   - Comprehensive nil checking for all 6 nillable Go types
//   - Type-safe collection operations (Contains, Len)
//   - Clear, actionable error messages in UK English
//   - Proper t.Helper() integration for accurate test stack traces
//
// Example usage:
//
//	assert := assertions.New(t)
//	assert.Equal(got, want)        // Fast-path for comparable types, deep equality fallback
//	assert.Contains(slice, item)   // Works with slices, maps, strings
//	assert.Nil(err)               // Handles interface nil gotcha correctly
//	assert.True(condition)        // Clear boolean assertions
//
// All assertions are designed for minimal allocation and maximum clarity,
// following Go's stdlib-only philosophy.
package assertions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	"gowise/pkg/assertions/internal/diff"
)

// isComparable checks if two values can be compared with ==.
// This is a fast-path optimisation for common types.
func isComparable(a, b interface{}) bool {
	if a == nil || b == nil {
		return true
	}

	va, vb := reflect.ValueOf(a), reflect.ValueOf(b)

	// Must be same type to be comparable
	if va.Type() != vb.Type() {
		return false
	}

	// Check if the type is comparable
	return va.Type().Comparable()
}

// DiffFormat specifies the preferred format for multi-line string diffs
type DiffFormat int

const (
	// DiffFormatAuto automatically selects the best format based on content complexity
	DiffFormatAuto DiffFormat = iota
	// DiffFormatContext shows context lines with +/- indicators
	DiffFormatContext
	// DiffFormatUnified shows unified diff format with @@ headers
	DiffFormatUnified
	// DiffFormatSideBySide shows side-by-side comparison (not yet implemented in error messages)
	DiffFormatSideBySide
)

// Assert is a struct that holds the testing context and error message.
type Assert struct {
	t          interface{}
	errorMsg   string
	diffFormat DiffFormat // Preferred format for multi-line string diffs
}

// New creates a new Assert instance with the given testing context.
func New(t interface{}) *Assert {
	return &Assert{
		t:          t,
		diffFormat: DiffFormatAuto, // Default to automatic format selection
	}
}

// WithDiffFormat returns a new Assert instance with the specified diff format preference.
// This follows GoWise principles of immutable configuration.
func (a *Assert) WithDiffFormat(format DiffFormat) *Assert {
	return &Assert{
		t:          a.t,
		errorMsg:   a.errorMsg,
		diffFormat: format,
	}
}

// Equal asserts that two values are equal.
// Uses fast-path comparison for comparable types, falls back to reflect.DeepEqual.
func (a *Assert) Equal(got, want interface{}) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	// Fast path for nil comparison
	if got == nil && want == nil {
		return
	}
	if (got == nil) != (want == nil) {
		a.reportError(got, want, "values differ")
		return
	}

	// Fast path for comparable types using type assertion
	if isComparable(got, want) && got == want {
		return
	}
	if isComparable(got, want) && got != want {
		a.reportError(got, want, "values differ")
		return
	}

	// Fallback to deep equality
	if !reflect.DeepEqual(got, want) {
		a.reportError(got, want, "values differ")
	}
}

// NotEqual asserts that two values are not equal.
// Uses fast-path comparison for comparable types, falls back to reflect.DeepEqual.
func (a *Assert) NotEqual(got, want interface{}) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	// Fast path for nil comparison
	if got == nil && want == nil {
		a.reportError(got, want, "values should not be equal")
		return
	}
	if (got == nil) != (want == nil) {
		return // different nil states = not equal, which is what we want
	}

	// Fast path for comparable types
	if isComparable(got, want) {
		if got == want {
			a.reportError(got, want, "values should not be equal")
		}
		return
	}

	// Fallback to deep equality
	if reflect.DeepEqual(got, want) {
		a.reportError(got, want, "values should not be equal")
	}
}

// DeepEqual asserts that two values are deeply equal.
// Always uses reflect.DeepEqual for thorough comparison.
func (a *Assert) DeepEqual(got, want interface{}) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	if !reflect.DeepEqual(got, want) {
		a.reportError(got, want, "values differ")
	}
}

// Same asserts that two values have the same pointer identity.
// Uses == comparison which tests for pointer identity.
func (a *Assert) Same(got, want interface{}) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	// Use == for pointer identity comparison
	// This works for pointers, interfaces, channels, maps, slices, and functions
	if got == want {
		return
	}

	a.reportError(got, want, "expected same pointer identity")
}

// True asserts that a boolean condition is true.
func (a *Assert) True(condition bool) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	if !condition {
		a.reportError(true, condition, "expected condition to be true")
	}
}

// False asserts that a boolean condition is false.
func (a *Assert) False(condition bool) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	if condition {
		a.reportError(false, condition, "expected condition to be false")
	}
}

// reportError is a helper function to report test failures.
// Uses lazy formatting and UK English.
func (a *Assert) reportError(got, want interface{}, message string) {
	// Check if both values are strings and use diff for better error messages
	if gotStr, gotOK := got.(string); gotOK {
		if wantStr, wantOK := want.(string); wantOK {
			a.reportStringError(gotStr, wantStr, message)
			return
		}
	}

	// Default error message for non-string types
	a.errorMsg = fmt.Sprintf("%s\n  got:  %#v\n  want: %#v", message, got, want)
}

// reportCollectionError provides enhanced error messages for collection comparisons using diff infrastructure
func (a *Assert) reportCollectionError(result diff.CollectionDiffResult) {
	var errorMsg strings.Builder
	errorMsg.WriteString(result.Summary)

	if result.Detail != "" {
		errorMsg.WriteString("\n  ")
		errorMsg.WriteString(strings.ReplaceAll(result.Detail, "\n", "\n  "))
	}

	a.errorMsg = errorMsg.String()
}

// reportStringError provides enhanced error messages for string comparisons using diff infrastructure
func (a *Assert) reportStringError(got, want string, message string) {
	// Choose appropriate diff function based on string characteristics
	var result diff.DiffResult

	// Use enhanced multi-line diff for strings containing newlines
	if strings.Contains(got, "\n") || strings.Contains(want, "\n") {
		// Use more context for complex diffs
		contextLines := 3
		if len(strings.Split(got, "\n")) > 10 || len(strings.Split(want, "\n")) > 10 {
			contextLines = 5
		}
		enhanced := diff.EnhancedMultiLineStringDiff(got, want, contextLines)

		var errorMsg strings.Builder
		errorMsg.WriteString(message)
		errorMsg.WriteString(fmt.Sprintf("\n  got:  %q", got))
		errorMsg.WriteString(fmt.Sprintf("\n  want: %q", want))

		if enhanced.HasDiff && enhanced.LineNumber != nil {
			errorMsg.WriteString(fmt.Sprintf("\n  difference at line %d", *enhanced.LineNumber))
		}

		// Choose diff format based on configuration and content complexity
		if enhanced.ContextLines != "" {
			contextLines := strings.Split(enhanced.ContextLines, "\n")

			// Determine which format to use
			var useUnified bool
			switch a.diffFormat {
			case DiffFormatUnified:
				useUnified = true
			case DiffFormatContext:
				useUnified = false
			case DiffFormatAuto:
				// Count number of differing lines for automatic selection
				diffLineCount := 0
				for _, line := range contextLines {
					if strings.HasPrefix(strings.TrimSpace(line), "+") || strings.HasPrefix(strings.TrimSpace(line), "-") {
						diffLineCount++
					}
				}
				useUnified = diffLineCount > 4
			}

			if useUnified && enhanced.UnifiedDiff != "" {
				errorMsg.WriteString("\n  unified diff:\n")
				unifiedLines := strings.Split(enhanced.UnifiedDiff, "\n")
				for _, line := range unifiedLines {
					if strings.TrimSpace(line) != "" {
						errorMsg.WriteString("    " + line + "\n")
					}
				}
			} else {
				// Use context format
				errorMsg.WriteString("\n  context:\n")
				for _, line := range contextLines {
					if strings.TrimSpace(line) != "" {
						errorMsg.WriteString("    " + line + "\n")
					}
				}
			}
		}

		a.errorMsg = strings.TrimSuffix(errorMsg.String(), "\n")
		return
	} else if hasUnicodeChars(got) || hasUnicodeChars(want) {
		// Use Unicode diff for strings with multi-byte characters
		result = diff.UnicodeStringDiff(got, want)
	} else {
		// Use context diff for better readability on longer strings
		contextSize := 10
		if len(got) > 50 || len(want) > 50 {
			result = diff.StringDiffWithContext(got, want, contextSize)
		} else {
			result = diff.StringDiff(got, want)
		}
	}

	// Build enhanced error message
	var errorMsg strings.Builder
	errorMsg.WriteString(message)
	errorMsg.WriteString("\n")

	if result.Summary != "" {
		errorMsg.WriteString("  ")
		errorMsg.WriteString(result.Summary)
		errorMsg.WriteString("\n")
	}

	if result.Context != "" {
		errorMsg.WriteString("  diff: ")
		errorMsg.WriteString(result.Context)
		errorMsg.WriteString("\n")
	}

	// Always show the full values for reference
	errorMsg.WriteString(fmt.Sprintf("  got:  %q\n", got))
	errorMsg.WriteString(fmt.Sprintf("  want: %q", want))

	a.errorMsg = errorMsg.String()
}

// hasUnicodeChars checks if a string contains non-ASCII characters
func hasUnicodeChars(s string) bool {
	for _, r := range s {
		if r > 127 {
			return true
		}
	}
	return false
}

// Error returns the error message if the assertion failed.
func (a *Assert) Error() string {
	return a.errorMsg
}

// Nil asserts that a value is nil.
// Supports pointers, interfaces, slices, maps, channels, and functions.
func (a *Assert) Nil(value interface{}) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	if !isNil(value) {
		a.reportError(nil, value, "expected value to be nil")
	}
}

// NotNil asserts that a value is not nil.
// Supports pointers, interfaces, slices, maps, channels, and functions.
func (a *Assert) NotNil(value interface{}) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	if isNil(value) {
		a.reportError("not nil", value, "expected value to not be nil")
	}
}

// isNil is a helper function to check if a value is nil.
// Handles both untyped nil and typed nil values across different types.
func isNil(value interface{}) bool {
	// Handle untyped nil
	if value == nil {
		return true
	}

	// Use reflection to check typed nil values
	valueRef := reflect.ValueOf(value)
	valueType := valueRef.Kind()

	// Check types that can be nil
	switch valueType {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return valueRef.IsNil()
	default:
		return false
	}
}

// Contains asserts that a container includes a specific item.
// Supports strings (substring), slices, arrays, and maps (key lookup).
func (a *Assert) Contains(container, item interface{}) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	result := diff.CollectionContainsDiff(container, item)
	if result.HasDiff {
		a.reportCollectionError(result)
	}
}

// Greater asserts that the first value is greater than the second.
func (a *Assert) Greater(v1, v2 float64) {
	if v1 <= v2 {
		a.reportError(v1, v2, "expected to be greater")
	}
}

// Less asserts that the first value is less than the second.
func (a *Assert) Less(v1, v2 float64) {
	if v1 >= v2 {
		a.reportError(v1, v2, "expected to be less")
	}
}

// HasPrefix asserts that a string starts with a certain substring.
func (a *Assert) HasPrefix(s, prefix string) {
	if !strings.HasPrefix(s, prefix) {
		a.reportError(prefix, s, "expected to have prefix")
	}
}

// HasSuffix asserts that a string ends with a certain substring.
func (a *Assert) HasSuffix(s, suffix string) {
	if !strings.HasSuffix(s, suffix) {
		a.reportError(suffix, s, "expected to have suffix")
	}
}

// InDelta asserts that the difference between two numeric values is within a certain range.
func (a *Assert) InDelta(expected, actual, delta float64) {
	if math.Abs(expected-actual) > delta {
		a.reportError(expected, actual, fmt.Sprintf("expected difference to be within %v", delta))
	}
}

// InEpsilon asserts that the difference between two numeric values is within a certain percentage.
func (a *Assert) InEpsilon(expected, actual, epsilon float64) {
	if math.Abs((expected-actual)/((expected+actual)/2)) > epsilon {
		a.reportError(expected, actual, fmt.Sprintf("expected difference to be within %v percent", epsilon*100))
	}
}

// Regexp asserts that a string matches a regular expression.
func (a *Assert) Regexp(pattern, str string) {
	if matched, _ := regexp.MatchString(pattern, str); !matched {
		a.reportError(pattern, str, "expected to match regular expression")
	}
}

// NoError asserts that a function call returns no error.
func (a *Assert) NoError(err error) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	if err != nil {
		a.reportError("no error", err, "expected no error")
	}
}

// ErrorType asserts that a function call returns a specific error type.
func (a *Assert) ErrorType(expected, actual error) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	expectedType := reflect.TypeOf(expected)
	actualType := reflect.TypeOf(actual)

	if expectedType != actualType {
		// Use cleaner type name formatting for error message
		expectedTypeName := expectedType.String()
		actualTypeName := actualType.String()
		a.errorMsg = fmt.Sprintf("expected error of different type\n  expected: %s\n  actual:   %s", expectedTypeName, actualTypeName)
	}
}

// HasError asserts that an error occurred (error is not nil).
func (a *Assert) HasError(err error) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	if err == nil {
		a.reportError("an error", nil, "expected an error but got none")
	}
}

// ErrorIs asserts that an error matches a target error using errors.Is.
// This follows Go 1.13+ error wrapping patterns.
func (a *Assert) ErrorIs(err, target error) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	if !errors.Is(err, target) {
		a.reportError(target, err, "expected error to match target")
	}
}

// ErrorAs asserts that an error can be assigned to a target type using errors.As.
// This follows Go 1.13+ error wrapping patterns.
func (a *Assert) ErrorAs(err error, target interface{}) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	if !errors.As(err, target) {
		a.reportError(reflect.TypeOf(target).Elem(), err, "expected error to be assignable to target type")
	}
}

// ErrorContains asserts that an error's message contains a specific substring.
func (a *Assert) ErrorContains(err error, substring string) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	if err == nil {
		a.reportError(substring, nil, "expected error but got nil")
		return
	}

	errorMessage := err.Error()
	if !strings.Contains(errorMessage, substring) {
		a.reportError(substring, errorMessage, "expected error message to contain")
	}
}

// ErrorMatches asserts that an error's message matches a regular expression pattern.
func (a *Assert) ErrorMatches(err error, pattern string) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	if err == nil {
		a.reportError(pattern, nil, "expected error but got nil")
		return
	}

	errorMessage := err.Error()
	matched, regexErr := regexp.MatchString(pattern, errorMessage)
	if regexErr != nil {
		a.reportError(pattern, regexErr, "invalid regular expression pattern")
		return
	}

	if !matched {
		// Use direct error message format to avoid string diff confusion
		// Show raw pattern (no quotes) for better readability
		a.errorMsg = fmt.Sprintf("expected error message to match pattern\n  pattern: %s\n  error:   %q", pattern, errorMessage)
	}
}

// IsEmpty asserts that a given array, slice, map, or string is empty.
func (a *Assert) IsEmpty(value interface{}) {
	valueType := reflect.TypeOf(value)
	switch valueType.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String:
		if reflect.ValueOf(value).Len() != 0 {
			a.reportError(0, reflect.ValueOf(value).Len(), "expected to be empty")
		}
	default:
		a.reportError(nil, value, "invalid type for IsEmpty")
	}
}

// IsNotEmpty asserts that a given array, slice, map, or string is not empty.
func (a *Assert) IsNotEmpty(value interface{}) {
	valueType := reflect.TypeOf(value)
	switch valueType.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String:
		if reflect.ValueOf(value).Len() == 0 {
			a.reportError("not empty", reflect.ValueOf(value).Len(), "expected to be not empty")
		}
	default:
		a.reportError(nil, value, "invalid type for IsNotEmpty")
	}
}

// Len asserts that a container has the expected length.
// Supports strings, slices, arrays, maps, and channels.
func (a *Assert) Len(container interface{}, expectedLen int) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	result := diff.CollectionLenDiff(container, expectedLen)
	if result.HasDiff {
		a.reportCollectionError(result)
	}
}

// Implements asserts that an object implements a certain interface.
func (a *Assert) Implements(object, interfaceObj interface{}) {
	objectType := reflect.TypeOf(object)
	interfaceType := reflect.TypeOf(interfaceObj).Elem()
	if !objectType.Implements(interfaceType) {
		a.reportError(interfaceType, objectType, "expected to implement interface")
	}
}

// IsZero asserts that a given numeric value or time.Time is zero.
func (a *Assert) IsZero(value interface{}) {
	switch v := value.(type) {
	case int, int8, int16, int32, int64:
		if v != 0 {
			a.reportError(0, v, "expected to be zero")
		}
	case uint, uint8, uint16, uint32, uint64:
		if v != 0 {
			a.reportError(0, v, "expected to be zero")
		}
	case float32, float64:
		if v != 0 {
			a.reportError(0, v, "expected to be zero")
		}
	case time.Time:
		if !v.IsZero() {
			a.reportError(time.Time{}, v, "expected to be zero")
		}
	default:
		a.reportError(nil, value, "invalid type for IsZero")
	}
}

// IsNotZero asserts that a given numeric value or time.Time is not zero.
func (a *Assert) IsNotZero(value interface{}) {
	switch v := value.(type) {
	case int, int8, int16, int32, int64:
		if v == 0 {
			a.reportError("not zero", v, "expected to be not zero")
		}
	case uint, uint8, uint16, uint32, uint64:
		if v == 0 {
			a.reportError("not zero", v, "expected to be not zero")
		}
	case float32, float64:
		if v == 0 {
			a.reportError("not zero", v, "expected to be not zero")
		}
	case time.Time:
		if v.IsZero() {
			a.reportError("not zero", v, "expected to be not zero")
		}
	default:
		a.reportError(nil, value, "invalid type for IsNotZero")
	}
}

// IsWithinDuration asserts that a given time.Time is within a certain duration from another time.Time.
func (a *Assert) IsWithinDuration(t1, t2 time.Time, d time.Duration) {
	if abs := t1.Sub(t2); abs > d {
		a.reportError(d, abs, "expected to be within duration")
	}
}

// MatchesPattern asserts that a string matches a certain pattern.
func (a *Assert) MatchesPattern(pattern, s string) {
	if matched, _ := regexp.MatchString(pattern, s); !matched {
		a.reportError(pattern, s, "expected to match pattern")
	}
}

// Panics asserts that a certain function panics.
func (a *Assert) Panics(f func()) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	defer func() {
		if r := recover(); r == nil {
			a.reportError("panic", nil, "expected to panic")
		}
	}()

	f()
}

// PanicsWith asserts that a certain function panics with a specific value.
func (a *Assert) PanicsWith(f func(), expected interface{}) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	defer func() {
		if r := recover(); r != expected {
			a.reportError(expected, r, "expected to panic with")
		}
	}()

	f()
}

// NotPanics asserts that a certain function does not panic.
func (a *Assert) NotPanics(f func()) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	defer func() {
		if r := recover(); r != nil {
			a.reportError("no panic", r, "expected not to panic")
		}
	}()

	f()
}

// SliceDiff asserts that two slices are equal with enhanced diff output for failures.
// Provides detailed context showing which elements differ and their positions.
func (a *Assert) SliceDiff(got, want []int) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	// Check lengths first
	if len(got) != len(want) {
		a.errorMsg = fmt.Sprintf("slices differ in length\n  got: %d\n  want: %d", len(got), len(want))
		return
	}

	// Find first difference
	for i, gotVal := range got {
		if gotVal != want[i] {
			a.errorMsg = fmt.Sprintf("slices differ at index %d\n  got: %d\n  want: %d", i, gotVal, want[i])
			return
		}
	}

	// Slices are identical - no error
}

// SliceDiffGeneric asserts that two slices of any comparable type are equal with enhanced diff output.
// Provides detailed context showing which elements differ and their positions.
func (a *Assert) SliceDiffGeneric(got, want any) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	// Use reflection to handle any slice type
	gotReflect := reflect.ValueOf(got)
	wantReflect := reflect.ValueOf(want)

	// Ensure both are slices
	if gotReflect.Kind() != reflect.Slice {
		a.errorMsg = fmt.Sprintf("got is not a slice: %T", got)
		return
	}
	if wantReflect.Kind() != reflect.Slice {
		a.errorMsg = fmt.Sprintf("want is not a slice: %T", want)
		return
	}

	// Check lengths first
	gotLen := gotReflect.Len()
	wantLen := wantReflect.Len()
	if gotLen != wantLen {
		a.errorMsg = fmt.Sprintf("slices differ in length\n  got: %d\n  want: %d", gotLen, wantLen)
		return
	}

	// Find first difference
	for i := 0; i < gotLen; i++ {
		gotVal := gotReflect.Index(i).Interface()
		wantVal := wantReflect.Index(i).Interface()

		if !reflect.DeepEqual(gotVal, wantVal) {
			a.errorMsg = fmt.Sprintf("slices differ at index %d\n  got: %v\n  want: %v", i, gotVal, wantVal)
			return
		}
	}

	// Slices are identical - no error
}

// MapDiff asserts that two maps are equal with enhanced diff output for failures.
// Provides detailed context showing missing keys, extra keys, and value differences.
func (a *Assert) MapDiff(got, want any) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	// Use reflection to handle any map type
	gotReflect := reflect.ValueOf(got)
	wantReflect := reflect.ValueOf(want)

	// Ensure both are maps
	if gotReflect.Kind() != reflect.Map {
		a.errorMsg = fmt.Sprintf("got is not a map: %T", got)
		return
	}
	if wantReflect.Kind() != reflect.Map {
		a.errorMsg = fmt.Sprintf("want is not a map: %T", want)
		return
	}

	// Check for missing keys (in want but not in got)
	wantKeys := wantReflect.MapKeys()
	for _, wantKey := range wantKeys {
		if !gotReflect.MapIndex(wantKey).IsValid() {
			wantValue := wantReflect.MapIndex(wantKey).Interface()
			a.errorMsg = fmt.Sprintf("maps differ: missing key %q\n  expected value: %v", wantKey.Interface(), wantValue)
			return
		}
	}

	// Check for extra keys (in got but not in want)
	gotKeys := gotReflect.MapKeys()
	for _, gotKey := range gotKeys {
		if !wantReflect.MapIndex(gotKey).IsValid() {
			gotValue := gotReflect.MapIndex(gotKey).Interface()
			a.errorMsg = fmt.Sprintf("maps differ: unexpected key %q\n  got value: %v", gotKey.Interface(), gotValue)
			return
		}
	}

	// Check for value differences
	for _, key := range wantKeys {
		gotValue := gotReflect.MapIndex(key).Interface()
		wantValue := wantReflect.MapIndex(key).Interface()

		if !reflect.DeepEqual(gotValue, wantValue) {
			a.errorMsg = fmt.Sprintf("maps differ at key %q\n  got: %v\n  want: %v", key.Interface(), gotValue, wantValue)
			return
		}
	}

	// Maps are identical - no error
}

// StructDiff asserts that two structs are equal with enhanced diff output for failures.
// Provides detailed context showing which fields differ and their values.
func (a *Assert) StructDiff(got, want any) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	// Use reflection to handle any struct type
	gotReflect := reflect.ValueOf(got)
	wantReflect := reflect.ValueOf(want)

	// Ensure both are structs
	if gotReflect.Kind() != reflect.Struct {
		a.errorMsg = fmt.Sprintf("got is not a struct: %T", got)
		return
	}
	if wantReflect.Kind() != reflect.Struct {
		a.errorMsg = fmt.Sprintf("want is not a struct: %T", want)
		return
	}

	// Ensure same struct type
	gotType := gotReflect.Type()
	wantType := wantReflect.Type()
	if gotType != wantType {
		a.errorMsg = fmt.Sprintf("struct types differ: got %s, want %s", gotType, wantType)
		return
	}

	// Compare each field
	numFields := gotType.NumField()
	for i := 0; i < numFields; i++ {
		field := gotType.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		gotFieldValue := gotReflect.Field(i).Interface()
		wantFieldValue := wantReflect.Field(i).Interface()

		if !reflect.DeepEqual(gotFieldValue, wantFieldValue) {
			a.errorMsg = fmt.Sprintf("structs differ at field %q\n  got: %v\n  want: %v", field.Name, gotFieldValue, wantFieldValue)
			return
		}
	}

	// Structs are identical - no error
}

// DeepDiff asserts that two values of any type are equal with intelligent diff routing.
// Automatically selects the most appropriate diff method based on type:
// - Slices use SliceDiffGeneric for enhanced slice comparison
// - Maps use MapDiff for key/value analysis
// - Structs use StructDiff for field-level comparison
// - Other types use standard deep equality with clear error reporting
func (a *Assert) DeepDiff(got, want any) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	// Quick equality check first
	if reflect.DeepEqual(got, want) {
		return // Values are identical
	}

	// Get reflection values and types
	gotValue := reflect.ValueOf(got)
	wantValue := reflect.ValueOf(want)
	gotType := gotValue.Type()
	wantType := wantValue.Type()

	// Check if types are different
	if gotType != wantType {
		a.errorMsg = fmt.Sprintf("types differ\n  got: %s\n  want: %s", gotType, wantType)
		return
	}

	// Route to specialized diff methods based on type
	switch gotValue.Kind() {
	case reflect.Slice:
		a.SliceDiffGeneric(got, want)
		return
	case reflect.Map:
		a.MapDiff(got, want)
		return
	case reflect.Struct:
		a.StructDiff(got, want)
		return
	default:
		// For primitives and other types, provide basic comparison
		a.errorMsg = fmt.Sprintf("values differ\n  got: %v\n  want: %v", got, want)
		return
	}
}

// Condition asserts that a certain condition is true.
func (a *Assert) Condition(condition bool) {
	if !condition {
		a.reportError(true, condition, "expected condition to be true")
	}
}

// Conditionf asserts that a certain condition is true with a formatted message.
func (a *Assert) Conditionf(condition bool, format string, args ...interface{}) {
	if !condition {
		a.reportError(true, condition, fmt.Sprintf(format, args...))
	}
}

// HttpStatus asserts that a HTTP response has a certain status code.
func (a *Assert) HttpStatus(response *http.Response, expected int) {
	if response.StatusCode != expected {
		a.reportError(expected, response.StatusCode, "expected different HTTP status")
	}
}

// JsonEqual asserts that a JSON string is equivalent to another JSON string or object.
func (a *Assert) JsonEqual(expected, actual string) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	// First check if JSON strings are identical (fast path)
	if expected == actual {
		return
	}

	// Parse both JSON strings
	var obj1, obj2 interface{}
	err1 := json.Unmarshal([]byte(expected), &obj1)
	err2 := json.Unmarshal([]byte(actual), &obj2)

	// Handle parse errors
	if err1 != nil && err2 != nil {
		a.reportStringError(actual, expected, "both JSON strings are invalid")
		return
	}
	if err1 != nil {
		a.reportError(actual, expected, "expected JSON is invalid: "+err1.Error())
		return
	}
	if err2 != nil {
		a.reportError(actual, expected, "actual JSON is invalid: "+err2.Error())
		return
	}

	// Compare parsed objects
	if !reflect.DeepEqual(obj1, obj2) {
		// Objects differ - provide enhanced JSON string comparison
		a.reportStringError(actual, expected, "JSON objects differ")
	}
}

// HasHeader asserts that a HTTP response has a certain header.
func (a *Assert) HasHeader(response *http.Response, header string) {
	if _, ok := response.Header[header]; !ok {
		a.reportError(header, nil, "expected to have header")
	}
}

// HeaderEqual asserts that a HTTP response header has a certain value.
func (a *Assert) HeaderEqual(response *http.Response, header, expected string) {
	if value, ok := response.Header[header]; !ok || len(value) == 0 || value[0] != expected {
		a.reportError(expected, response.Header[header], "expected different header value")
	}
}

// BodyContains asserts that a HTTP response body contains a certain string.
func (a *Assert) BodyContains(response *http.Response, expected string) {
	body, _ := ioutil.ReadAll(response.Body)
	if !strings.Contains(string(body), expected) {
		a.reportError(expected, string(body), "expected body to contain")
	}
}

// BodyMatches asserts that a HTTP response body matches a certain regular expression.
func (a *Assert) BodyMatches(response *http.Response, pattern string) {
	body, _ := ioutil.ReadAll(response.Body)
	if matched, _ := regexp.MatchString(pattern, string(body)); !matched {
		a.reportError(pattern, string(body), "expected body to match")
	}
}

// IsRedirect asserts that a HTTP response is a redirect.
func (a *Assert) IsRedirect(response *http.Response) {
	if response.StatusCode < 300 || response.StatusCode >= 400 {
		a.reportError("3xx", response.StatusCode, "expected redirect status code")
	}
}

// IsSuccess asserts that a HTTP response is a success.
func (a *Assert) IsSuccess(response *http.Response) {
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		a.reportError("2xx", response.StatusCode, "expected success status code")
	}
}

// IsClientError asserts that a HTTP response is a client error.
func (a *Assert) IsClientError(response *http.Response) {
	if response.StatusCode < 400 || response.StatusCode >= 500 {
		a.reportError("4xx", response.StatusCode, "expected client error status code")
	}
}

// IsServerError asserts that a HTTP response is a server error.
func (a *Assert) IsServerError(response *http.Response) {
	if response.StatusCode < 500 || response.StatusCode >= 600 {
		a.reportError("5xx", response.StatusCode, "expected server error status code")
	}
}

// BodyJsonEqual asserts that a HTTP response body is equivalent to a given JSON object.
func (a *Assert) BodyJsonEqual(response *http.Response, expected interface{}) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		a.reportError(nil, expected, "failed to read response body: "+err.Error())
		return
	}

	// Parse actual JSON from response body
	var actual interface{}
	if err := json.Unmarshal(body, &actual); err != nil {
		// JSON parsing failed - show the raw body with enhanced diff if expected is string
		if expectedStr, ok := expected.(string); ok {
			a.reportStringError(string(body), expectedStr, "response body is not valid JSON: "+err.Error())
		} else {
			a.reportError(string(body), expected, "response body is not valid JSON: "+err.Error())
		}
		return
	}

	// Compare parsed objects
	if !reflect.DeepEqual(actual, expected) {
		// If expected is a string (JSON), show enhanced string comparison of raw JSON
		if expectedStr, ok := expected.(string); ok {
			a.reportStringError(string(body), expectedStr, "response JSON differs from expected")
		} else {
			// Expected is an object, use standard comparison
			a.reportError(actual, expected, "response JSON differs from expected object")
		}
	}
}

// HasCookie asserts that a HTTP response has a certain cookie.
func (a *Assert) HasCookie(response *http.Response, name string) {
	var hasCookie bool
	for _, cookie := range response.Cookies() {
		if cookie.Name == name {
			hasCookie = true
			break
		}
	}
	if !hasCookie {
		a.reportError(name, nil, "expected to have cookie")
	}
}

// HeaderContains asserts that a HTTP response header contains a certain value.
func (a *Assert) HeaderContains(response *http.Response, header, expected string) {
	values, ok := response.Header[header]
	if !ok {
		a.reportError(expected, nil, "expected to have header")
	} else {
		var found bool
		for _, value := range values {
			if value == expected {
				found = true
				break
			}
		}
		if !found {
			a.reportError(expected, values, "expected header to contain value")
		}
	}
}

func (a *Assert) ResponseTime(url string, maxTime time.Duration) {
	start := time.Now()
	_, err := http.Get(url)
	if err != nil {
		a.reportError(maxTime, err, "error making request")
	}
	elapsed := time.Since(start)
	if elapsed > maxTime {
		a.reportError(maxTime, elapsed, "response time exceeded")
	}
}

func (a *Assert) IsSorted(slice []int) {
	if !sort.IntsAreSorted(slice) {
		a.reportError(nil, slice, "slice is not sorted")
	}
}

func (a *Assert) IsSortedFloat64(slice []float64) {
	if !sort.Float64sAreSorted(slice) {
		a.reportError(nil, slice, "slice is not sorted")
	}
}

func (a *Assert) FileExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		a.reportError(path, nil, "file does not exist")
	}
}

func (a *Assert) DirectoryExists(path string) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) || !info.IsDir() {
		a.reportError(path, nil, "directory does not exist")
	}
}

// EventuallyConfig holds configuration for Eventually and Never assertions.
type EventuallyConfig struct {
	// Timeout is the maximum duration to wait for the condition.
	Timeout time.Duration
	// Interval is the polling interval between condition checks.
	Interval time.Duration
	// BackoffFactor multiplies the interval after each failed attempt (exponential backoff).
	// Set to 1.0 for constant interval. Must be >= 1.0.
	BackoffFactor float64
	// MaxInterval is the maximum interval between polls when using backoff.
	// If zero, no maximum is enforced.
	MaxInterval time.Duration
}

// defaultEventuallyConfig provides sensible defaults for async testing.
func defaultEventuallyConfig() EventuallyConfig {
	return EventuallyConfig{
		Timeout:       5 * time.Second,
		Interval:      100 * time.Millisecond,
		BackoffFactor: 1.0, // No backoff by default
		MaxInterval:   0,   // No maximum by default
	}
}

// Eventually asserts that a condition becomes true within a timeout period.
// Uses configurable polling with optional exponential backoff.
// Follows GoWise principles of deterministic timing and resource cleanup.
func (a *Assert) Eventually(condition func() bool, timeout, interval time.Duration) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	config := EventuallyConfig{
		Timeout:       timeout,
		Interval:      interval,
		BackoffFactor: 1.0,
		MaxInterval:   0,
	}

	a.eventuallyWithConfig(condition, config)
}

// Never asserts that a condition never becomes true within a timeout period.
// Uses configurable polling with optional exponential backoff.
// Fails immediately if the condition becomes true at any point.
func (a *Assert) Never(condition func() bool, timeout, interval time.Duration) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	config := EventuallyConfig{
		Timeout:       timeout,
		Interval:      interval,
		BackoffFactor: 1.0,
		MaxInterval:   0,
	}

	a.neverWithConfig(condition, config)
}

// EventuallyWith asserts that a condition becomes true using custom configuration.
// Provides fine-grained control over timeout, polling, and backoff behaviour.
func (a *Assert) EventuallyWith(condition func() bool, config EventuallyConfig) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	// Validate and apply defaults
	if config.Timeout <= 0 {
		config.Timeout = defaultEventuallyConfig().Timeout
	}
	if config.Interval <= 0 {
		config.Interval = defaultEventuallyConfig().Interval
	}
	if config.BackoffFactor < 1.0 {
		config.BackoffFactor = 1.0
	}

	a.eventuallyWithConfig(condition, config)
}

// NeverWith asserts that a condition never becomes true using custom configuration.
// Provides fine-grained control over timeout, polling, and backoff behaviour.
func (a *Assert) NeverWith(condition func() bool, config EventuallyConfig) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	// Validate and apply defaults
	if config.Timeout <= 0 {
		config.Timeout = defaultEventuallyConfig().Timeout
	}
	if config.Interval <= 0 {
		config.Interval = defaultEventuallyConfig().Interval
	}
	if config.BackoffFactor < 1.0 {
		config.BackoffFactor = 1.0
	}

	a.neverWithConfig(condition, config)
}

// eventuallyWithConfig implements the core Eventually logic with proper resource management.
func (a *Assert) eventuallyWithConfig(condition func() bool, config EventuallyConfig) {
	// Create context with timeout for clean cancellation
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	// Track timing for error reporting
	startTime := time.Now()
	attempts := 0
	currentInterval := config.Interval

	// First check without delay
	attempts++
	if condition() {
		return // Success on first try
	}

	// Start polling loop
	ticker := time.NewTicker(currentInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Timeout reached - report failure with timing context
			elapsed := time.Since(startTime)
			a.errorMsg = fmt.Sprintf("Eventually: condition not met within timeout\n  timeout: %v\n  elapsed: %v\n  attempts: %d\n  final interval: %v",
				config.Timeout, elapsed, attempts, currentInterval)
			return

		case <-ticker.C:
			attempts++
			if condition() {
				return // Success
			}

			// Apply exponential backoff if configured
			if config.BackoffFactor > 1.0 {
				newInterval := time.Duration(float64(currentInterval) * config.BackoffFactor)

				// Respect maximum interval if set
				if config.MaxInterval > 0 && newInterval > config.MaxInterval {
					newInterval = config.MaxInterval
				}

				if newInterval != currentInterval {
					currentInterval = newInterval
					ticker.Stop()
					ticker = time.NewTicker(currentInterval)
				}
			}
		}
	}
}

// neverWithConfig implements the core Never logic with proper resource management.
func (a *Assert) neverWithConfig(condition func() bool, config EventuallyConfig) {
	// Create context with timeout for clean cancellation
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	// Track timing for error reporting
	startTime := time.Now()
	attempts := 0
	currentInterval := config.Interval

	// First check without delay
	attempts++
	if condition() {
		elapsed := time.Since(startTime)
		a.errorMsg = fmt.Sprintf("Never: condition became true unexpectedly\n  elapsed: %v\n  attempts: %d\n  interval: %v",
			elapsed, attempts, config.Interval)
		return
	}

	// Start polling loop
	ticker := time.NewTicker(currentInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Timeout reached successfully - condition never became true
			return

		case <-ticker.C:
			attempts++
			if condition() {
				elapsed := time.Since(startTime)
				a.errorMsg = fmt.Sprintf("Never: condition became true unexpectedly\n  elapsed: %v\n  attempts: %d\n  final interval: %v",
					elapsed, attempts, currentInterval)
				return
			}

			// Apply exponential backoff if configured
			if config.BackoffFactor > 1.0 {
				newInterval := time.Duration(float64(currentInterval) * config.BackoffFactor)

				// Respect maximum interval if set
				if config.MaxInterval > 0 && newInterval > config.MaxInterval {
					newInterval = config.MaxInterval
				}

				if newInterval != currentInterval {
					currentInterval = newInterval
					ticker.Stop()
					ticker = time.NewTicker(currentInterval)
				}
			}
		}
	}
}

// WithinTimeout asserts that a function completes execution within the specified timeout.
// Uses proper resource management and provides detailed error messages with timing context.
func (a *Assert) WithinTimeout(f func(), timeout time.Duration) {
	if t, ok := a.t.(interface{ Helper() }); ok {
		t.Helper()
	}

	// Create context with timeout for clean cancellation
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Track execution time
	startTime := time.Now()
	done := make(chan bool, 1)

	// Execute function in goroutine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Function panicked - still signal completion
				done <- true
			}
		}()
		f()
		done <- true
	}()

	// Wait for completion or timeout
	select {
	case <-done:
		// Function completed successfully within timeout
		return
	case <-ctx.Done():
		// Timeout exceeded
		elapsed := time.Since(startTime)
		a.errorMsg = fmt.Sprintf("WithinTimeout: function did not complete within timeout\n  timeout: %v\n  elapsed: %v", timeout, elapsed)
		return
	}
}
