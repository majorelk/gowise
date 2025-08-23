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
| Assertions                  | Rich, fluent assertions with readable diffs.                                                     | ✅     |
| Parallel test execution     | Run tests concurrently with sensible scheduling and isolation.                                   | 🚧     |
| Test lifecycle              | `BeforeAll/AfterAll`, `BeforeEach/AfterEach`, per-test timeouts.                                 | 🚧     |
| Focus & skip                | `wise.Focus(...)`, `wise.Skip(...)` helpers.                                                     | 🚧     |
| CLI runner                  | Filter by pattern, shuffle, seed, fail-fast, JSON/pretty output.                                | 🚧     |
| Reporting                   | Text and machine-readable (JSON) reports; JUnit XML _optional_ (still standard library only).   | 📝     |
| Benchmarks                  | Micro-bench helpers integrated with the runner.                                                  | 📝     |
| Fuzzing                     | First-class wrapper for Go 1.18+ fuzzing (stdlib).                                               | 📝     |

Legend: ✅ implemented · 🚧 in progress · 📝 planned

---

## Installation

```bash
go get github.com/majorelk/gowise
```

GoWise is a libbrary and a runner. You can either:
1. Use the library with `go test`, or
2. Use the GoWise runner (planned) as a small cli over your test packages.

---

### Quick Start
1. Write a test with GoWise assertions
```go
package maths_test

import (
  "testing"

  "github.com/majorelk/gowise/assert"
)

func TestAdd(t *testing.T) {
  got := 2 + 3
  want := 5

  assert.Equal(t, got, want)
  assert.NotEqual(t, got, 6)
  assert.True(t, got > 0)
  assert.False(t, got < 0)
  assert.InDelta(t, 3.14159, 3.1416, 0.0002)
}

```

Run with standard tooling:
```bash
go test ./...
```

2. Suites & lifecycle (library API)
```go
package store_test

import (
  "testing"
  "github.com/majorelk/gowise/wise"
  "github.com/majorelk/gowise/assert"
)

func TestStoreSuite(t *testing.T) {
  wise.Suite(t, func(s *wise.SuiteCtx) {
    var db *InMemoryDB

    s.BeforeAll(func() { db = NewInMemoryDB() })
    s.AfterAll(func() { db.Close() })

    s.BeforeEach(func() { db.Reset() })

    s.Test("put/gets a value", func(t *testing.T) {
      db.Put("k", "v")
      got, ok := db.Get("k")
      assert.True(t, ok)
      assert.Equal(t, got, "v")
    })

    s.Test("missing key", func(t *testing.T) {
      _, ok := db.Get("nope")
      assert.False(t, ok)
    })
  })
}
```

> The suite API mirrors familiar patterns from Jest/Mocha (beforeAll, beforeEach) using only stdlib constructs under the hood.

---

### Assertion Library (examples)

```go
assert.Equal(t, got, want)              // Deep equality (structs, maps, slices)
assert.NotEqual(t, got, dontWant)

assert.Same(t, &a, &b)                  // Pointer identity
assert.Nil(t, err)                      // err == nil
assert.NotNil(t, v)

assert.True(t, cond)
assert.False(t, cond)

assert.InDelta(t, 3.0, 3.001, 0.01)     // Float tolerance
assert.Contains(t, "hello world", "world")
assert.Len(t, mySlice, 3)

assert.Panics(t, func() { must() })
assert.NotPanics(t, func() { safe() })

type User struct{ ID int; Name string }
assert.Match(t, User{1, "A"}, func(u User) bool { return u.ID > 0 && u.Name != "" })
```
> All assertions allocate minimally and avoid reflection where possible. Where reflection is required (e.g. deep equality), it is kept tight and well-tested.

---

### Parralel Exection
Enable parallelism at the suite or test level:
```go
s.Parallel() // entire suite
s.Test("fast path", func(t *testing.T) {
  t.Parallel() // or per test
  // ...
})

```
The runner will ensure deterministic seeding and ordered reporting when running in parallel.

---

### Integration Tests
GoWise encourages end-to-end tests alongside units:

- Spin up lightweight in-process components (no external services).
- Use temporary directories via t.TempDir().
- Set timeouts per test: wise.WithTimeout(t, 2*time.Second).
- Recommend `-race` for CI by default.

---

### CLI Runner (planned)
```text
gowise [flags]

  -run REGEX        filter tests by name
  -parallel N       number of workers (default: GOMAXPROCS)
  -shuffle          randomise test order
  -seed INT         random seed (implies -shuffle)
  -failfast         stop on first failure
  -timeout DURATION per-test timeout (e.g. 2s, 500ms, 1m)
  -json             emit JSON report
  -junit FILE       emit JUnit XML (stdlib encoding/xml)
  -v                verbose
```
---

### Test Policy & Coverage
- Unit tests for every exported function (and the majority of internal helpers).
- Integration tests for the runner, suites and reporters.
- Fuzz tests (Go ≥ 1.18) for parsing and boundary-heavy functions.
- CI runs with -race, -timeout=5m, and -vet=all.
- Golden files are stored under testdata/ and are reviewed like code.

Suggested layout:
```
pkg/
  assert/
  wise/          # suites, lifecycle, runner plumbing
  internal/      # small, well-documented internals
examples/
  maths/
  http/
```
---
### Performance

- Avoid allocations in hot paths (assertions, reporting).
- Benchmarks live in _test.go files (go test -bench=.).
- If a change regresses a benchmark by >5%, it should be investigated or justified in the PR.

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
