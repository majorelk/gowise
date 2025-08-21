# GitHub Issues Backlog

Copy-paste these into GitHub Issues for https://github.com/users/majorelk/projects/5

## [testing,docs] Testing philosophy: behaviour over implementation
**Goal:** Document and enforce behaviour-focused testing across unit, integration, and contract layers.  
**Why:** Keeps tests stable across refactors; encodes intended outcomes not internals.  
**Definition of done:**
- Add section to README referencing CLAUDE.md philosophy
- Add CI note to PR template ("no white-box tests")
- Add 2–3 example tests demonstrating the approach

## [assertions] Assertions v0 (core set, examples)
**Goal:** Implement fast and ergonomic core assertions with examples.  
**Why:** Provide a stable, minimal surface for early adopters.  
**Definition of done:**
- Equal, NotEqual, DeepEqual, Same, Nil, NotNil, True, False, Len, Contains
- Error, NoError, ErrorIs, ErrorAs; Panics, NotPanics, PanicsWith
- InDelta, WithinDuration, Match
- Unit tests (success + failure) and `ExampleXxx` for each

## [internal,diff] Diff helpers for readable failures
**Goal:** First-mismatch printers for strings/slices/maps; compact failure output.  
**Why:** Actionable failures reduce debug time.  
**Definition of done:**
- String window with index caret
- Slice/map index/key diff summary
- Wired into assertion failure messages

## [suite,parallel] Suite lifecycle & parallelism
**Goal:** `BeforeAll/Each`, `AfterAll/Each`, suite-level `Parallel`, deterministic reporting.  
**Why:** Structure and speed without flakiness.  
**Definition of done:**
- Suite API, unit tests, integration tests (parallel)
- Deterministic event ordering
- Examples in `examples/suite`

## [runner,mvp] Thin runner MVP
**Goal:** Wrap `go test -json` with pretty console output and essential flags.  
**Why:** Better UX without replacing Go's toolchain.  
**Definition of done:**
- Flags: -run, -timeout, -failfast, -json, -junit
- Cancel on first failure with -failfast
- JSON passthrough + JUnit XML via `encoding/xml`
- Integration test that boots real `go test` on a fixture package

## [reporting] Reporters: JSON and JUnit XML
**Goal:** Machine-readable outputs for CI systems.  
**Why:** Enables dashboards and PR annotations.  
**Definition of done:**
- JSON stream emitter
- JUnit writer (structs via `encoding/xml`)
- Golden tests under `testdata/`

## [examples] Examples: maths, http, suite
**Goal:** Runnable examples demonstrating assertions, http testing, and lifecycle/parallelism.  
**Why:** Documentation you can execute.  
**Definition of done:**
- `examples/maths/*_test.go`
- `examples/http/*_test.go`
- `examples/suite/*_test.go`

## [ci,quality] CI: stdlib-only guard, race, coverage (soft), fuzz smoke
**Goal:** Lock in quality bar across platforms.  
**Why:** Prevents drift; catches data races early.  
**Definition of done:**
- Workflow with no-deps guard (POSIX + PowerShell)
- `-race` on push/PR; coverage report surfaced
- Optional nightly `-fuzz=Fuzz -fuzztime=10s`

## [snapshots] Snapshot testing (text-only)
**Goal:** Golden file snapshots for text outputs with `UPDATE_GOLDEN=1`.  
**Why:** Safe, stdlib-only snapshotting.  
**Definition of done:**
- Snapshot reader/writer under `testdata/`
- Update switch via env var
- Tests for diffing and update flow

## [concurrency] Goroutine leak probe (best-effort)
**Goal:** Detect unexpected goroutine growth in tests.  
**Why:** Catch resource leaks early.  
**Definition of done:**
- Helper that samples `runtime.NumGoroutine()` before/after
- Grace period and allow-list for known backgrounds
- Document limitations (best-effort)