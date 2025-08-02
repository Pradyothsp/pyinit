# pyinit

🚀 **Interactive Python Project Scaffolding Tool**

An easy-to-use CLI tool that helps you create well-structured Python projects with modern development tools pre-configured. No more starting from scratch or copying project templates!

## ✨ Why pyinit?

- **Interactive Setup** - Guided project creation with sensible defaults
- **Modern Tools** - Pre-configured with [uv](https://docs.astral.sh/uv/) for fast dependency management, [ruff](https://docs.astral.sh/ruff/) for lightning-fast linting and formatting, and [pyright](https://github.com/microsoft/pyright) for robust type checking
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
2. **Project Configuration** - Project name, type, and description
3. **Setup Options** - Python version and development environment

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
# Format code and fix issues
uv run fmt

# Check code quality (linting + type checking)
uv run fmt-check
```

## 🌟 Upcoming Features

- **More Package Managers** - APT and Yarn registry support
- **Enhanced Project Types** - CLI, web, data-science, and library templates
- **Git Integration** - Automatic git initialization and first commit
- **Platform Support** - GitHub and GitLab specific configurations and workflows

## 💻 Requirements

- **Python 3.9+**
- **Platforms**: macOS, Linux (Windows support coming soon)

## 🤝 Contributing

Interested in contributing? Check out our [Developer Guide](CONTRIBUTING.md) for setup instructions and development workflows.

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

---

**Made with ❤️ for Python developers who value clean project structure and modern tooling.**
