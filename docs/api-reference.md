# GoWise API Reference

Complete reference for the GoWise assertion library and testing framework.

## Core Assertion API

### Creating an Assertion Context

```go
import "github.com/majorelk/gowise/pkg/assertions"

func TestExample(t *testing.T) {
    assert := assertions.New(t)
    // Use assert for all assertions in this test
}
```

### `func New(t TestingT) *Assert`

Creates a new assertion context for the given test.

**Parameters:**
- `t TestingT`: Any type implementing the TestingT interface (typically `*testing.T`)

**Returns:**
- `*Assert`: New assertion context

**Example:**
```go
func TestBasicUsage(t *testing.T) {
    assert := assertions.New(t)
    assert.Equal(2+2, 4)
}
```

## Equality Assertions

### `func (a *Assert) Equal(got, want interface{}) *Assert`

Asserts that two values are equal using fast-path comparison for comparable types, falling back to deep equality.

**Performance:**
- **Comparable types**: ~2-5ns (int, string, bool, etc.)
- **Complex types**: ~50-100ns (structs, slices, maps via reflection)

**Examples:**
```go
// Basic equality
assert.Equal(42, 42)
assert.Equal("hello", "hello")
assert.Equal(true, true)

// Complex types (uses reflect.DeepEqual)
assert.Equal([]int{1, 2, 3}, []int{1, 2, 3})
assert.Equal(map[string]int{"a": 1}, map[string]int{"a": 1})

// Structs
type Person struct { Name string; Age int }
assert.Equal(Person{"Alice", 30}, Person{"Alice", 30})
```

**Error Output:**
```
Equal: values differ
  got:  42
  want: 24
```

### `func (a *Assert) NotEqual(got, want interface{}) *Assert`

Asserts that two values are not equal.

**Example:**
```go
assert.NotEqual(42, 24)
assert.NotEqual("hello", "world")
```

### `func (a *Assert) DeepEqual(got, want interface{}) *Assert`

Explicitly uses deep equality comparison via reflection.

**Use when:**
- You want to force deep comparison
- Working with complex nested structures
- Comparing types that aren't comparable

**Example:**
```go
complex1 := map[string][]int{"a": {1, 2}, "b": {3, 4}}
complex2 := map[string][]int{"a": {1, 2}, "b": {3, 4}}
assert.DeepEqual(complex1, complex2)
```

### `func (a *Assert) Same(got, want interface{}) *Assert`

Asserts that two pointers point to the same memory address.

**Example:**
```go
x := 42
ptr1 := &x
ptr2 := ptr1
assert.Same(ptr1, ptr2)  // Same pointer

ptr3 := &x  // Different pointer to same value
assert.Same(ptr1, ptr3)  // This would fail
```

## Nil Assertions

### `func (a *Assert) Nil(value interface{}) *Assert`

Asserts that a value is nil. Works with all nillable types.

**Supports:**
- Pointers (`*T`)
- Interfaces
- Slices (`[]T`)
- Maps (`map[K]V`)
- Channels (`chan T`)
- Functions (`func(...)`)

**Examples:**
```go
var ptr *int
assert.Nil(ptr)

var slice []string
assert.Nil(slice)

var m map[string]int
assert.Nil(m)

var err error
assert.Nil(err)
```

### `func (a *Assert) NotNil(value interface{}) *Assert`

Asserts that a value is not nil.

**Example:**
```go
ptr := &someValue
assert.NotNil(ptr)

slice := make([]string, 0)
assert.NotNil(slice)  // Empty but not nil
```

## Boolean Assertions

### `func (a *Assert) True(value bool) *Assert`

Asserts that a boolean expression is true.

**Examples:**
```go
assert.True(2 > 1)
assert.True(len("hello") == 5)
assert.True(user.IsActive())
```

### `func (a *Assert) False(value bool) *Assert`

Asserts that a boolean expression is false.

