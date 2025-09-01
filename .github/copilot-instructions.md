# Copilot PR Review Instructions — GoWise

> **Scope:** Use these instructions only for **code review during PRs**. Do **not** generate net-new features or large refactors. Keep suggestions small, safe, and directly tied to the changes in the PR.

## Project overview (for context)
**GoWise** is a fast, memory-safe, zero-dependency testing framework and runner for Go. It layers a focused runner and a rich assertion library on top of the Go standard library. The public API aims to be small, explicit, and deterministic.

## Non‑negotiables (apply to every PR)
- **No third‑party dependencies** — standard library only.
- **Performance and safety first** — prefer low/zero‑allocation paths; avoid reflection unless essential and well‑justified.
- **Determinism & parallel‑safety** — tests/suites must behave deterministically and support safe parallel execution.
- **Clear, minimal APIs** — small, explicit surface; strong doc comments.
- **UK English** for docs, comments, and messages.

---

## What Copilot should do in reviews
- **Assess the diff** for correctness, determinism, performance, API clarity, test coverage, and documentation quality.
- **Make concrete, minimal suggestions** (≤ 5 lines) using GitHub’s suggestion blocks where appropriate.
- **Prefer explanations over code**. When code is required, keep it very small and local to the change.
- **Flag risks** (races, allocations, reflection, flaky tests) and ask for evidence (benchmarks, race runs, fuzzing) when useful.
- **Respect existing conventions** (naming, layout, error messages).

## What Copilot should _not_ do
- Propose new dependencies, large rewrites, new features, or broad architectural changes.
- Introduce reflection or generics unless the PR already uses them and the change is clearly safer/simpler.
- Nitpick style that’s already consistent with the repository (formatting is handled by `gofmt`).

---

## Review focus areas & checklists

### 1) Correctness
- Are edge cases handled (zero values, `nil`, empty containers, NaNs, negative lengths)?
- Assertions mark helpers with `t.Helper()` for accurate call sites.
- Comparable fast‑paths (`==`, `!=`) are used before deep checks.
- Avoid undefined behaviour (modifying maps while ranging, shared slice aliasing).
- Panics are only used where intentionally tested (`Panics` assertions).

### 2) Determinism & parallelism
- `t.Parallel()` is only used when no shared mutable state is accessed.
- No hidden global state; suite lifecycle hooks are predictable.
- Randomness (if any) is seeded and controllable; outputs are order‑stable (e.g. map order normalised).

### 3) Performance
- Avoid unnecessary allocations in hot paths (formatting, conversions, `fmt` in loops).
- Reflection is isolated and justified; prefer type switches or constraints when already in use.
- Consider preallocation and zero‑alloc formatting on critical paths.
- Benchmarks exist or are requested when changes affect hot paths.

### 4) API design & stability
- Public API (under `pkg/...`) is coherent, documented, and minimal.
- Names are clear and consistent; exported identifiers have doc comments and examples where appropriate.
- Potential breaking changes are called out explicitly, with migration notes.

### 5) Tests & coverage
- New code has tests; changed behaviour updates existing tests.
- Tests are deterministic (no time sleeps without deadlines; seeded randomness).
- Use table‑driven tests; keep fixtures under `testdata/` when needed.
- Consider fuzz tests for parser/boundary‑heavy code (Go ≥ 1.18).

### 6) Error messages & developer experience
- Messages are concise, actionable, and in UK English.
- When helpful, include _got/want_ and minimal context. Keep diffs readable and bounded.
- CLI and JSON outputs (if touched) remain stable and machine‑readable; exit codes are meaningful.

---

## Comment structure and tone

Use this structure for each review point:

**[LEVEL] Area — short title**  
_File:_ `path/to/file.go:42–57` • _Risk:_ Low/Medium/High  
**Rationale:** Why this matters (reference project principles above).  
**Action:** What to change, add, or verify (tests/benchmarks/flags).

Levels:
- **[BLOCKER]** correctness, data race, determinism, or public API break without migration.
- **[SUGGESTION]** clear improvement with low risk.
- **[NIT]** minor polish; do not block merge alone.
- **[QUESTION]** request for clarification or evidence (bench/race/fuzz).

### Small change example
Keep suggestions ≤ 5 lines and local to the diff.

```suggestion
t.Helper() // Improve failure call site for this assertion helper
```

### Evidence requests
- “Please run: `go test ./... -race -count=1` and paste the output.”
- “This path looks allocation‑heavy. Could you add a micro‑benchmark and share `-benchmem` results?”
- “If reflection is necessary here, can we isolate it and add tests covering comparable and non‑comparable cases?”

---

## Quick commands authors (or CI) should run
- `go fmt ./... && go vet ./...`
- `go test ./... -race -count=1`
- `go test -bench=. -benchmem ./...` (when performance‑sensitive paths change)

---

## Area‑specific notes

### Assertions (`pkg/assertions`)
- Prefer comparable fast‑paths before deep equality.
- Keep failure output compact; show `got`/`want`, types, and _relevant_ context only.
- Avoid allocations on the success path; pay costs only on failure.

### Suites & lifecycle (`pkg/wise`)
- Hooks are minimal and predictable; no global state or hidden ordering.
- Suite/test parallelism is explicit and safe; helpers don’t share mutable state.
- Utilities like timeouts should be opt‑in and clearly scoped.

### Runner / CLI (when touched)
- Stdlib only (`flag`, `testing`, `encoding/json`/`xml`).
- Flags are consistent with `go test` semantics where it makes sense.
- JSON/XML outputs are stable; document any schema change and provide migration notes.

---

## Approval rubric
- **Approve** — no blockers; only nits or optional suggestions remain.
- **Approve with nits** — small fixes suggested; safe to merge either way.
- **Request changes** — one or more blockers remain (correctness, race, determinism, API break, perf regression).

---

## Language & style
- UK English in all comments and messages.
- Be concise, constructive, and specific. Prefer clarity over cleverness.
- Default to the repository’s existing patterns; do not introduce new conventions ad hoc.

---

_Last updated: keep in sync with the README and current conventions._

