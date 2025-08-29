# Security Policy

## Supported Versions

We actively support the following versions of GoWise with security updates:

| Version | Supported          | Go Versions    | Status |
| ------- | ------------------ | -------------- | ------ |
| 0.3.x   | :white_check_mark: | 1.21, 1.22, 1.23 | Current |
| 0.2.x   | :white_check_mark: | 1.21, 1.22, 1.23 | Maintenance |
| 0.1.x   | :x:                | 1.21, 1.22, 1.23 | End of life |
| < 0.1   | :x:                | -              | End of life |

### Support Timeline

- **Current version**: Full feature development and security updates
- **Previous minor version**: Security updates and critical bug fixes for 6 months
- **Older versions**: No security updates (please upgrade)

## Security Standards

GoWise follows these security principles:

### Zero External Dependencies
- **Supply Chain Security**: No third-party dependencies reduces attack surface
- **Dependency Scanning**: Not applicable - only Go standard library used
- **Vulnerability Management**: Relies on Go team's security practices for standard library

### Secure Development Practices
- **Input Validation**: All user inputs validated and sanitised
- **Memory Safety**: Go's memory safety features prevent common vulnerabilities
- **Code Review**: All changes reviewed by maintainers before merge
- **Automated Testing**: Comprehensive test suite including security-relevant scenarios

### CI/CD Security
- **GitHub Actions**: Minimal permissions, official actions only
- **Secrets Management**: No secrets stored in repository
- **Build Reproducibility**: Deterministic builds with locked Go versions
- **Supply Chain**: Go module checksums verified

## Reporting a Vulnerability

We take all security vulnerabilities seriously. If you discover a security vulnerability in GoWise, please report it responsibly.

### How to Report

**DO NOT** create a public GitHub issue for security vulnerabilities.

Instead, please:

1. **Email**: Send details to `security@majorelk.dev` (replace with actual contact)
2. **Subject**: Include "SECURITY: GoWise Vulnerability Report"
3. **Encryption**: Use PGP encryption if possible (key available on request)

### What to Include

Please include as much of the following information as possible:

- **Description**: Clear description of the vulnerability
- **Impact**: What could an attacker achieve?
- **Reproduction**: Step-by-step instructions to reproduce the issue
- **Environment**: Go version, GoWise version, operating system
- **Proof of Concept**: Minimal code example demonstrating the vulnerability
- **Suggested Fix**: If you have ideas for remediation

### Example Report Structure

```
Subject: SECURITY: GoWise Vulnerability Report - Potential Code Injection

Description:
The assertion library may be vulnerable to code injection when processing 
user-controlled error messages in [specific component].

Impact:
An attacker could potentially execute arbitrary code during test execution 
if they control assertion input in [specific scenario].

Reproduction:
1. Create a test with the following assertion...
2. Pass the following malicious input...
3. Observe that...

Environment:
- GoWise: v0.3.1
- Go: 1.23.0
- OS: Ubuntu 22.04

Proof of Concept:
[Minimal code example]

Suggested Fix:
Input validation should be added to [specific function] to prevent...
```

## Response Timeline

We aim to respond to security reports according to this timeline:

| Timeline | Action |
|----------|--------|
| 24 hours | Initial response acknowledging receipt |
| 72 hours | Initial assessment and severity classification |
| 1 week   | Detailed analysis and reproduction attempt |
| 2 weeks  | Fix development and testing |
| 3 weeks  | Patch release (for confirmed vulnerabilities) |

### Severity Classification

We use the following severity levels:

#### Critical (CVSS 9.0-10.0)
- Remote code execution
- Complete system compromise
- **Response**: Emergency patch within 48-72 hours

#### High (CVSS 7.0-8.9)
- Significant data exposure
- Privilege escalation
- **Response**: Patch within 1 week

#### Medium (CVSS 4.0-6.9)
- Limited data exposure
- Denial of service
- **Response**: Patch within 2 weeks

#### Low (CVSS 0.1-3.9)
- Information disclosure
- Minor functionality bypass
- **Response**: Patch in next regular release

## Security Advisories

When we release security patches, we will:

1. **GitHub Security Advisory**: Create advisory with CVE if applicable
2. **Release Notes**: Include security information in CHANGELOG.md
3. **Notification**: Announce on relevant channels (GitHub, mailing list if applicable)
4. **Credit**: Acknowledge responsible reporters (with permission)

## Vulnerability Disclosure Policy

We follow responsible disclosure principles:

### Our Commitments
- Acknowledge receipt within 24 hours
- Provide regular status updates
- Credit researchers who report responsibly
- Fix confirmed vulnerabilities promptly
- Communicate transparently about timeline

### We Ask That You
- Give us reasonable time to fix issues before public disclosure
- Avoid accessing or modifying data beyond what's necessary to demonstrate the vulnerability
- Don't perform testing that could harm our users or infrastructure
- Don't violate privacy or disrupt services

## Security Best Practices for Users

When using GoWise in your projects:

### Test Environment Security
- **Isolated Testing**: Run tests in isolated environments
- **No Production Data**: Never use real credentials or sensitive data in tests
- **Clean Up**: Ensure tests clean up temporary files and resources
- **Network Isolation**: Consider network isolation for integration tests

### Code Security
- **Input Validation**: Validate inputs even in test code
- **Secret Management**: Don't hardcode secrets in test files
- **Dependency Management**: Keep Go version updated for security patches
- **Code Review**: Review test code as thoroughly as production code

### CI/CD Security
- **Secure Pipelines**: Use secure CI/CD practices
- **Permission Minimisation**: Grant minimal permissions to test runners  
- **Artifact Security**: Secure test artifacts and reports
- **Environment Variables**: Properly manage test environment variables

## Known Security Considerations

### Testing Framework Specific

#### Test Isolation
- Tests should not affect each other or system state
- Temporary files should be properly cleaned up
- Network resources should be isolated

#### Error Messages
- Error messages may contain sensitive information from assertions
- Consider log sanitisation in production environments
- Be cautious with error message content in public CI logs

#### Reflection Usage
- GoWise uses reflection for deep equality checks
- Reflection operations are memory-safe but may expose internal state
- No known security implications, but worth noting for security audits

### Go Language Security
GoWise inherits Go's security characteristics:
- Memory safety (no buffer overflows, use-after-free, etc.)
- Type safety prevents common injection attacks
- Standard library security maintained by Go team

## Security Updates

We will announce security updates through:

1. **GitHub Releases**: Security releases clearly marked
2. **Security Advisories**: GitHub security advisory system
3. **CHANGELOG.md**: Security fixes documented
4. **Go Security Database**: CVEs published where applicable

## Contact Information

- **Security Email**: `security@majorelk.dev` (update with real contact)
- **PGP Key**: Available on request
- **Response Time**: 24 hours for initial response
- **Time Zone**: UTC+0 (adjust as needed)

## Acknowledgements

We would like to thank the following researchers for responsibly disclosing security vulnerabilities:

- None yet (template for future acknowledgements)

## Security Auditing

- **Last Security Review**: [Date when available]
- **Review Scope**: Full codebase security review
- **External Audits**: None completed yet
- **Automated Scanning**: GitHub security scanning enabled

## Legal

This security policy does not create any legal obligations or rights. It describes our current security practices and may be updated as needed.

For questions about this security policy, please contact us at the email address listed above.