**Examples:**
```go
assert.False(2 < 1)
assert.False(len("") > 0)
assert.False(user.IsDeleted())
```

## Collection Assertions

### `func (a *Assert) Len(container interface{}, expectedLength int) *Assert`

Asserts that a container has the expected length.

**Supports:**
- Strings
- Slices (`[]T`)
- Arrays (`[N]T`)
- Maps (`map[K]V`)
- Channels (`chan T`)

**Examples:**
```go
assert.Len("hello", 5)
assert.Len([]int{1, 2, 3}, 3)
assert.Len(map[string]int{"a": 1, "b": 2}, 2)
```

**Error Output:**
```
Len: wrong length
  got length:  5
  want length: 3
  collection content: [1 2 3 4 5]
```

### `func (a *Assert) Contains(container, item interface{}) *Assert`

Asserts that a container contains the specified item.

**Supports:**
- **Strings**: substring search
- **Slices/Arrays**: element membership
- **Maps**: key existence

**Examples:**
```go
// String contains substring
assert.Contains("hello world", "world")

// Slice contains element
assert.Contains([]int{1, 2, 3}, 2)

// Map contains key
assert.Contains(map[string]int{"a": 1, "b": 2}, "a")
```

**Error Output:**
```
Contains: element not found
  missing from collection: "orange"
  collection content: ["apple", "banana", "cherry"]
```

## Error Assertions

### `func (a *Assert) NoError(err error) *Assert`

Asserts that no error occurred.

**Example:**
```go
result, err := someOperation()
assert.NoError(err)
// Continue using result safely
```

**Error Output:**
```
NoError: unexpected error occurred
  error: "connection refused"
```

### `func (a *Assert) HasError(err error) *Assert`

Asserts that an error occurred (err is not nil).

**Example:**
```go
_, err := riskyOperation()
assert.HasError(err)
```

### `func (a *Assert) ErrorIs(err, target error) *Assert`

Asserts that an error matches a target error using `errors.Is`.

**Example:**
```go
_, err := os.Open("nonexistent.txt")
assert.ErrorIs(err, os.ErrNotExist)

// With wrapped errors
wrappedErr := fmt.Errorf("failed to read config: %w", os.ErrPermission)
assert.ErrorIs(wrappedErr, os.ErrPermission)
```

### `func (a *Assert) ErrorAs(err error, target interface{}) *Assert`

Asserts that an error can be assigned to a target type using `errors.As`.

**Example:**
```go
var pathErr *os.PathError
_, err := os.Open("nonexistent.txt")
assert.ErrorAs(err, &pathErr)
// pathErr now contains the specific error details
```

### `func (a *Assert) ErrorContains(err error, substring string) *Assert`

Asserts that an error message contains a specific substring.

**Example:**
```go
err := errors.New("database connection failed: timeout")
assert.ErrorContains(err, "connection failed")
assert.ErrorContains(err, "timeout")
```

### `func (a *Assert) ErrorMatches(err error, pattern string) *Assert`

Asserts that an error message matches a regular expression pattern.

**Example:**
```go
err := errors.New("invalid user ID: 12345")
assert.ErrorMatches(err, `invalid user ID: \d+`)
assert.ErrorMatches(err, `^invalid user.*`)
```

## Panic Assertions

### `func (a *Assert) Panics(fn func()) *Assert`

Asserts that a function panics when called.

**Example:**
```go
assert.Panics(func() {
    panic("something went wrong")
})

assert.Panics(func() {
    slice := []int{1, 2, 3}
    _ = slice[10]  // Index out of bounds
})
```

### `func (a *Assert) NotPanics(fn func()) *Assert`

Asserts that a function does not panic when called.

**Example:**
```go
assert.NotPanics(func() {
    result := safeOperation()
    _ = result
})
```

### `func (a *Assert) PanicsWith(fn func(), expected interface{}) *Assert`

Asserts that a function panics with a specific value.

**Example:**
```go
assert.PanicsWith(func() {
    panic("specific error")
}, "specific error")

assert.PanicsWith(func() {
    panic(42)
}, 42)
```

