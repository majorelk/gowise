package diff

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// CollectionDiffResult represents the result of comparing two collections.
type CollectionDiffResult struct {
	HasDiff     bool   // Whether the collections differ
	Summary     string // Human-readable summary of the difference
	Detail      string // Detailed breakdown of differences
	CollectionType string // "slice", "array", "map", "string"
	Truncated   bool   // Whether the display was truncated due to size
}

// CollectionContainsDiff compares a collection and item to determine if the item is missing.
// Returns enhanced diff information for better error messages.
func CollectionContainsDiff(container, item interface{}) CollectionDiffResult {
	if container == nil {
		return CollectionDiffResult{
			HasDiff:     true,
			Summary:     "cannot check containment in nil container",
			Detail:      "",
			CollectionType: "nil",
			Truncated:   false,
		}
	}

	containerValue := reflect.ValueOf(container)
	
	switch containerValue.Kind() {
	case reflect.Slice, reflect.Array:
		return sliceContainsDiff(containerValue, item)
	case reflect.Map:
		return mapContainsDiff(containerValue, item)
	case reflect.String:
		return stringContainsDiff(containerValue.String(), item)
	default:
		return CollectionDiffResult{
			HasDiff:     true,
			Summary:     fmt.Sprintf("unsupported container type: %T", container),
			Detail:      "",
			CollectionType: containerValue.Kind().String(),
			Truncated:   false,
		}
	}
}

// CollectionLenDiff compares the length of a collection against expected length.
// Returns enhanced diff information showing collection contents.
func CollectionLenDiff(container interface{}, expectedLen int) CollectionDiffResult {
	if container == nil {
		return CollectionDiffResult{
			HasDiff:     true,
			Summary:     "cannot get length of nil container",
			Detail:      "",
			CollectionType: "nil",
			Truncated:   false,
		}
	}

	containerValue := reflect.ValueOf(container)
	actualLen := containerValue.Len()

	if actualLen == expectedLen {
		return CollectionDiffResult{
			HasDiff:     false,
			Summary:     "",
			Detail:      "",
			CollectionType: containerValue.Kind().String(),
			Truncated:   false,
		}
	}

	// Generate collection content display
	collectionDisplay, truncated := formatCollectionContent(containerValue, 10) // Show first 10 elements
	
	var summary strings.Builder
	summary.WriteString(fmt.Sprintf("got length: %d, want length: %d", actualLen, expectedLen))
	
	var detail strings.Builder
	if actualLen == 0 {
		detail.WriteString("collection is empty")
	} else {
		detail.WriteString(fmt.Sprintf("collection content: %s", collectionDisplay))
		if truncated {
			detail.WriteString(" ... (showing first 10 elements)")
		}
	}

	return CollectionDiffResult{
		HasDiff:     true,
		Summary:     summary.String(),
		Detail:      detail.String(),
		CollectionType: containerValue.Kind().String(),
		Truncated:   truncated,
	}
}

// sliceContainsDiff handles slice and array containment checking
func sliceContainsDiff(containerValue reflect.Value, item interface{}) CollectionDiffResult {
	// Check if item is contained
	for i := 0; i < containerValue.Len(); i++ {
		if reflect.DeepEqual(containerValue.Index(i).Interface(), item) {
			return CollectionDiffResult{
				HasDiff:     false,
				Summary:     "",
				Detail:      "",
				CollectionType: containerValue.Kind().String(),
				Truncated:   false,
			}
		}
	}

	// Item not found - generate diff
	collectionDisplay, truncated := formatCollectionContent(containerValue, 5)
	
	var summary strings.Builder
	summary.WriteString("expected to contain element")
	
	var detail strings.Builder
	detail.WriteString(fmt.Sprintf("missing from collection: %v\n", item))
	if containerValue.Len() == 0 {
		detail.WriteString("collection is empty")
	} else {
		detail.WriteString(fmt.Sprintf("collection content: %s", collectionDisplay))
		if truncated {
			detail.WriteString(fmt.Sprintf(" ... (showing first 5 of %d elements)", containerValue.Len()))
		}
	}

	return CollectionDiffResult{
		HasDiff:     true,
		Summary:     summary.String(),
		Detail:      detail.String(),
		CollectionType: containerValue.Kind().String(),
		Truncated:   truncated,
	}
}

