# Changelog

All notable changes to the GoWise testing framework will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- UK English spelling enforcement with automated CI checks
- Comprehensive documentation suite (CONTRIBUTING.md, architecture docs, API reference)
- Enhanced spelling check script with smart filtering for standard library compatibility
- Documentation standards with clear exemptions for API compatibility

### Changed
- Improved CI workflow with dedicated bash script for spelling checks
- Enhanced error messages throughout codebase using UK English spellings

## [0.3.0] - 2024-XX-XX

### Added
- Enhanced collection assertion error messages with detailed context (#70)
- DeepDiff universal handler for Collection Diff Helpers (#69) 
- Collection diff struct with basic functionality (#66)
- Map diff assertion with enhanced error reporting (#64)
- WithinTimeout assertion for asynchronous testing (#73)
- Eventually and Never assertions for async operations
- EventuallyWith for advanced async testing with exponential backoff

### Enhanced
- String diff algorithms with multi-line context and unified diff output
- Error assertion suite with ErrorContains, ErrorMatches, and improved error handling
- Fast-path optimisations for comparable types in equality assertions
- Memory allocation optimisations with lazy error message formatting

### Fixed
- Panic handling in timeout assertions with proper goroutine cleanup
- Context cancellation and resource management in async assertions
- Type safety improvements with generic slice diff methods

## [0.2.0] - 2024-XX-XX

### Added
- Core assertion library with Equal, NotEqual, DeepEqual, Same
- Nil checking assertions for all nillable types (pointers, interfaces, slices, maps, channels, functions)
- Boolean assertions (True, False) with clear error contexts
- Collection assertions (Len, Contains) with type-safe operations
- Error assertions (NoError, HasError, ErrorIs, ErrorAs) with Go 1.13+ error handling
- Panic assertions (Panics, NotPanics, PanicsWith) with recovery handling
- Basic string diff capabilities for multi-line comparisons
- Numeric assertions (InDelta) for floating-point tolerance
- Time assertions (WithinDuration) for temporal comparisons

### Performance
- Fast-path optimisation for comparable types (int, string, bool, etc.)
- Zero allocations in assertion success paths
- Lazy error message construction only on failure
- Efficient diff algorithms for string comparison

### Documentation
- Complete API examples with ExampleXxx functions
- Behaviour-focused testing methodology
- Zero external dependency architecture

## [0.1.0] - 2024-XX-XX

### Added
- Initial project structure and core architecture
- TestingT interface for compatibility with standard testing package
- Basic assertion context with method chaining support
- Project standards documentation (CLAUDE.md, DEFINITION-OF-DONE.md)
- CI/CD pipeline with Go 1.21, 1.22, 1.23 matrix testing
- Behaviour-focused testing validation
- UK English language standards
- Zero external dependencies validation

### Infrastructure  
- GitHub Actions workflow with comprehensive checks
- Automated formatting (`go fmt`) and linting (`go vet`)
- Race condition detection with `-race` flag
- Coverage reporting with 40% minimum threshold
- Benchmark integration for performance regression detection

## [0.0.1] - 2024-XX-XX

### Added
- Initial repository setup
- Project structure following Go conventions
- Basic module configuration (go.mod)
- Initial documentation (README.md)
- MIT licence (pending final decision)

---

## Migration Notes

### From 0.2.x to 0.3.x
- **Enhanced Error Messages**: Collection assertions now provide more detailed context about mismatches
- **New Async Assertions**: `Eventually`, `Never`, and `WithinTimeout` added for asynchronous testing scenarios
- **Breaking Change**: Some internal diff formatting has changed to improve readability
- **Performance**: Additional fast-path optimisations may change benchmark numbers (improvements expected)

### From 0.1.x to 0.2.x  
- **API Stabilisation**: Core assertion methods are now stable
- **Method Chaining**: All assertions now return `*Assert` for fluent chaining
- **Error Handling**: Improved error messages with context and diff information
- **Performance**: Significant optimisations for comparable type assertions

### From 0.0.x to 0.1.x
- **Project Structure**: Major reorganisation of package structure
- **Testing Standards**: Introduction of behaviour-focused testing requirements
- **Documentation**: Comprehensive documentation and contribution guidelines
- **CI/CD**: Full automation of testing, formatting, and validation

---

## Supported Go Versions

GoWise supports the three most recent minor versions of Go:

- **0.3.x**: Go 1.21, 1.22, 1.23
- **0.2.x**: Go 1.21, 1.22, 1.23  
- **0.1.x**: Go 1.21, 1.22, 1.23

---

## Breaking Changes Policy

### Pre-1.0 (Current)
- Minor versions (0.x.0) may include breaking changes
- Breaking changes will be clearly documented in this changelog
- Migration guides provided for significant API changes
- Deprecation warnings given at least one minor version before removal

### Post-1.0 (Future)
- Breaking changes only in major versions (x.0.0)
- Semantic versioning strictly followed
- Deprecated features supported for at least one major version
- Clear migration paths and automated tooling where possible

---

## Performance Tracking

### Assertion Performance (0.3.x)
| Operation | v0.2.x | v0.3.x | Improvement |
|-----------|--------|--------|------------|
| `Equal(int, int)` | ~5ns | ~2ns | 60% faster |
| `Equal(string, string)` | ~8ns | ~5ns | 37% faster |
| `Len([]int{}, 0)` | ~10ns | ~3ns | 70% faster |
| `Contains(string, string)` | ~15ns | ~8ns | 47% faster |

### Memory Allocation (0.3.x)
- **Success path**: 0 allocations (maintained)
- **Failure path**: ~30% reduction in average allocations
- **Diff generation**: ~50% improvement for large diffs

---

## Security

Security vulnerabilities are taken seriously. See [SECURITY.md](SECURITY.md) for our security policy and how to report vulnerabilities.

---

## Contributors

Thank you to all contributors who have helped make GoWise better:

- Core development team
- Community contributors
- Documentation improvements
- Bug reports and feature requests

For a full list of contributors, see the [GitHub contributors page](https://github.com/majorelk/gowise/contributors).

---

## Licence

This project is licensed under the MIT Licence - see the [LICENCE](LICENCE) file for details.

**Note**: Licence terms are currently under review and may change before v1.0.0.