## Advanced Diff Assertions

### `func (a *Assert) SliceDiff(got, want []int) *Assert`

Integer slice comparison with detailed diff output.

**Example:**
```go
got := []int{1, 2, 3, 5}
want := []int{1, 2, 4, 5}
assert.SliceDiff(got, want)
```

**Error Output:**
```
SliceDiff: slices differ at index 2
  got:  [1 2 3 5]
  want: [1 2 4 5]
  diff:
    [0]: 1 ✓
    [1]: 2 ✓  
    [2]: 3 ≠ 4 ❌
    [3]: 5 ✓
```

### `func (a *Assert) SliceDiffGeneric[T comparable](got, want []T) *Assert`

Generic slice comparison for any comparable type.

**Example:**
```go
got := []string{"a", "b", "c"}
want := []string{"a", "x", "c"}
assert.SliceDiffGeneric(got, want)
```

### `func (a *Assert) MapDiff(got, want map[string]int) *Assert`

Map comparison with detailed diff showing missing, extra, and differing values.

**Example:**
```go
got := map[string]int{"a": 1, "b": 2, "c": 3}
want := map[string]int{"a": 1, "b": 5, "d": 4}
assert.MapDiff(got, want)
```

**Error Output:**
```
MapDiff: maps differ
  missing keys: ["d"]
  extra keys: ["c"]  
  different values:
    ["b"]: got 2, want 5
```

## JSON Assertions

### `func (a *Assert) JsonEqual(got, want string) *Assert`

Compares JSON strings semantically, ignoring formatting differences.

**Example:**
```go
got := `{"name":"Alice","age":30}`
want := `{
  "name": "Alice",
  "age": 30
}`
assert.JsonEqual(got, want)  // Passes despite different formatting
```

**Error Output:**
```
JsonEqual: JSON objects differ semantically
  got:  {"name":"Bob","age":25}
  want: {"name":"Alice","age":30}
  differences:
    .name: "Bob" ≠ "Alice"
    .age: 25 ≠ 30
```

## Numeric Assertions

### `func (a *Assert) InDelta(got, want, delta float64) *Assert`

Asserts that two floating-point numbers are within a specified tolerance.

**Example:**
```go
assert.InDelta(3.14159, 3.14160, 0.0001)  // Passes
assert.InDelta(math.Pi, 3.14159, 0.00001) // Passes
```

**Use cases:**
- Floating-point calculations
- Tolerating rounding errors
- Approximate comparisons

## Time Assertions

### `func (a *Assert) WithinDuration(got, want time.Time, tolerance time.Duration) *Assert`

Asserts that two times are within a specified duration of each other.

**Example:**
```go
start := time.Now()
time.Sleep(10 * time.Millisecond)
end := time.Now()

assert.WithinDuration(end, start, 50*time.Millisecond)
```

## Async Assertions

### `func (a *Assert) Eventually(condition func() bool, timeout, interval time.Duration) *Assert`

Repeatedly checks a condition until it becomes true or timeout is reached.

**Parameters:**
- `condition`: Function returning bool to check repeatedly
- `timeout`: Maximum time to wait for condition
- `interval`: Time between condition checks

**Example:**
```go
var result string
go func() {
    time.Sleep(100 * time.Millisecond)
    result = "ready"
}()

// Wait up to 1 second, checking every 50ms
assert.Eventually(func() bool {
    return result == "ready"
}, 1*time.Second, 50*time.Millisecond)
```

**Advanced Example:**
```go
// Database integration test
assert.Eventually(func() bool {
    user, err := db.GetUser("alice")
    return err == nil && user.Status == "active"
}, 5*time.Second, 100*time.Millisecond)
```

### `func (a *Assert) Never(condition func() bool, duration, interval time.Duration) *Assert`

Asserts that a condition never becomes true during the specified duration.

