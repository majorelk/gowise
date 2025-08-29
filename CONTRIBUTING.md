# Contributing to GoWise

Thank you for your interest in contributing to GoWise! This document provides guidelines for contributing to the project.

## Development Philosophy

GoWise follows these core principles:

- **Zero external dependencies** - Only Go standard library
- **Behaviour-focused testing** - Tests assert observable behaviour, not implementation
- **UK English** throughout documentation and comments
- **Performance conscious** - Fast-path optimisations and minimal allocations
- **Trunk-based development** with short-lived feature branches

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- Understanding of Go testing patterns

### Setting Up Development Environment

1. **Fork and clone the repository:**
   ```bash
   git clone https://github.com/YOUR_USERNAME/gowise.git
   cd gowise
   ```

2. **Verify everything works:**
   ```bash
   go test ./... -race -count=1
   go fmt ./...
   go vet ./...
   ./scripts/check-uk-spelling.sh
   ```

3. **Run benchmarks:**
   ```bash
   go test ./... -bench=. -benchtime=1s
   ```

## Development Workflow

### 1. Test-Driven Development (TDD)

GoWise follows strict TDD practices:

1. **Red**: Write a failing test that describes the desired behaviour
2. **Green**: Write the minimal code to make the test pass
3. **Refactor**: Improve the code whilst keeping tests green

Example workflow:
```bash
# Create feature branch
git checkout -b feature-new-assertion

# Write failing test first
# Implement minimal code to pass
# Refactor and optimise
# Ensure all checks pass

go test ./... -race -count=1
go fmt ./...
go vet ./...
./scripts/check-uk-spelling.sh
```

### 2. Branch Naming

Use descriptive branch names:
- `feature/add-timeout-assertion`
- `fix/memory-leak-in-diff`
- `docs/update-api-reference`
- `refactor/optimise-string-diff`

### 3. Commit Messages

Follow the conventional commit format:

```
type: brief description in present tense

Optional longer description explaining the why, not the what.
Include any breaking changes or migration notes.

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Formatting, no code change
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `perf`: Performance improvement
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

## Code Standards

### 1. Go Code Quality

- **Go idioms**: Follow effective Go practices
- **No external dependencies**: Use only standard library
- **Fast-path optimisations**: Avoid reflection where possible
- **Minimal allocations**: Especially in assertion hot paths
- **Clear naming**: Functions and variables should be self-documenting

### 2. Testing Requirements

**All exported functions must have:**
- Unit tests demonstrating correct usage
- Unit tests demonstrating error conditions
- `ExampleXxx` functions showing idiomatic usage
- Benchmarks for performance-critical paths

**Test characteristics:**
- **Behaviour-focused**: Test observable behaviour, not implementation
- **Deterministic**: No sleeps, use proper synchronisation
- **Isolated**: Tests should not depend on each other
- **Fast**: Unit tests should complete quickly

### 3. Documentation Standards

- **UK English**: behaviour, initialise, optimise, licence, colour
- **Doc comments**: All exported functions, types, and constants
- **Examples**: Demonstrate real-world usage patterns
- **Clear error messages**: Help users understand what went wrong

### 4. Performance Considerations

- **Benchmark everything**: New features need benchmarks
- **Profile guided**: Use profiling to identify bottlenecks
- **Memory conscious**: Track allocations in hot paths
- **Lazy formatting**: Error messages only on failure

## Testing Guidelines

### Unit Tests

```go
func TestEqual_Success(t *testing.T) {
    assert := New(t)
    
    // Test behaviour, not implementation
    assert.Equal(42, 42)
    
    // Assert no error occurred
    if assert.Error() != "" {
        t.Errorf("Expected no error, got: %s", assert.Error())
    }
}

func TestEqual_Failure(t *testing.T) {
    assert := New(&mockT{})
    
    assert.Equal(42, 24)
    
    // Assert appropriate error occurred
    if assert.Error() == "" {
        t.Error("Expected error for unequal values")
    }
}
```

### Integration Tests

```go
func TestAssertionsIntegration(t *testing.T) {
    // Test cross-package behaviour
    assert := New(t)
    
    // Use real testing.T, test actual behaviour
    result := someComplexOperation()
    assert.Equal(result.Status, "success")
    assert.Len(result.Items, 3)
}
```

### Benchmarks

```go
func BenchmarkEqual_FastPath(b *testing.B) {
    b.ReportAllocs()
    assert := New(&mockT{})
    
    for i := 0; i < b.N; i++ {
        assert.Equal(42, 42)
    }
}
```

## Pull Request Process

### 1. Before Submitting

Ensure your PR meets all requirements:

```bash
# Run full test suite
go test ./... -race -count=1

# Check formatting
go fmt ./...

# Check for issues
go vet ./...

# Verify UK English compliance
./scripts/check-uk-spelling.sh

# Run benchmarks
go test ./... -bench=. -benchtime=1s
```

### 2. PR Requirements

- [ ] All tests pass
- [ ] Code formatted with `go fmt`
- [ ] No `go vet` warnings
- [ ] UK English spelling check passes
- [ ] Unit tests for all new exported functions
- [ ] Integration tests for cross-package features
- [ ] Benchmarks for performance-critical code
- [ ] Documentation updated
- [ ] Examples provided for new APIs

### 3. PR Description Template

```markdown
## Summary
Brief description of what this PR does.

## Changes
- Specific change 1
- Specific change 2

## Testing
- Unit tests: [describe what you tested]
- Integration tests: [if applicable]
- Manual testing: [if applicable]

## Performance
- Benchmarks added: [yes/no]
- Performance impact: [none/improved/degraded with justification]

## Breaking Changes
[List any breaking changes and migration path]

## Documentation
- [ ] Updated relevant documentation
- [ ] Added examples for new APIs
- [ ] Updated CHANGELOG.md (if applicable)
```

### 4. Review Process

1. **Automated checks** must pass (CI/CD)
2. **Code review** by maintainers
3. **Performance review** for critical paths
4. **Documentation review** for user-facing changes

## Issue Guidelines

### Reporting Bugs

Use the bug report template:
- **Description**: What happened vs what you expected
- **Reproduction**: Minimal code example
- **Environment**: Go version, OS, GoWise version
- **Context**: What were you trying to achieve?

### Feature Requests

Use the feature request template:
- **Problem**: What problem does this solve?
- **Proposed solution**: How should it work?
- **Alternatives**: What other solutions did you consider?
- **Examples**: Show usage examples

## Release Process

GoWise follows semantic versioning:

- **Major** (1.0.0): Breaking changes
- **Minor** (0.1.0): New features, backwards compatible
- **Patch** (0.0.1): Bug fixes, backwards compatible

### Pre-1.0 Development

- Breaking changes may occur in minor versions
- All breaking changes will be documented in CHANGELOG.md
- Migration guides provided for significant changes

## Community Guidelines

- **Be respectful**: Treat all contributors with respect
- **Be constructive**: Provide helpful feedback
- **Be patient**: Reviews take time
- **Be collaborative**: We're building this together

## Questions?

- **Documentation**: Check README.md and docs/
- **Issues**: Search existing issues first
- **Discussions**: Use GitHub Discussions for general questions
- **Security**: See SECURITY.md for vulnerability reporting

## Recognition

All contributors are recognised in our release notes and contributor list. Thank you for helping make GoWise better!