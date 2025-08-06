# pyinit Future Improvements TODO

This document outlines potential improvements and enhancements for the pyinit project.

## 🏗️ Architecture & Code Quality

### Testing Infrastructure
- [ ] Add unit tests for all packages (currently no test files exist)
- [ ] Integration tests for full project generation workflows
- [ ] Template validation tests to ensure all templates render correctly
- [ ] Mock uv commands for reliable CI/CD testing

### Error Handling & Validation
- [ ] Input validation for project names, emails, Python versions
- [ ] Template existence checks before rendering
- [ ] Dependency conflict detection (e.g., incompatible FastAPI + library versions)
- [ ] Better error messages with actionable suggestions

## 🚀 Feature Enhancements

### Additional Project Types
- [ ] Flask web framework (placeholder exists)
- [ ] Django support 
- [ ] Data Science projects (Jupyter, pandas, matplotlib)
- [ ] CLI tools with Click/Typer templates
- [ ] Library projects with proper packaging

### Dependency Management
- [ ] Version pinning options (latest vs stable vs specific)
- [ ] Development vs production dependency separation
- [ ] Optional dependencies (like database drivers for FastAPI)
- [ ] Dependency groups for different use cases (testing, docs, etc.)

### Configuration & Customization
- [ ] Project templates users can define/share
- [ ] Custom template directories
- [ ] Per-project configuration files
- [ ] Environment-specific settings (.env file generation)

## 🛠️ Developer Experience

### CLI Improvements
- [ ] Non-interactive mode with flags (`--name`, `--type`, etc.)
- [ ] Configuration presets for common setups
- [ ] Project update/migration commands
- [ ] Dry-run mode to preview what will be generated
- [ ] Verbose mode for debugging

### Template System
- [ ] Template inheritance for shared components
- [ ] Conditional sections in templates
- [ ] Custom template functions/filters
- [ ] Template validation during build

## 📦 Distribution & Compatibility

### Platform Support
- [ ] Windows compatibility testing
- [ ] Different Python version testing (3.8-3.13)
- [ ] uv alternatives (pip, poetry fallbacks)
- [ ] Docker containerization for consistent environments

### Package Management
- [ ] Homebrew formula maintenance automation
- [ ] PyPI package optimization
- [ ] Version synchronization between Go binary and Python wrapper
- [ ] Checksum verification improvements

## 🔧 Performance & Reliability

### Optimization
- [ ] Template caching for faster subsequent runs
- [ ] Parallel dependency installation
- [ ] Progress indicators for long-running operations
- [ ] Cleanup on failure (partial project cleanup)

### Monitoring
- [ ] Usage analytics (optional, privacy-conscious)
- [ ] Error reporting for debugging
- [ ] Performance metrics collection

## 📚 Documentation & Community

### Documentation
- [ ] API documentation for Go packages
- [ ] Template development guide
- [ ] Contributing guidelines with examples
- [ ] Architecture decision records (ADRs)

### Examples & Tutorials
- [ ] Generated project examples in repository
- [ ] Video tutorials for complex workflows
- [ ] Best practices guide
- [ ] Migration guides from other tools

## 🔐 Security & Maintenance

### Security
- [ ] Input sanitization for all user inputs
- [ ] Template security (prevent arbitrary code execution)
- [ ] Dependency vulnerability scanning
- [ ] SBOM generation for compliance

### Maintenance
- [ ] Automated dependency updates
- [ ] Template freshness monitoring
- [ ] Breaking changes detection and migration
- [ ] Deprecation warnings for old features

## 🎯 Quick Wins (High Impact, Low Effort)

Priority items to tackle first:

- [x] Add `--version` flag to show version info
- [ ] Improve error messages with help text
- [ ] Add `--help` examples for common use cases
- [ ] Template syntax validation during build
- [ ] Basic unit tests for core functions
- [ ] Non-interactive mode with environment variables
- [ ] Project validation after generation
- [ ] Better logging with different levels

## Implementation Notes

- Start with testing infrastructure and non-interactive mode as they'll make development and CI/CD much easier
- Focus on quick wins first to build momentum
- Consider community feedback when prioritizing features
- Maintain backward compatibility when possible

## Recent Progress

### Completed ✅
- [x] Add `--version` flag to show version info
- [x] Update CI script to use proper version embedding with ldflags
- [x] Added Windows build support to release workflow

---

*Generated on: August 6, 2025*
*Last updated: August 6, 2025*
*Status: Planning phase (with recent improvements)*