**Example:**
```go
var counter int32
go func() {
    // Increment counter in background
    for i := 0; i < 5; i++ {
        time.Sleep(20 * time.Millisecond)
        atomic.AddInt32(&counter, 1)
    }
}()

// Assert counter never exceeds 10 during test duration
assert.Never(func() bool {
    return atomic.LoadInt32(&counter) > 10
}, 200*time.Millisecond, 25*time.Millisecond)
```

### `func (a *Assert) EventuallyWith(condition func() bool, config EventuallyConfig) *Assert`

Advanced eventually assertion with configurable backoff and timeout behaviour.

**Config Options:**
```go
type EventuallyConfig struct {
    Timeout       time.Duration // Total timeout
    Interval      time.Duration // Initial interval
    BackoffFactor float64      // Exponential backoff multiplier
    MaxInterval   time.Duration // Maximum interval between checks
}
```

**Example:**
```go
config := EventuallyConfig{
    Timeout:       30 * time.Second,
    Interval:      100 * time.Millisecond,
    BackoffFactor: 1.5,  // Exponential backoff
    MaxInterval:   2 * time.Second,
}

assert.EventuallyWith(func() bool {
    return expensiveCheck()
}, config)
```

## Timeout Assertions

### `func (a *Assert) WithinTimeout(fn func(), timeout time.Duration) *Assert`

Asserts that a function completes within the specified timeout.

**Example:**
```go
// Function should complete quickly
assert.WithinTimeout(func() {
    fastOperation()
}, 100*time.Millisecond)

// Function with potential blocking
assert.WithinTimeout(func() {
    result := databaseQuery()
    processResult(result)
}, 5*time.Second)
```

**Error Output:**
```
WithinTimeout: function did not complete within timeout
  timeout: 100ms
  elapsed: ~150ms
```

**Panic Handling:**
- If the function panics, it's considered "completed"
- Panics are recovered and don't crash the test
- Choose this behaviour for timeout testing vs panic testing

## Configuration and Chaining

### Method Chaining

All assertion methods return `*Assert`, enabling fluent chaining:

```go
assert.Equal(user.Name, "Alice").
       True(user.IsActive()).
       Len(user.Permissions, 3).
       Contains(user.Roles, "admin")
```

### Diff Format Configuration

### `func (a *Assert) WithDiffFormat(format DiffFormat) *Assert`

Configure how string differences are displayed.

**Options:**
```go
const (
    DiffFormatAuto    DiffFormat = iota // Auto-select based on content
    DiffFormatContext                   // Show context around changes
    DiffFormatUnified                   // Unified diff format
)
```

**Example:**
```go
// Force context diff format
assert := New(t).WithDiffFormat(DiffFormatContext)
assert.Equal(longText1, longText2)

// Force unified diff format  
assertUnified := New(t).WithDiffFormat(DiffFormatUnified)
assertUnified.Equal(config1, config2)
```

## Error Handling and Reporting

### `func (a *Assert) Error() string`

Returns the accumulated error message from failed assertions.

**Example:**
```go
assert := New(&mockT{})  // Don't fail immediately
assert.Equal(1, 2)
assert.Equal("a", "b")

if assert.Error() != "" {
    t.Logf("Assertion errors: %s", assert.Error())
}
```

### `func (a *Assert) Failed() bool`

Returns whether any assertions have failed.

**Example:**
```go
assert := New(&mockT{})
assert.Equal(1, 1)  // Passes
assert.Equal(1, 2)  // Fails

if assert.Failed() {
    // Handle multiple failures
    fmt.Printf("Some assertions failed: %s", assert.Error())
}
```

## Custom Extensions

### TestingT Interface

GoWise assertions work with any type implementing `TestingT`:

```go
type TestingT interface {
    Errorf(format string, args ...interface{})
    FailNow()
}
```

