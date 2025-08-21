# Definition of Done Template

Use this template for all GitHub issues to ensure consistent quality standards.

## Standard Definition of Done (All Issues)

**Code Quality:**
- [ ] Code follows Go idioms and project style
- [ ] All exported functions have doc comments (UK English)
- [ ] `go fmt ./...` passes
- [ ] `go vet ./...` passes  
- [ ] No external dependencies beyond Go stdlib

**Testing Requirements:**
- [ ] Unit tests for all exported functions (behaviour-focused)
- [ ] Integration tests where cross-package interaction occurs
- [ ] All tests pass with `go test ./... -race -count=1`
- [ ] Test coverage maintained (minimum 40% overall)
- [ ] No unsafe usage or reflection to access unexported fields in tests

**Documentation:**
- [ ] `ExampleXxx` functions for all new public APIs
- [ ] UK English spelling throughout (behaviour, initialise, licence)
- [ ] Code comments explain "why" not "what"

**CI/CD:**
- [ ] All GitHub Actions workflows pass
- [ ] No external modules guard passes
- [ ] Behaviour-focused testing validation passes

## Feature-Specific Additions

### For Assertion Functions:
- [ ] Success test case demonstrating correct usage
- [ ] Failure test case demonstrating error conditions  
- [ ] Clear, actionable error messages on assertion failure
- [ ] Fast path for comparable types where applicable
- [ ] Lazy error message formatting (only on failure)

### For Suite/Runner Features:
- [ ] Deterministic behaviour under parallelism
- [ ] Proper cleanup in failure scenarios
- [ ] Integration test with real `go test` process

### For Reporting Features:
- [ ] Golden file tests in `testdata/` directory
- [ ] JSON output validates against expected schema
- [ ] Works correctly with `encoding/json` and `encoding/xml` only

### For Examples:
- [ ] Runnable with `go test`
- [ ] Demonstrates realistic use cases
- [ ] Self-contained (no external setup required)
- [ ] Clear comments explaining the approach

## Quality Gates

**Before PR Creation:**
- [ ] All Definition of Done items completed
- [ ] Self-review performed
- [ ] Breaking changes documented

**PR Review Checklist:**
- [ ] Code review by maintainer
- [ ] All CI checks green
- [ ] No regression in existing functionality
- [ ] Documentation updated if needed

## Notes

- This is a living document - update as the project evolves
- For work-in-progress features, clearly mark incomplete items
- When in doubt, over-document rather than under-document
- Remember: behaviour over implementation in all testing