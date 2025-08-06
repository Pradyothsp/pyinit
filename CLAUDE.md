# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

pyinit is a hybrid architecture CLI tool for Python project scaffolding:
- **Core Engine**: Go-based CLI using Cobra framework for cross-platform compatibility
- **Python Wrapper**: Python package for seamless pip installation and Python ecosystem integration  
- **Template System**: Pongo2 template engine for flexible project generation
- **Distribution**: Available via Homebrew (Go binary) and PyPI (Python wrapper)

## Common Development Commands

### Go Development
```bash
# Build the Go binary
go build -o pyinit ./cmd/pyinit

# Run tests
go test ./...

# Install dependencies
go mod download
```

### Python Package Development
```bash
# Install in development mode
cd python-package && pip install -e .

# Format Python code
uv run fmt

# Lint and type check Python code
uv run fmt-check
```

### Testing
```bash
# Test Go binary directly
./pyinit

# Test Python wrapper
cd python-package && pyinit

# Integration test - run and verify generated project
./pyinit
# Follow prompts and verify output
```

### Release Process
```bash
# Automated release via GitHub Actions
git tag v0.1.0
git push origin v0.1.0
# This triggers build for all platforms, Homebrew formula update, and PyPI publish
```

## Code Architecture

### Core Components

**Command Structure** (`internal/commands/`):
- `commands.go`: Cobra command setup and routing
- `interactive.go`: Main interactive flow execution
- `config.go`: Configuration management (`~/.pyinitrc`)
- `banner.go`: CLI banner display

**Project Generation** (`internal/generator/`):
- `generator.go`: Main project generation logic with type-specific handlers
- `utils.go`: File generation utilities

**Template System** (`pkg/template/`):
- `engine.go`: Pongo2-based template rendering with embedded template loader
- Uses `templates.go` embedded filesystem for template access

**User Interface** (`internal/prompts/`, `pkg/ui/`):
- `prompts.go`: Survey-based user prompts and validation
- `ui/`: Reusable UI components

### Template Structure
- `templates/core/`: Common files (gitignore, python-version, format scripts)
- `templates/basic/`: Basic Python project templates
- `templates/web/`: Web framework templates (FastAPI planned)

### Python Wrapper
- `python-package/pyinit_cli/`: Core Python CLI wrapper
- `python-package/scripts/`: Development utility scripts
- Downloads and executes Go binary via HTTP requests

## Key Design Patterns

1. **Embedded Templates**: Templates are embedded in Go binary using `//go:embed`
2. **Type-Specific Generation**: `GenerateProject()` routes to type-specific generators
3. **Configuration Persistence**: User preferences stored in `~/.pyinitrc`
4. **Hybrid Distribution**: Single codebase supporting both Go binary and Python package distribution

## Template Development

Templates use Pongo2 (Django-like) syntax:
- Variables: `{{ project_name }}`
- Conditionals: `{% if web_framework == "fastapi" %}`
- Located in `templates/` directory and embedded at build time

## Configuration

The tool uses a configuration system with:
- User config at `~/.pyinitrc`
- Project-specific settings in `internal/config/config.go`
- Runtime prompts for project customization