**Custom Implementation:**
```go
type CustomT struct {
    errors []string
}

func (c *CustomT) Errorf(format string, args ...interface{}) {
    c.errors = append(c.errors, fmt.Sprintf(format, args...))
}

func (c *CustomT) FailNow() {
    panic("test failed")
}

// Usage
customT := &CustomT{}
assert := New(customT)
assert.Equal(1, 2)
// Check customT.errors for failure details
```

### Domain-Specific Assertions

Extend the `Assert` type with custom methods:

```go
// Add custom assertion methods
func (a *Assert) IsValidEmail(email string) *Assert {
    emailRegex := regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)
    if !emailRegex.MatchString(email) {
        a.t.Errorf("IsValidEmail: invalid email format: %s", email)
    }
    return a
}

func (a *Assert) HasStatus(response *http.Response, expectedStatus int) *Assert {
    if response.StatusCode != expectedStatus {
        a.t.Errorf("HasStatus: wrong status code\n  got: %d\n  want: %d", 
                   response.StatusCode, expectedStatus)
    }
    return a
}

// Usage
assert.IsValidEmail("user@example.com")
assert.HasStatus(response, http.StatusOK)
```

## Performance Considerations

### Fast Path vs Reflection

| Assertion | Fast Path Types | Reflection Types | Performance Difference |
|-----------|----------------|------------------|----------------------|
| `Equal` | `comparable` types | Structs, slices, maps | ~25x faster |
| `Len` | Strings, arrays | Slices, maps, channels | ~10x faster |
| `Contains` | Strings | Slices, maps | ~15x faster |

### Memory Allocation

- **Success path**: Zero allocations
- **Failure path**: Minimal allocations for error formatting
- **Diff generation**: Allocations proportional to difference size

### Best Practices

1. **Use specific assertions**: `True(x > 0)` vs `Equal(x > 0, true)`
2. **Fast path when possible**: Prefer comparable types
3. **Chain related assertions**: Reduces context switching
4. **Avoid complex diff**: For performance-critical tests

## Migration Guide

### From testify/assert

```go
// testify
assert.Equal(t, expected, actual)
assert.NotNil(t, value)
assert.Contains(t, slice, element)

// GoWise
assert := assertions.New(t)
assert.Equal(actual, expected)  // Note: order reversed
assert.NotNil(value)
assert.Contains(slice, element)
```

### From standard testing

```go
// Standard testing
if got != want {
    t.Errorf("got %v, want %v", got, want)
}

// GoWise
assert := assertions.New(t)
assert.Equal(got, want)  // Automatic detailed error messages
```

## Examples and Patterns

### Basic Test Structure

```go
func TestUserService(t *testing.T) {
    assert := assertions.New(t)
    
    service := NewUserService()
    user, err := service.CreateUser("alice", "alice@example.com")
    
    assert.NoError(err).
           NotNil(user).
           Equal(user.Username, "alice").
           True(user.IsActive())
}
```

### Integration Test Pattern

```go
func TestDatabaseIntegration(t *testing.T) {
    assert := assertions.New(t)
    db := setupTestDB(t)
    defer db.Close()
    
    // Test data insertion
    err := db.InsertUser("bob", 25)
    assert.NoError(err)
    
    // Test data retrieval with eventually for consistency
    assert.Eventually(func() bool {
        user, err := db.GetUser("bob")
        return err == nil && user.Age == 25
    }, 2*time.Second, 100*time.Millisecond)
}
```

### Error Testing Pattern

```go
func TestErrorHandling(t *testing.T) {
    assert := assertions.New(t)
    
    service := NewService()
    
    // Test specific error conditions
    _, err := service.ProcessInvalidData(nil)
    assert.HasError(err).
           ErrorContains(err, "invalid data").
           ErrorIs(err, ErrInvalidInput)
           
    // Test error wrapping
    var validationErr *ValidationError
    assert.ErrorAs(err, &validationErr)
}
```

This API reference covers the complete GoWise assertion library. For more examples and patterns, see the [examples directory](../examples/) and [CONTRIBUTING.md](../CONTRIBUTING.md).