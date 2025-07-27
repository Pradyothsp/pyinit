# pyinit

An interactive CLI tool to create Python project scaffolds with a customizable structure.

## Features

- **Interactive Setup**: Guides you through creating a new Python project.
- **Customizable Project Structure**: Choose between a `src` layout or a direct layout.
- **Project Types**: Supports different project types like `cli` and `package`.
- **Automated Environment Setup**: Automatically creates a virtual environment and installs dependencies.
- **Pre-configured Tools**: Comes with pre-configured tools for formatting (`black`, `isort`) and linting (`ruff`).
- **GitHub Actions**: Includes basic GitHub Actions for CI/CD.

## Usage

To use `pyinit`, simply run the command:

```bash
pyinit
```

The tool will then guide you through the setup process, asking for the following information:

- Your name and email
- Project name, type, and description
- Project structure (`src` or direct)
- Python version

After gathering the required information, `pyinit` will generate the project structure and offer to set up the development environment.

## Generated Project Structure

The generated project will have the following structure:

```
.
├── .gitignore
├── .python-version
├── pyproject.toml
├── README.md
├── <project_name>/
│   └── __init__.py
└── tests/
    └── __init__.py
```

If you choose the `src` layout, the structure will be:

```
.
├── .gitignore
├── .python-version
├── pyproject.toml
├── README.md
├── src/
│   └── <project_name>/
│       └── __init__.py
└── tests/
    └── __init__.py
```

## License

This project is licensed under the terms of the MIT license. See [LICENSE](LICENSE) for more details.
