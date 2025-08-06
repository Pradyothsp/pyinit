package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Pradyothsp/pyinit/internal/config"
	"github.com/Pradyothsp/pyinit/internal/prompts"
)

// GenerateCommonProject creates common project files and directories
func (g *Generator) GenerateCommonProject(cfg *config.ProjectConfig) error {
	// Check and create a project directory
	if err := g.createProjectDirectory(cfg); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create a scripts directory
	if err := g.createScriptsDirectory(cfg); err != nil {
		return fmt.Errorf("failed to create scripts directory: %w", err)
	}

	// Generate .gitignore
	if err := g.generateFileFromTemplate(cfg, "core/gitignore.j2", ".gitignore"); err != nil {
		return fmt.Errorf("failed to generate .gitignore: %w", err)
	}

	// Generate .python-version
	if err := g.generateFileFromTemplate(cfg, "core/python-version.j2", ".python-version"); err != nil {
		return fmt.Errorf("failed to generate .python-version: %w", err)
	}

	return nil
}

// createProjectDirectory creates the main project directory with user confirmation
func (g *Generator) createProjectDirectory(cfg *config.ProjectConfig) error {
	// Ask for confirmation if the directory exists
	confirmed, err := prompts.ConfirmDirectoryCreation(cfg.ProjectPath)
	if err != nil {
		return err
	}
	if !confirmed {
		return fmt.Errorf("project creation cancelled %w", err)
	}

	// Create the directory
	if err := os.MkdirAll(cfg.ProjectPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}

// createScriptsDirectory creates the scripts directory with development utilities
func (g *Generator) createScriptsDirectory(cfg *config.ProjectConfig) error {
	scriptsDir := filepath.Join(cfg.ProjectPath, "scripts")
	if err := os.MkdirAll(scriptsDir, 0755); err != nil {
		return fmt.Errorf("failed to create scripts directory: %w", err)
	}

	// Create an empty __init__.py file
	initPath := filepath.Join(scriptsDir, "__init__.py")
	if err := os.WriteFile(initPath, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to create __init__.py in scripts: %w", err)
	}

	// Generate fmt.py
	if err := g.generateFileFromTemplate(cfg, "core/fmt.py.j2", filepath.Join("scripts", "fmt.py")); err != nil {
		return fmt.Errorf("failed to generate fmt.py: %w", err)
	}

	// Generate fmt_check.py
	if err := g.generateFileFromTemplate(cfg, "core/fmt_check.py.j2", filepath.Join("scripts", "fmt_check.py")); err != nil {
		return fmt.Errorf("failed to generate fmt_check.py: %w", err)
	}

	return nil
}

// createMainDirectory creates the main project package directory
func (g *Generator) createMainDirectory(cfg *config.ProjectConfig) error {
	mainDir := filepath.Join(cfg.ProjectPath, cfg.MainDirName)
	if err := os.MkdirAll(mainDir, 0755); err != nil {
		return fmt.Errorf("failed to create main project directory: %w", err)
	}

	// Generate __init__.py in the main project directory
	initPath := filepath.Join(mainDir, "__init__.py")
	if err := os.WriteFile(initPath, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to create __init__.py in main project directory: %w", err)
	}

	// Generate main.py
	if err := g.generateFileFromTemplate(cfg, "basic/main.py.j2", filepath.Join(cfg.MainDirName, "main.py")); err != nil {
		return fmt.Errorf("failed to generate main.py: %w", err)
	}

	return nil
}