// mapContainsDiff handles map key containment checking
func mapContainsDiff(containerValue reflect.Value, item interface{}) CollectionDiffResult {
	itemValue := reflect.ValueOf(item)
	
	// Check if key exists
	if containerValue.MapIndex(itemValue).IsValid() {
		return CollectionDiffResult{
			HasDiff:     false,
			Summary:     "",
			Detail:      "",
			CollectionType: "map",
			Truncated:   false,
		}
	}

	// Key not found - generate diff
	keys := containerValue.MapKeys()
	var keyStrings []string
	for _, key := range keys {
		keyStrings = append(keyStrings, fmt.Sprintf("%v", key.Interface()))
	}
	sort.Strings(keyStrings) // Sort for consistent output

	var summary strings.Builder
	summary.WriteString("expected to contain key")
	
	var detail strings.Builder
	detail.WriteString(fmt.Sprintf("missing from map: %v\n", item))
	if len(keys) == 0 {
		detail.WriteString("map is empty")
	} else {
		detail.WriteString("available keys: [")
		detail.WriteString(strings.Join(keyStrings, ", "))
		detail.WriteString("]")
	}

	return CollectionDiffResult{
		HasDiff:     true,
		Summary:     summary.String(),
		Detail:      detail.String(),
		CollectionType: "map",
		Truncated:   false,
	}
}

// stringContainsDiff handles string character/substring containment checking
func stringContainsDiff(container string, item interface{}) CollectionDiffResult {
	var searchString string
	var searchType string
	
	// Handle different item types
	switch v := item.(type) {
	case string:
		searchString = v
		searchType = "substring"
	case rune:
		searchString = string(v)
		searchType = "character"
	case byte:
		searchString = string(v)
		searchType = "character"
	default:
		return CollectionDiffResult{
			HasDiff:     true,
			Summary:     fmt.Sprintf("cannot search for %T in string", item),
			Detail:      fmt.Sprintf("item must be string, rune, or byte for string containers, got: %T", item),
			CollectionType: "string",
			Truncated:   false,
		}
	}

	// Check if contained
	if strings.Contains(container, searchString) {
		return CollectionDiffResult{
			HasDiff:     false,
			Summary:     "",
			Detail:      "",
			CollectionType: "string",
			Truncated:   false,
		}
	}

	// Not found - generate diff
	var summary strings.Builder
	summary.WriteString(fmt.Sprintf("expected to contain %s", searchType))
	
	var detail strings.Builder
	detail.WriteString(fmt.Sprintf("%s %q not found in string\n", searchType, searchString))
	detail.WriteString(fmt.Sprintf("string content: %q", container))

	return CollectionDiffResult{
		HasDiff:     true,
		Summary:     summary.String(),
		Detail:      detail.String(),
		CollectionType: "string",
		Truncated:   false,
	}
}

// formatCollectionContent creates a readable display of collection contents
func formatCollectionContent(containerValue reflect.Value, maxElements int) (string, bool) {
	length := containerValue.Len()
	if length == 0 {
		return "[]", false
	}

	var elements []string
	truncated := length > maxElements
	displayCount := length
	if truncated {
		displayCount = maxElements
	}

	for i := 0; i < displayCount; i++ {
		element := containerValue.Index(i).Interface()
		elements = append(elements, fmt.Sprintf("%v", element))
	}

	var result strings.Builder
	result.WriteString("[")
	result.WriteString(strings.Join(elements, " "))
	result.WriteString("]")

	return result.String(), truncated
}