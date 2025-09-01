# GoWise

_A fast, memory-safe, zero-dependency testing framework and runner for Go._

> **Status:** Work in progress — the API may change.  
> **Design goals:** speed, memory safety, parallelism by default, full assertion library, zero external libraries.

---

## Why GoWise?

Go’s built-in `testing` package is excellent but intentionally minimal. **GoWise** adds a focused runner and assertion layer inspired by Jest, JUnit and Mocha while **staying within the Go standard library**. No third-party deps, no magic — just pragmatic ergonomics, speed, and clear output.

---

## Features (current & planned)

| Area                         | Description                                                                                     | Status |
|-----------------------------|-------------------------------------------------------------------------------------------------|--------|
| Zero dependencies           | Only the Go standard library.                                                                   | ✅     |
| Core assertions             | Equality, nil checking, boolean assertions with fast-path optimisation.                          | ✅     |
| Collection assertions       | Length, contains, type-safe operations on slices, arrays, maps, strings.                         | ✅     |
| Enhanced string diff        | Multi-line string comparison with context and unified diff output.                               | ✅     |
| Error assertions            | NoError, HasError, ErrorIs, ErrorAs, ErrorContains, ErrorMatches, Panics, NotPanics.           | ✅     |
| Collection diff helpers     | SliceDiff, SliceDiffGeneric, MapDiff with enhanced error reporting.                              | 🚧     |
| CLI runner (thin MVP)       | Basic test runner wrapping `go test -json` with enhanced filtering and reporting.                | 📝     |
| Suite lifecycle             | `BeforeAll/AfterAll`, `BeforeEach/AfterEach`, per-test timeouts using stdlib patterns.           | 📝     |
| Parallel test execution     | Concurrent test execution with deterministic reporting and proper isolation.                     | 📝     |
| Focus & skip helpers        | `wise.Focus(...)`, `wise.Skip(...)` test filtering for development workflows.                    | 📝     |
| Machine-readable output     | JSON and JUnit XML generation for CI/CD system consumption (stdlib only).                        | 📝     |
| Advanced diff features      | Side-by-side visualisation, JSON-aware semantic diff, configurable output formats.              | 📝     |

Legend: ✅ implemented · 🚧 in progress · 📝 planned

---

## Installation

```bash
go get github.com/majorelk/gowise
```

GoWise is a library and a runner. You can either:
1. Use the library with `go test` (current), or
2. Use the GoWise runner (planned) as a small CLI over your test packages.

---

### Quick Start
1. Write a test with GoWise assertions
```go
package maths_test

import (
  "testing"

  "github.com/majorelk/gowise/pkg/assertions"
)

func TestAdd(t *testing.T) {
  assert := assertions.New(t)
  
  got := 2 + 3
  want := 5

  assert.Equal(got, want)
  assert.NotEqual(got, 6)
  assert.True(got > 0)
  assert.False(got < 0)
  assert.InDelta(3.14159, 3.1416, 0.0002)
}

```

Run with standard tooling:
```bash
go test ./...
```

2. Suites & lifecycle (planned library API)
```go
package store_test

import (
  "testing"
  "github.com/majorelk/gowise/wise"
  "github.com/majorelk/gowise/pkg/assertions"
)

func TestStoreSuite(t *testing.T) {
  wise.Suite(t, func(s *wise.SuiteCtx) {
    var db *InMemoryDB

    s.BeforeAll(func() { db = NewInMemoryDB() })
    s.AfterAll(func() { db.Close() })

    s.BeforeEach(func() { db.Reset() })

    s.Test("put/gets a value", func(t *testing.T) {
      assert := assertions.New(t)
      
      db.Put("k", "v")
      got, ok := db.Get("k")
      assert.True(ok)
      assert.Equal(got, "v")
    })

    s.Test("missing key", func(t *testing.T) {
      assert := assertions.New(t)
      
      _, ok := db.Get("nope")
      assert.False(ok)
    })
  })
}
```

> **Planned Feature:** The suite API will mirror familiar patterns from Jest/Mocha (beforeAll, beforeEach) using only stdlib constructs under the hood. This provides familiar lifecycle management whilst maintaining zero external dependencies.

