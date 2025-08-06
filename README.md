# pyinit

[![GitHub Release](https://img.shields.io/github/v/release/Pradyothsp/pyinit)](https://github.com/Pradyothsp/pyinit/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/Pradyothsp/pyinit/release.yml)](https://github.com/Pradyothsp/pyinit/actions)
[![PyPI](https://img.shields.io/pypi/v/pyinit-cli)](https://pypi.org/project/pyinit-cli/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

🚀 **Interactive Python Project Scaffolding Tool**

An easy-to-use CLI tool that helps you create well-structured Python projects with modern development tools pre-configured. No more starting from scratch or copying project templates!

## ✨ Why pyinit?

- **Interactive Setup** - Guided project creation with sensible defaults
- **Multiple Project Types** - Support for basic Python projects and web frameworks (FastAPI)
- **Smart Dependency Management** - Interactive selection of libraries with automatic installation via [uv](https://docs.astral.sh/uv/)
- **Modern Tools** - Pre-configured with [ruff](https://docs.astral.sh/ruff/) for lightning-fast linting and formatting, and [pyright](https://github.com/microsoft/pyright) for robust type checking
- **Cross-Platform** - Works on macOS, Linux, and Windows with native binaries
- **Zero Configuration** - Everything works out of the box, no complex setup required

## 📦 Installation

### Homebrew (macOS/Linux)
```bash
brew install Pradyothsp/pyinit/pyinit
```

### pip
```bash
pip install pyinit-cli
```

## 🚀 Quick Start

Simply run:
```bash
pyinit
```

The tool will guide you through:
1. **Basic Information** - Your name, email, and project details
2. **Project Configuration** - Project name, type (basic, web), and description
3. **Framework Selection** - For web projects, choose FastAPI (more coming soon)
4. **Dependency Selection** - Pick libraries to install automatically
5. **Development Environment** - Automated setup with formatting and linting

## 📁 Generated Project Structure

Here's what you get with a basic project:

```
my-awesome-project/
├── .gitignore              # Comprehensive Python .gitignore
├── .python-version         # Python version specification
├── pyproject.toml          # Modern Python project configuration
├── README.md               # Project documentation
├── my_awesome_project/     # Main package directory
│   ├── __init__.py
│   └── main.py             # Entry point with "Hello, World!"
└── scripts/                # Development scripts
    ├── __init__.py
    ├── fmt.py              # Code formatting (ruff)
    └── fmt_check.py        # Linting and type checking
```

## 🔧 Development Commands

After project creation, you can use these commands for development:

```bash
# Check version information
pyinit --version

# Format code and fix issues
uv run fmt

# Check code quality (linting + type checking)
uv run fmt-check
```

## 🆕 What's New in v0.0.6

- **🪟 Windows Support** - Now available for Windows users
- **📋 Interactive Dependencies** - Choose FastAPI libraries during project creation
- **🔍 Version Information** - Use `--version` or `-v` to see detailed build info
- **🚀 Enhanced FastAPI** - Automatic dependency installation with `uv`

## 💻 Requirements

- **Python 3.9+**
- **Platforms**: macOS, Linux, and Windows

## 🤝 Contributing

Interested in contributing? Check out our [Developer Guide](CONTRIBUTING.md) for setup instructions and development workflows.

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

---

**Made with ❤️ for Python developers who value clean project structure and modern tooling.**
