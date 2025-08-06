# pyinit-cli

An interactive CLI tool to create Python project scaffolds with a customizable structure.

This is a Python wrapper for the [pyinit](https://github.com/Pradyothsp/pyinit) Go binary.

## Installation

```bash
pip install pyinit-cli
```

Or with uv:
```bash
uv add --dev pyinit-cli
```

## Usage

```bash
pyinit
```

The tool will guide you through creating a new Python project with:

- **Interactive Setup**: Guides you through project creation with sensible defaults
- **Multiple Project Types**: Support for basic Python projects and web frameworks (FastAPI)
- **Smart Dependency Management**: Interactive selection of libraries with automatic installation
- **Automated Environment Setup**: Automatically installs dependencies using `uv`
- **Pre-configured Tools**: Comes with pre-configured tools for formatting and linting (`ruff`, `pyright`)

## Requirements

- Python 3.9+
- macOS, Linux, or Windows

## How it Works

This package downloads the appropriate `pyinit` binary for your platform on first use and stores it in `~/.pyinit/bin/`. The binary is verified using SHA256 checksums for security.

## Source Code

The main pyinit tool is written in Go and available at: https://github.com/Pradyothsp/pyinit

## License

MIT License - see the main repository for details.