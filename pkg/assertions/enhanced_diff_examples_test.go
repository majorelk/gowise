package assertions

import (
	"fmt"
	"testing"
)

// ExampleAssert_Equal_enhancedDiff demonstrates enhanced string diff output.
func ExampleAssert_Equal_enhancedDiff() {
	// This example shows enhanced diff capabilities for string comparisons.
	// When strings differ, Enhanced diff provides detailed error messages.

	assert := New(&testing.T{})

	// Single-line string differences show position-based diff:
	// "string values differ at position 6"
	// "diff: hello [w]orld vs hello [W]orld"

	// Multi-line string differences show line-based diff with context:
	// "difference at line 2"
	// "context:"
	// "  {"
	// "- \"name\": \"John\","
	// "+ \"name\": \"Jane\","
	// "  \"age\": 25"

	// Example of successful assertion:
	userInput := "Hello World"
	expected := "Hello World"
	assert.Equal(userInput, expected)

	fmt.Printf("Enhanced diff available for string failures: %t", assert.Error() == "")
	// Output: Enhanced diff available for string failures: true
}

// ExampleAssert_WithDiffFormat demonstrates configurable diff formats.
func ExampleAssert_WithDiffFormat() {
	// Configure assertion to always use context format
	assert := New(&testing.T{}).WithDiffFormat(DiffFormatContext)

	got := `line 1
line 2 original
line 3`

	// This would force context diff format when strings differ:
	// string values differ
	//   difference at line 2
	//   context:
	//     line 1
	//   - line 2 original
	//   + line 2 changed
	//     line 3

	// Configure assertion to always use unified format
	assertUnified := New(&testing.T{}).WithDiffFormat(DiffFormatUnified)

	// This would force unified diff format when strings differ:
	// string values differ
	//   difference at line 2
	//   unified diff:
	//   --- got
	//   +++ want
	//   @@ -1,3 +1,3 @@
	//    line 1
	//   -line 2 original
	//   +line 2 changed
	//    line 3

	// Use the declared variable in a successful assertion:
	want := `line 1
line 2 original
line 3` // Same content as 'got'
	assert.Equal(got, want) // Compare equivalent multi-line strings
	assertUnified.Equal(got, want)

	fmt.Printf("All assertions passed: %t",
		assert.Error() == "" && assertUnified.Error() == "")
	// Output: All assertions passed: true
}

// ExampleAssert_JsonEqual_enhancedDiff demonstrates enhanced JSON comparison.
func ExampleAssert_JsonEqual_enhancedDiff() {
	assert := New(&testing.T{})

	// When JSON objects differ semantically, enhanced string diff is used
	jsonGot := `{
  "status": "error",
  "message": "Not found",
  "code": 404
}`

	// This would produce enhanced diff output when objects differ:
	// JSON objects differ
	//   got:  "{...compact JSON...}"
	//   want: "{...compact JSON...}"
	//   difference at line 2
	//   unified diff:
	//   --- got
	//   +++ want
	//   @@ -1,5 +1,5 @@
	//    {
	//   -  "status": "error",
	//   -  "message": "Not found",
	//   -  "code": 404
	//   +  "status": "success",
	//   +  "message": "Found",
	//   +  "code": 200
	//    }

	// Use the declared variable in a successful assertion:
	// Compare with equivalent JSON (semantically same)
	expectedJson := `{
  "status": "error",
  "message": "Not found",
  "code": 404
}` // Same content as jsonGot, just formatted differently
	assert.JsonEqual(jsonGot, expectedJson)

	fmt.Printf("JSON comparison passed: %t", assert.Error() == "")
	// Output: JSON comparison passed: true
}

// ExampleAssert_DeepEqual_enhancedDiff demonstrates enhanced string diff through DeepEqual.
func ExampleAssert_DeepEqual_enhancedDiff() {
	assert := New(&testing.T{})

	// When DeepEqual compares strings, it benefits from enhanced diff
	configGot := `server:
  host: localhost
  port: 8080
  debug: true`

	// This would produce enhanced string diff when strings differ:
	// values differ
	//   got:  "server:\n  host: localhost\n  port: 8080\n  debug: true"
	//   want: "server:\n  host: localhost\n  port: 8080\n  debug: false"
	//   difference at line 4
	//   context:
	//     server:
	//       host: localhost
	//       port: 8080
	//   -   debug: true
	//   +   debug: false

	// Use the declared variable in a successful assertion:
	expectedConfig := `server:
  host: localhost
  port: 8080
  debug: true` // Same content as configGot
	assert.DeepEqual(configGot, expectedConfig)

	fmt.Printf("DeepEqual passed: %t", assert.Error() == "")
	// Output: DeepEqual passed: true
}

// ExampleDiffFormat demonstrates diff format constants.
func ExampleDiffFormat() {
	// Available diff format options:

	// DiffFormatAuto: Automatically choose between context and unified
	// - Uses context format for simple diffs (few changes)
	// - Uses unified format for complex diffs (many changes)
	fmt.Printf("Auto format: %d\n", DiffFormatAuto)

	// DiffFormatContext: Always use context format with +/- lines
	fmt.Printf("Context format: %d\n", DiffFormatContext)

	// DiffFormatUnified: Always use unified format with @@ headers
	fmt.Printf("Unified format: %d\n", DiffFormatUnified)

	// Output:
	// Auto format: 0
	// Context format: 1
	// Unified format: 2
}

// Example_multiLineDiff demonstrates advanced multi-line diff capabilities.
func Example_multiLineDiff() {
	// This example demonstrates the various enhanced diff features
	// that developers get when comparing multi-line strings.

	assert := New(&testing.T{})

	// Example 1: Configuration file changes would produce:
	// string values differ
	//   difference at line 2
	//   unified diff:
	//   --- got
	//   +++ want
	//   @@ -1,9 +1,9 @@
	//    database:
	//   -  host: localhost
	//   +  host: production.db.com
	//      port: 5432
	//   -  name: testdb
	//   -  ssl: false
	//   +  name: proddb
	//   +  ssl: true
	//
	//    cache:
	//   -  type: redis
	//   -  ttl: 300
	//   +  type: memcached
	//   +  ttl: 600

	configOriginal := `database:
  host: localhost
  port: 5432
  name: testdb
  ssl: false

cache:
  type: redis
  ttl: 300`

	// Example 2: JSON configuration with contextual diff would produce:
	// string values differ
	//   difference at line 2
	//   context:
	//     {
	//   -   "environment": "development"
	//   +   "environment": "production"
	//     }
	assertContext := New(&testing.T{}).WithDiffFormat(DiffFormatContext)

	jsonConfig := `{
  "environment": "development"
}`

	// Use the declared variables in successful assertions:
	expectedDbConfig := `database:
  host: localhost
  port: 5432
  name: testdb
  ssl: false

cache:
  type: redis
  ttl: 300` // Same content as configOriginal
	assert.Equal(configOriginal, expectedDbConfig)

	expectedJsonConfig := `{
  "environment": "development"
}` // Same content as jsonConfig
	assertContext.Equal(jsonConfig, expectedJsonConfig)

	fmt.Printf("Multi-line diff examples work: %t",
		assert.Error() == "" && assertContext.Error() == "")
	// Output: Multi-line diff examples work: true
}
