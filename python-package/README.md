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

- **Interactive Setup**: Guides you through creating a new Python project
- **Customizable Project Structure**: Choose between `src` layout or direct layout
- **Project Types**: Supports different project types like `cli`, `web`, `library`, `data-science`
- **Automated Environment Setup**: Automatically creates virtual environment and installs dependencies
- **Pre-configured Tools**: Comes with pre-configured tools for formatting and linting (`ruff`, `pyright`)

## Requirements

- Python 3.9+
- macOS or Linux

## How it Works

This package downloads the appropriate `pyinit` binary for your platform on first use and stores it in `~/.pyinit/bin/`. The binary is verified using SHA256 checksums for security.

## Source Code

The main pyinit tool is written in Go and available at: https://github.com/Pradyothsp/pyinit

## License

MIT License - see the main repository for details.