---

### Assertion Library (examples)

```go
assert := assertions.New(t)

// Core equality assertions (fast-path optimised)
assert.Equal(got, want)              // Deep equality with fast-path for comparable types
assert.NotEqual(got, dontWant)       // Inverse equality checking
assert.DeepEqual(got, want)          // Explicit deep equality via reflection
assert.Same(&a, &a)                  // Pointer identity comparison

// Comprehensive nil checking (all nillable types)
assert.Nil(err)                      // Handles interfaces, pointers, slices, maps, channels, functions  
assert.NotNil(value)                 // Non-nil verification with type safety

// Boolean condition assertions
assert.True(condition)               // Boolean true with clear error context
assert.False(condition)              // Boolean false with clear error context

// Collection assertions (type-safe and performant)
assert.Len(container, 3)             // Length verification for strings, slices, arrays, maps, channels
assert.Contains(container, item)     // Membership testing for strings, slices, arrays, maps

// Error assertions (enhanced error handling)
assert.NoError(err)                   // Verify no error occurred
assert.HasError(err)                  // Verify an error occurred
assert.ErrorIs(err, target)           // Error wrapping with errors.Is
assert.ErrorAs(err, &target)          // Error type assertion with errors.As
assert.ErrorContains(err, "text")     // Error message contains substring
assert.ErrorMatches(err, "pattern")   // Error message matches regex

// Collection diff assertions (enhanced failure reporting)
assert.SliceDiff(got, want)           // Integer slice comparison with detailed diff
assert.SliceDiffGeneric(got, want)    // Any slice type with enhanced error context
assert.MapDiff(got, want)             // Map comparison showing missing/extra keys and value differences

// Numeric and misc assertions
assert.InDelta(3.0, 3.001, 0.01)     // Float tolerance
assert.Panics(func() { must() })      // Panic detection
assert.NotPanics(func() { safe() })   // No-panic verification
```
> All assertions allocate minimally and avoid reflection where possible. Where reflection is required (e.g. deep equality), it is kept tight and well-tested.

---

### Parallel Execution (Planned)
**Planned Feature:** Enable parallelism at the suite or test level with deterministic reporting:

```go
// Suite-level parallelism
s.Parallel() // entire suite runs in parallel with others

// Test-level parallelism  
s.Test("fast path", func(t *testing.T) {
  t.Parallel() // individual test runs in parallel
  // test implementation
})

// Focus and skip for development
wise.Focus(t, "critical path")     // only run focused tests
wise.Skip(t, "slow integration")   // skip specific tests
```

**Key Features:**
- **Deterministic output** - parallel execution with ordered reporting
- **Resource isolation** - proper cleanup and state management
- **Development workflow** - focus/skip helpers for debugging
- **Performance monitoring** - execution time tracking and reporting

---

### Integration Testing Patterns
GoWise encourages comprehensive end-to-end testing alongside unit tests:

**Current Capabilities:**
- **Assertions library** works with any `testing.T` including integration tests
- **Enhanced diff output** provides clear feedback for complex data comparisons
- **Error assertions** handle real-world error scenarios and wrapping

**Planned Integration Features:**
- **Suite lifecycle** for setup/teardown of test resources
- **Timeout management** with `wise.WithTimeout(t, 2*time.Second)`
- **Parallel execution** with proper resource isolation
- **Focus/skip helpers** for debugging integration test failures

**Recommended Patterns:**
- Spin up lightweight in-process components (avoid external services)
- Use temporary directories via `t.TempDir()` for file system tests
- Leverage `net/http/httptest` for HTTP service testing
- Run with `-race` flag in CI to detect concurrency issues

---

### CLI Runner (Planned)
The GoWise runner will wrap `go test -json` to provide enhanced output and filtering:

