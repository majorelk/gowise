# GoWise Examples

This directory contains practical examples demonstrating various GoWise testing patterns and usage scenarios.

## Examples Overview

### 1. [Basic Usage](./basic-usage/)
**File**: `main.go`

Demonstrates fundamental GoWise assertion usage patterns including:
- Core equality assertions (`Equal`, `NotEqual`, `DeepEqual`, `Same`)
- Nil checking (`Nil`, `NotNil`) 
- Boolean assertions (`True`, `False`)
- Collection assertions (`Len`, `Contains`)
- Error assertions (`NoError`, `HasError`, `ErrorContains`, `ErrorMatches`)
- Numeric and time assertions (`InDelta`, `WithinDuration`)
- Method chaining for fluent assertion style

**Run Example:**
```bash
cd examples/basic-usage
go run main.go
```

### 2. [Integration Testing](./integration-testing/)
**File**: `database_test.go`

Shows how to use GoWise for integration testing scenarios:
- Database operation testing with mock database
- HTTP service testing with `httptest`
- Concurrent operation testing with goroutines
- Asynchronous testing with `Eventually`, `Never`, and `WithinTimeout`
- Business service layer testing
- Comprehensive error condition testing

**Run Tests:**
```bash
cd examples/integration-testing
go test -v
```

### 3. [Performance Testing](./performance-testing/)
**File**: `benchmarks_test.go`

Demonstrates performance testing and benchmarking with GoWise:
- Benchmarking core assertion performance
- Measuring memory allocation patterns
- Performance regression testing
- Comparing GoWise vs manual assertion performance
- Memory usage monitoring in tests
- Performance monitoring and profiling techniques

**Run Benchmarks:**
```bash
cd examples/performance-testing
go test -bench=. -benchmem
go test -bench=BenchmarkCoreAssertions -count=5
```

### 4. [Custom Assertions](./custom-assertions/)
**Files**: `domain_assertions.go`, `domain_test.go`

Shows how to create domain-specific custom assertions:
- Extending GoWise with business domain assertions
- Email validation assertions
- HTTP response validation
- JSON validation with key path traversal
- URL validation assertions
- Time and date validation
- String pattern matching
- Business logic validation (age, price, inventory)
- Collection business logic (uniqueness, sorting)
- Comprehensive domain model validation

**Run Tests:**
```bash
cd examples/custom-assertions
go test -v
```

## Usage Patterns

### Method Chaining
All examples demonstrate GoWise's fluent API with method chaining:

```go
assert := assertions.New(t)
assert.Equal(user.ID, 123).
       True(user.Active).
       Contains(user.Email, "@").
       Len(user.Roles, 2)
```

### Error Handling
Examples show proper error assertion patterns:

```go
result, err := operation()
assert.NoError(err).
       Equal(result.Status, "success")

// For expected errors
_, err = riskyOperation()
assert.HasError(err).
       ErrorContains(err, "expected message").
       ErrorIs(err, ErrSpecificType)
```

### Async Testing
Integration examples demonstrate asynchronous testing:

```go
// Wait for condition to become true
assert.Eventually(func() bool {
    return service.IsReady()
}, 5*time.Second, 100*time.Millisecond)

// Assert condition never becomes true
assert.Never(func() bool {
    return service.HasErrors()
}, 2*time.Second, 50*time.Millisecond)

// Function completes within timeout
assert.WithinTimeout(func() {
    service.ProcessData()
}, 1*time.Second)
```

### Custom Domain Assertions
Custom assertion examples show extension patterns:

```go
// Create domain-specific assertion context
domainAssert := NewDomainAssert(t)

// Use domain-specific assertions
domainAssert.IsValidEmail("user@example.com").
             HasEmailDomain("user@example.com", "example.com").
             IsValidUser(user).
             IsValidProduct(product)
```

## Running All Examples

### Individual Examples
```bash
# Basic usage (executable)
cd examples/basic-usage && go run main.go

# Integration tests
cd examples/integration-testing && go test -v

# Performance benchmarks
cd examples/performance-testing && go test -bench=. -benchmem

# Custom assertions
cd examples/custom-assertions && go test -v
```

### All Tests from Root
```bash
# Run all example tests
go test ./examples/... -v

# Run with race detection
go test ./examples/... -race -v

# Run benchmarks for all examples
go test ./examples/... -bench=. -benchmem
```

## Key Learning Points

### 1. Behaviour-Focused Testing
All examples demonstrate testing observable behaviour rather than implementation details:
- Test public API contracts
- Assert on return values and side effects
- Avoid testing internal state or private methods

### 2. Zero External Dependencies
Examples use only Go standard library plus GoWise:
- `net/http/httptest` for HTTP testing
- `time` package for temporal testing
- `sync` package for concurrency testing
- No third-party testing libraries required

### 3. Performance Awareness
Performance examples show:
- How to benchmark assertion performance
- Memory allocation patterns
- Performance regression testing techniques
- Monitoring and profiling approaches

### 4. Extensibility
Custom assertion examples demonstrate:
- How to extend GoWise for domain-specific needs
- Building reusable assertion libraries
- Maintaining consistent error message formatting
- Composing complex validations

## Best Practices Demonstrated

### 1. Test Organization
- Group related assertions with method chaining
- Use subtests for different scenarios
- Separate setup, execution, and assertion phases

### 2. Error Messages
- Provide clear, actionable error messages
- Include context about expected vs actual values
- Use domain-specific terminology in custom assertions

### 3. Performance
- Benchmark critical assertion paths
- Monitor memory allocations
- Test performance regression boundaries

### 4. Async Testing
- Use `Eventually` for waiting on conditions
- Use `Never` for invariant checking
- Use `WithinTimeout` for operation timing
- Prefer explicit synchronisation over sleeps

## Integration with GoWise Features

These examples showcase integration with GoWise's core features:

- **Fast-path optimisations**: Basic usage shows performance benefits
- **Enhanced diff output**: Integration tests demonstrate detailed error messages
- **Method chaining**: All examples use fluent assertion style
- **Async support**: Integration examples use `Eventually`, `Never`, `WithinTimeout`
- **Extensibility**: Custom assertions show extension patterns
- **UK English**: All error messages use UK English spellings

## Contributing Examples

When adding new examples:

1. Follow the established directory structure
2. Include both success and failure scenarios
3. Add comprehensive comments explaining the patterns
4. Use UK English in all documentation and comments
5. Include runnable code with clear output
6. Add appropriate `go test` commands in documentation
7. Follow behaviour-focused testing principles

For more information about contributing, see [CONTRIBUTING.md](../CONTRIBUTING.md).