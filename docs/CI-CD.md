# CI/CD Documentation

This document describes the Continuous Integration and Continuous Deployment setup for pyinit.

## Overview

The CI/CD pipeline ensures code quality, runs comprehensive tests, and automates releases across multiple platforms and package managers.

## Workflows

### 1. CI Workflow (`.github/workflows/ci.yml`)

**Triggers:** 
- Push to any branch
- Pull requests to any branch
- Manual dispatch

**Jobs:**
- **Test Matrix**: Runs tests on Ubuntu, macOS, and Windows with Go 1.21.x, 1.22.x, and 1.23.x
- **Lint**: Static code analysis with golangci-lint
- **Build**: Verifies binary builds on all platforms
- **Integration**: Runs integration test suite
- **Security**: Scans for vulnerabilities with gosec
- **Coverage**: Generates test coverage reports

**Quality Gates:**
- All tests must pass with race detection
- Code must be properly formatted (`go fmt`)
- Static analysis must pass
- Security scan must not find critical issues

### 2. Pre-Release Validation (`.github/workflows/pre-release.yml`)

**Triggers:** 
- Push of version tags (v*.*)
- Manual dispatch with tag input

**Jobs:**
- **Tag Validation**: Ensures proper semantic version format
- **Comprehensive Testing**: Full test matrix across all platforms
- **Build Validation**: Tests release binaries with version embedding
- **Security Audit**: Enhanced security scanning for releases
- **Documentation Check**: Verifies required documentation exists
- **Coverage Requirements**: Enforces minimum 50% test coverage

**Requirements for Release:**
- Tag format: `v1.2.3` or `v1.2.3-beta.1`
- All tests pass on all supported platforms
- Test coverage ≥ 50%
- No critical security issues
- Required documentation present
- Binary builds successfully with correct version info

### 3. Release Workflow (`.github/workflows/release.yml`)

**Triggers:**
- Push of version tags (after pre-release validation)
- Manual release creation

**Dependencies:**
- Waits for pre-release validation to complete successfully
- Fails if pre-release validation failed

**Jobs:**
- **Build Binaries**: Creates release binaries for all platforms with version info
- **GitHub Release**: Creates GitHub release with binaries
- **PyPI Publishing**: Publishes Python wrapper package
- **Homebrew Update**: Triggers Homebrew formula update

## Local Development

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-cover

# Run only unit tests
make test-unit

# Run only integration tests
make test-integration

# Run CI-style tests (with race detection)
make test-ci

# Watch tests during development
make test-watch
```

### Quality Checks

```bash
# Format code
make fmt

# Check formatting
make fmt-check

# Run linter
make lint

# Fix linting issues
make lint-fix

# Run all quality checks
make check
```

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Clean build artifacts
make clean
```

## Branch Protection

All branches should be protected with the following rules:
- Require pull request reviews
- Require status checks to pass (CI workflow)
- Require branches to be up to date before merging
- Include administrators in restrictions

## Release Process

### 1. Prepare Release

1. Update version in relevant files if needed
2. Update CHANGELOG.md (for non-prerelease)
3. Commit changes: `git commit -m "Prepare v1.2.3"`
4. Push to main branch and ensure CI passes

### 2. Create Release

1. Create and push tag: `git tag v1.2.3 && git push origin v1.2.3`
2. Pre-release validation will run automatically
3. If validation passes, release workflow will trigger
4. Release will be created with binaries, PyPI package, and Homebrew update

### 3. Prerelease

For beta/RC releases:
1. Tag with prerelease format: `git tag v1.2.3-beta.1`
2. Same validation process applies
3. GitHub release will be marked as prerelease
4. PyPI package will include prerelease version

## Troubleshooting

### Pre-Release Validation Failures

**Test Failures:**
- Check test output in CI logs
- Run `make test-ci` locally to reproduce
- Fix failing tests and push new commit

**Coverage Too Low:**
- Check coverage report in CI artifacts
- Add tests for uncovered code
- Run `make test-cover` locally to see coverage

**Security Issues:**
- Review gosec report in CI artifacts
- Fix security issues or add exceptions if false positives
- Run `make lint` locally to check

**Build Failures:**
- Check build logs for compilation errors
- Verify all platforms build with `make build-all`
- Fix build issues and push

### Release Failures

**Release Dependencies:**
- Ensure pre-release validation completed successfully
- Check validation workflow logs for failures
- Re-run validation if needed

**PyPI Publishing:**
- Check PyPI credentials and permissions
- Verify package builds correctly
- Check for version conflicts

**Homebrew Update:**
- Check webhook logs
- Verify Homebrew repository access
- Manually update formula if needed

## Security

### Secrets Management

Required secrets for CI/CD:
- `GITHUB_TOKEN`: Automatically provided, used for GitHub API
- `HOMEBREW_UPDATE_TOKEN`: Token for updating Homebrew formula
- PyPI trusted publishing configured for automatic PyPI uploads

### Security Scanning

- **gosec**: Scans for common security issues in Go code
- **Nancy**: Checks for known vulnerabilities in dependencies
- **SARIF Upload**: Security scan results uploaded to GitHub Security tab

## Monitoring

### Success Indicators

- ✅ All CI checks pass on every branch
- ✅ Pre-release validation succeeds before releases
- ✅ Releases complete successfully across all platforms
- ✅ Test coverage maintains or improves over time

### Failure Notifications

- GitHub status checks show failures in PR/branch views
- Workflow run logs provide detailed failure information
- Security issues appear in GitHub Security tab

## Configuration Files

- `.github/workflows/ci.yml`: Main CI configuration
- `.github/workflows/pre-release.yml`: Release validation
- `.github/workflows/release.yml`: Release automation
- `.golangci.yml`: Linter configuration
- `Makefile`: Local development commands
- `coverage.out`: Generated test coverage data