```text
gowise [flags]

  -run REGEX        filter tests by name pattern
  -parallel N       number of test workers (default: GOMAXPROCS)
  -shuffle          randomise test execution order
  -seed INT         random seed for shuffle (implies -shuffle)
  -failfast         stop execution on first test failure
  -timeout DURATION per-test timeout (e.g. 2s, 500ms, 1m)
  -json FILE        generate JSON test report for CI consumption
  -junit FILE       generate JUnit XML report for CI dashboards
  -v                verbose output with detailed progress
  -format FORMAT    output format: auto, compact, verbose (default: auto)
```

**Reporting Output Generation:**
- **JSON reports** for modern CI/CD systems (GitHub Actions, GitLab CI)
- **JUnit XML generation** for enterprise CI consumption (Jenkins, TeamCity)
- **Human-readable progress** with real-time test status
- **Zero external dependencies** - uses stdlib `encoding/json` and `encoding/xml`

The runner parses `go test -json` output and transforms it into enhanced formats whilst maintaining full compatibility with existing Go tooling.
---

### UK English Standards
GoWise enforces UK English spelling throughout documentation and code comments via automated CI checks:

**Automatically Checked Spellings:**
- behaviour (not behavior), initialise (not initialize), organise (not organize)  
- prioritise (not prioritize), optimise (not optimize), analyse (not analyze)
- catalogue (not catalog), licence (not license/US), colour (not color)
- travelled (not traveled), cancelled (not canceled), modelling (not modeling)

**Smart Filtering:**
- **Standard library compatibility**: `context.WithCancel`, `http.Header` etc. are exempt
- **API compatibility**: JSON/XML struct tags may use US spellings for external APIs
- **Existing code tolerance**: Only flags new violations in comments and documentation

**Check Locally:**
```bash
./scripts/check-uk-spelling.sh
```

The spelling checker focuses on comments and documentation where developers have control, while allowing necessary US spellings for Go standard library integration and external API compatibility.

---

### Development Standards & Testing Policy

**Current Testing Approach:**
- **Behaviour-focused testing** - tests assert observable behaviour, not implementation details
- **Unit tests** for every exported function following GoWise principles
- **Integration tests** demonstrate cross-package functionality  
- **Contract tests** validate public interfaces with reusable test suites
- **Zero external dependencies** - all tests use stdlib only

**Planned Testing Infrastructure:**
- **Fuzz testing integration** (Go ≥ 1.18) for parsing and boundary-heavy functions
- **Performance testing** with benchmark integration and regression detection
- **Snapshot testing** for complex output validation
- **CI/CD integration** with JSON and JUnit XML report generation

**Quality Standards:**
- CI runs with `-race`, `-timeout=5m`, and `-vet=all`
- **UK English spelling enforced** via automated CI checks on comments and documentation
- Golden files stored under `testdata/` and reviewed like code
- **No reflection** in hot paths where avoidable
- **UK English** in all documentation and error messages

**Project Structure:**
```
pkg/
  assertions/         # core assertion library (current)
  wise/              # suites, lifecycle, runner (planned)
  internal/          # optimised internals and diff algorithms
cmd/
  gowise/            # CLI runner binary (planned)
examples/
  basic/             # getting started examples  
  integration/       # real-world usage patterns
```
---
### Performance Philosophy

**Current Implementation:**
- **Fast-path optimisations** for comparable types avoiding reflection overhead
- **Thread-safe assertions** using atomic operations for concurrent access
- **Minimal allocations** in assertion hot paths
- **Lazy formatting** - error messages only constructed on failure
- **Efficient diff algorithms** for string and data structure comparison
- **Fail-fast chaining** preserving only first error in assertion chains

**Planned Performance Features:**
- **Benchmark integration** with performance regression detection
- **Memory profiling** assertions for allocation testing
- **Parallel execution** with efficient resource utilisation
- **Streaming output** for large test suites without memory buildup

**Development Standards:**
- Benchmarks live in `_test.go` files (run with `go test -bench=.`)
- Performance regression >5% requires investigation or justification
- Memory allocation tracking for assertion library hot paths
- Profile-guided optimisation for common usage patterns

---
### Versioning

- Pre-1.0: v0.y.z with occasional breaking changes.
- 1.0+: semantic versioning; deprecations are announced before removal.

---

## Contributing

**TO BE DECIDED**

---

## Licence
**TO BE DECIDED**
