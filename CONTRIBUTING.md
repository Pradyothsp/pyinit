# Contributing to pyinit

Thank you for your interest in contributing to pyinit! This guide will help you get set up for development.

## 🏗️ Architecture Overview

pyinit is built with a hybrid architecture:

- **Core Engine**: Written in Go for performance and cross-platform compatibility
- **Python Wrapper**: Provides seamless installation via pip and integrates with Python workflows
- **Template System**: Uses Pongo2 templates for flexible project generation
- **Distribution**: Available via Homebrew (Go binary) and PyPI (Python wrapper)

## 🚀 Development Setup

### Prerequisites

- **Go 1.24+**
- **Python 3.9+** 
- **uv** (recommended for Python package development)

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/Pradyothsp/pyinit.git
   cd pyinit
   ```

2. **Set up Go development**
   ```bash
   go mod download
   
   # Build with version info (recommended)
   make build
   
   # Or build directly
   go build -o pyinit ./cmd/pyinit
   ```

3. **Set up Python package development**
   ```bash
   cd python-package
   pip install -e .
   
   # Or with uv (recommended)
   uv sync
   ```

4. **Test your changes**
   ```bash
   # Test Go binary directly
   ./pyinit
   
   # Test Python wrapper
   cd python-package
   pyinit
   ```

## 📁 Project Structure

```
├── cmd/pyinit/             # Go CLI entry point
├── internal/               # Go internal packages
│   ├── config/            # Project configuration
│   ├── generator/         # Project generation logic
│   ├── prompts/           # User interaction
│   └── setup/             # Environment setup
├── pkg/                   # Go public packages
│   ├── template/          # Template engine
│   └── ui/                # User interface components
├── templates/             # Project templates (*.j2 files)
├── python-package/        # Python wrapper
│   ├── pyinit_cli/       # Python package
│   └── scripts/          # Python development scripts
├── Formula/               # Homebrew formula
└── .github/workflows/     # CI/CD pipelines
```

## 🧪 Testing

### Go Testing
```bash
go test ./...
```

### Python Testing
```bash
cd python-package
python -m pytest
```

### Integration Testing
```bash
# Test full workflow
./pyinit
# Follow the prompts and verify generated project
```

## 📦 Release Process

### Automated Release (Recommended)
1. Create and push a new tag:
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

2. GitHub Actions will automatically:
   - Build Go binaries for all platforms
   - Update Homebrew formula
   - Publish Python package to PyPI
   - Create GitHub release

### Manual Release
See `.github/workflows/release.yml` for the complete automation pipeline.

## 🎯 Areas for Contribution

### High Priority
- **Project Templates**: Add new project types (CLI, web, data-science, library)
- **Git Integration**: Automatic git initialization and first commit
- **Platform Support**: Windows compatibility
- **CI Templates**: GitHub Actions and GitLab CI configurations

### Medium Priority
- **Enhanced Prompts**: Better validation and user experience
- **Configuration**: User preferences and project defaults
- **Documentation**: Examples, tutorials, and guides

### Low Priority
- **Plugins**: Extension system for custom templates
- **Themes**: Different project structure themes
- **Integration**: IDE-specific configurations

## 🐛 Bug Reports

When reporting bugs, please include:
- Operating system and version
- Python version
- Go version (if building from source)
- Complete error messages
- Steps to reproduce

## 💡 Feature Requests

We love feature ideas! Please:
- Check existing issues first
- Describe the use case clearly
- Explain why it would benefit other users
- Consider if it fits the project's scope

## 📝 Code Style

### Go
- Follow standard Go formatting (`gofmt`)
- Use descriptive variable names
- Add comments for public functions
- Keep functions focused and small

### Python
- Follow PEP 8
- Use type hints where appropriate
- Add docstrings for public functions
- Use `ruff` for formatting and linting

## 🔍 Code Review Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Test thoroughly
5. Commit with clear messages
6. Push to your fork
7. Open a Pull Request

## 📞 Getting Help

- **GitHub Issues**: Bug reports and feature requests
- **Discussions**: Questions and general chat
- **Email**: contact@pradyoth-sp.me (for private matters)

## 📄 License

By contributing, you agree that your contributions will be licensed under the MIT License.
