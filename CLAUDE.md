# GoWise Project Rules

All AI/tooling contributions **must** follow this file. The goals are speed, memory safety, parallelism by default, a complete assertion library — and **no external dependencies** beyond the Go standard library.

---

## Non-negotiables

### Standard library only
- This repository must remain dependency-free beyond the Go standard library.
- Do **not** add third-party modules, scripts, or build steps.
- Allowed packages include only stdlib (e.g. `testing`, `time`, `sync`, `reflect`, `encoding/json`, `encoding/xml`, `flag`, `os/exec`, `io`, `fmt`, `net/http/httptest`, etc.).

### Behaviour-focused tests
Unit, integration, and contract tests must assert **observable behaviour** through public APIs and effects (return values, exported state, filesystem/network I/O via stdlib fakes) — **not** internal implementation.  
Do **not**: peek at unexported fields, count internal calls, rely on goroutine scheduling/timing artefacts, use white-box hooks, or `unsafe` to inspect private state.

### Tests for everything
- Every **exported** function has **unit tests**.
- **Integration tests** cover cross-package behaviour (assertions + suite lifecycle + reporters + runner).
- **Contract tests** validate public interfaces via reusable suites.
- Prefer deterministic synchronisation (channels, WaitGroups) over arbitrary `time.Sleep`.

### Documentation & language
- Use **UK English** in docs and comments (behaviour, initialise, organise, prioritise, optimised, licence).
- Provide examples for **every assertion** in docs (`ExampleXxx`).

### Style & design
- Small, explicit APIs; no global mutable state.
- Allocation-light hot paths; lazy formatting on failure.
- Deterministic output and ordering under parallelism.

---

## Testing philosophy (behaviour, not implementation)

**Unit tests**
- Drive via inputs/outputs and public side-effects.
- Avoid inspecting unexported fields or private helpers.

**Integration tests**
- Exercise behaviour across packages via public surfaces.
- Implementation changes that keep behaviour should not break tests.

**Contract tests**
- Express interface invariants as a reusable suite callable by any implementation.
- Check error semantics and concurrency guarantees with explicit synchronisation.

**Anti-patterns to avoid**
- Asserting exact log strings (unless they are part of the contract).
- Counting function calls or internal goroutines.
- Using `unsafe`/reflection to access private fields.
- Timing-sensitive sleeps instead of explicit signalling.

---

## Pull Request checklist (CI-enforced)

- [ ] `go fmt ./...` produces no diff  
- [ ] `go vet ./...` passes  
- [ ] `go test ./... -race -count=1` passes  
- [ ] **No external modules** guard passes  
      - POSIX:
        ```sh
        test "$(go list -m all | wc -l | tr -d ' ')" -eq 1
        ```
      - Windows PowerShell:
        ```powershell
        if ((go list -m all).Length -ne 1) { Write-Error "External modules detected"; exit 1 }
        ```
- [ ] Tests assert **behaviour, not implementation** (no unexported state access)  
- [ ] Unit tests for each exported function touched/added  
- [ ] Integration tests for cross-package behaviour (where applicable)  
- [ ] Contract tests for any new/changed public interface (where applicable)  
- [ ] Benchmarks for hot paths **or** a short justification if omitted  
- [ ] Doc comments + `ExampleXxx` for all new public APIs  
- [ ] UK English used in docs/comments

---

## Design guidance

### Assertions (initial set)
- Equality: `Equal`, `NotEqual`, `DeepEqual`, `Same` (pointer identity)
- Nilness & truthiness: `Nil`, `NotNil`, `True`, `False`
- Collections & strings: `Len`, `Contains`
- Errors: `Error`, `NoError`, `ErrorIs`, `ErrorAs`
- Panics: `Panics`, `NotPanics`, `PanicsWith`
- Numerics, time: `InDelta`, `WithinDuration`
- Predicate: `Match(func(T) bool)`

**Performance**
- Fast-path `comparable` types; fall back to `reflect.DeepEqual` only when required.
- Construct failure messages **lazily**; keep them short and actionable.
- Avoid reflection and allocations in hot paths.

**Parallelism**
- Suite-level and per-test parallelism.
- Deterministic, ordered reporting (buffer events and sort by (sequence, start time)).

**Reporting**
- Human output: `fmt` (compact and readable).
- Machine outputs: `encoding/json` and `encoding/xml` only (e.g. optional JUnit XML).

**Thin runner approach**
- Wrap `go test -json` via `os/exec`.
- Support flags: `-run`, `-parallel`, `-shuffle`, `-seed`, `-failfast`, `-timeout`, `-json`, `-junit`.
- Cancel remaining processes on first failure when `-failfast` is set.

