package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Pradyothsp/pyinit/internal/config"
)

// GenerateFastAPIProject creates a FastAPI web project structure
func (g *Generator) GenerateFastAPIProject(cfg *config.ProjectConfig) error {
	// Generate FastAPI-specific README.md
	if err := g.generateFileFromTemplate(cfg, "web/fastapi/README.md.j2", "README.md"); err != nil {
		return fmt.Errorf("failed to generate README.md: %w", err)
	}

	// Create the main project structure (without main.py)
	if err := g.createFastAPIMainDirectory(cfg); err != nil {
		return fmt.Errorf("failed to create main project directory: %w", err)
	}

	// Generate FastAPI-specific pyproject.toml
	if err := g.generateFileFromTemplate(cfg, "web/fastapi/pyproject.toml.j2", "pyproject.toml"); err != nil {
		return fmt.Errorf("failed to generate pyproject.toml: %w", err)
	}

	// Generate main FastAPI application
	if err := g.generateFileFromTemplate(cfg, "web/fastapi/main.py.j2", filepath.Join(cfg.MainDirName, "main.py")); err != nil {
		return fmt.Errorf("failed to generate main.py: %w", err)
	}

	// Create API directory structure
	if err := g.createFastAPIDirectories(cfg); err != nil {
		return fmt.Errorf("failed to create FastAPI directories: %w", err)
	}

	// Create tests directory structure
	if err := g.createFastAPITestsDirectory(cfg); err != nil {
		return fmt.Errorf("failed to create tests directory: %w", err)
	}

	return nil
}

// createFastAPIDirectories creates the FastAPI-specific directory structure
func (g *Generator) createFastAPIDirectories(cfg *config.ProjectConfig) error {
	mainDir := filepath.Join(cfg.ProjectPath, cfg.MainDirName)

	// Create api directory
	apiDir := filepath.Join(mainDir, "api")
	if err := os.MkdirAll(apiDir, 0755); err != nil {
		return fmt.Errorf("failed to create api directory: %w", err)
	}

	// Create __init__.py in api directory
	initPath := filepath.Join(apiDir, "__init__.py")
	if err := os.WriteFile(initPath, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to create __init__.py in api: %w", err)
	}

	// Generate routes.py
	if err := g.generateFileFromTemplate(cfg, "web/fastapi/api/routes.py.j2", filepath.Join(cfg.MainDirName, "api", "routes.py")); err != nil {
		return fmt.Errorf("failed to generate routes.py: %w", err)
	}

	// Create core directory
	coreDir := filepath.Join(mainDir, "core")
	if err := os.MkdirAll(coreDir, 0755); err != nil {
		return fmt.Errorf("failed to create core directory: %w", err)
	}

	// Create __init__.py in core directory
	coreInitPath := filepath.Join(coreDir, "__init__.py")
	if err := os.WriteFile(coreInitPath, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to create __init__.py in core: %w", err)
	}

	// Generate config.py
	if err := g.generateFileFromTemplate(cfg, "web/fastapi/core/config.py.j2", filepath.Join(cfg.MainDirName, "core", "config.py")); err != nil {
		return fmt.Errorf("failed to generate config.py: %w", err)
	}

	// Create schemas directory with placeholder
	schemasDir := filepath.Join(mainDir, "schemas")
	if err := os.MkdirAll(schemasDir, 0755); err != nil {
		return fmt.Errorf("failed to create schemas directory: %w", err)
	}

	schemasInitPath := filepath.Join(schemasDir, "__init__.py")
	if err := os.WriteFile(schemasInitPath, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to create __init__.py in schemas: %w", err)
	}

	// Generate placeholder schemas/user.py
	if err := g.generateFileFromTemplate(cfg, "web/fastapi/schemas/user.py.j2", filepath.Join(cfg.MainDirName, "schemas", "user.py")); err != nil {
		return fmt.Errorf("failed to generate schemas/user.py: %w", err)
	}

	// Create models directory with placeholder
	modelsDir := filepath.Join(mainDir, "models")
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		return fmt.Errorf("failed to create models directory: %w", err)
	}

	modelsInitPath := filepath.Join(modelsDir, "__init__.py")
	if err := os.WriteFile(modelsInitPath, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to create __init__.py in models: %w", err)
	}

	// Generate placeholder models/user.py
	if err := g.generateFileFromTemplate(cfg, "web/fastapi/models/user.py.j2", filepath.Join(cfg.MainDirName, "models", "user.py")); err != nil {
		return fmt.Errorf("failed to generate models/user.py: %w", err)
	}

	return nil
}

// createFastAPITestsDirectory creates the tests directory structure
func (g *Generator) createFastAPITestsDirectory(cfg *config.ProjectConfig) error {
	testsDir := filepath.Join(cfg.ProjectPath, "tests")
	if err := os.MkdirAll(testsDir, 0755); err != nil {
		return fmt.Errorf("failed to create tests directory: %w", err)
	}

	// Create __init__.py in tests directory
	initPath := filepath.Join(testsDir, "__init__.py")
	if err := os.WriteFile(initPath, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to create __init__.py in tests: %w", err)
	}

	// Generate test_main.py
	if err := g.generateFileFromTemplate(cfg, "web/fastapi/tests/test_main.py.j2", filepath.Join("tests", "test_main.py")); err != nil {
		return fmt.Errorf("failed to generate test_main.py: %w", err)
	}

	return nil
}

// createFastAPIMainDirectory creates the main project directory without generating main.py
func (g *Generator) createFastAPIMainDirectory(cfg *config.ProjectConfig) error {
	mainDir := filepath.Join(cfg.ProjectPath, cfg.MainDirName)
	if err := os.MkdirAll(mainDir, 0755); err != nil {
		return fmt.Errorf("failed to create main project directory: %w", err)
	}

	// Generate __init__.py in the main project directory
	initPath := filepath.Join(mainDir, "__init__.py")
	if err := os.WriteFile(initPath, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to create __init__.py in main project directory: %w", err)
	}

	return nil
}