---

## Required examples for each assertion

Each assertion must include:
- A doc comment describing behaviour.
- At least **one success** and **one failure** test.
- An `Example...` function showing idiomatic usage.

---

## Contract test harness pattern

```go
// ContractKV runs behaviour checks for any KV implementation.
func ContractKV(t *testing.T, name string, newStore func(t *testing.T) KV) {
  t.Run(name+"/put_get", func(t *testing.T) {
    s := newStore(t)
    s.Put("k", "v")
    if got, ok := s.Get("k"); !ok || got != "v" {
      t.Fatalf("behaviour: expected (ok=true, v), got (%v, %q)", ok, got)
    }
  })
  t.Run(name+"/missing", func(t *testing.T) {
    s := newStore(t)
    if _, ok := s.Get("nope"); ok {
      t.Fatalf("behaviour: expected ok=false on missing key")
    }
  })
  t.Run(name+"/concurrency", func(t *testing.T) {
    s := newStore(t)
    var wg sync.WaitGroup
    wg.Add(2)
    go func(){ defer wg.Done(); s.Put("k","v") }()
    go func(){ defer wg.Done(); _, _ = s.Get("k") }()
    wg.Wait() // explicit synchronisation, no sleeps
  })
}
```

Implementations call it from their own tests:
```go
func TestMyStore_Contract(t *testing.T) {
  ContractKV(t, "mystore", func(t *testing.T) KV { return NewMyStore() })
}
```

---

## Helpful stdlib-only snippets

**Fast-path equality & deep equality**
```go
// Equal for comparable types (zero alloc fast path).
func Equal[T comparable](t *testing.T, got, want T) {
  t.Helper()
  if got != want {
    t.Fatalf("Equal: mismatch
  got:  %#v
  want: %#v", got, want)
  }
}

// Deep equality for any (fallback).
func DeepEqual(t *testing.T, got, want any) {
  t.Helper()
  if !reflect.DeepEqual(got, want) {
    t.Fatalf("DeepEqual: values differ
  got:  %#v
  want: %#v", got, want)
  }
}
```

**Eventually helper (behavioural, no arbitrary sleeps)**
```go
func Eventually(t *testing.T, within, every time.Duration, cond func() bool) {
  t.Helper()
  deadline := time.Now().Add(within)
  for {
    if cond() { return }
    if time.Now().After(deadline) {
      t.Fatalf("Eventually: condition not met within %v", within)
    }
    time.Sleep(every)
  }
}
```

**Minimal JUnit XML structs (`encoding/xml`)**
```go
type junitFailure struct {
  Message string `xml:"message,attr"`
  Type    string `xml:"type,attr"`
  Body    string `xml:",chardata"`
}
type junitCase struct {
  Name      string        `xml:"name,attr"`
  Classname string        `xml:"classname,attr"`
  Time      string        `xml:"time,attr"`
  Failure   *junitFailure `xml:"failure,omitempty"`
}
type junitSuite struct {
  XMLName  xml.Name    `xml:"testsuite"`
  Name     string      `xml:"name,attr"`
  Tests    int         `xml:"tests,attr"`
  Failures int         `xml:"failures,attr"`
  Time     string      `xml:"time,attr"`
  Cases    []junitCase `xml:"testcase"`
}
```

---

## Examples demonstrating behaviour-based testing

**Unit (string diff not required to pass — behaviour only)**
```go
func TestSum_Behaviour(t *testing.T) {
  got := Sum([]int{1,2,3})
  if got != 6 {
    t.Fatalf("behaviour: expected 6, got %d", got)
  }
}
```

**Integration (HTTP handler)**
```go
func TestHealthHandler(t *testing.T) {
  srv := httptest.NewServer(http.HandlerFunc(Health))
  t.Cleanup(srv.Close)

  resp, err := http.Get(srv.URL + "/healthz")
  if err != nil { t.Fatal(err) }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    t.Fatalf("behaviour: expected 200 OK, got %d", resp.StatusCode)
  }
  b, _ := io.ReadAll(resp.Body)
  if !bytes.Contains(b, []byte("ok")) {
    t.Fatalf("behaviour: expected body to contain 'ok'")
  }
}
```

**Contract (see ContractKV above)**

---

## Prohibited

- Third-party dependencies, tools, or generated code beyond stdlib.
- Heavy reflection in hot paths; global mutable state.
- White-box or timing-fragile tests.

---

## Glossary (UK spellings)

behaviour · initialise · organise · prioritise · optimised · licence

---

## Licence

MIT (unless otherwise stated in `LICENCE`).

---
