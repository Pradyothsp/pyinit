package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Pradyothsp/pyinit/internal/config"
	"github.com/Pradyothsp/pyinit/internal/prompts"
	"github.com/Pradyothsp/pyinit/pkg/template"
)

// Generator handles project generation
type Generator struct {
	templateEngine *template.Engine
}

// New creates a new Generator instance
func New() *Generator {
	return &Generator{
		templateEngine: template.NewEngine(),
	}
}

// GenerateProject creates the complete project structure
func (g *Generator) GenerateProject(cfg *config.ProjectConfig) error {
	// Check and create a project directory
	if err := g.createProjectDirectory(cfg); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Generate README.md
	if err := g.generateFileFromTemplate(cfg, "README.md.j2", "README.md"); err != nil {
		return fmt.Errorf("failed to generate README.md: %w", err)
	}

	// Create a scripts directory
	if err := g.createScriptsDirectory(cfg); err != nil {
		return fmt.Errorf("failed to create scripts directory: %w", err)
	}

	return nil
}

func (g *Generator) createProjectDirectory(cfg *config.ProjectConfig) error {
	// Ask for confirmation if the directory exists
	confirmed, err := prompts.ConfirmDirectoryCreation(cfg.ProjectPath)
	if err != nil {
		return err
	}
	if !confirmed {
		return fmt.Errorf("project creation cancelled")
	}

	// Create the directory
	if err := os.MkdirAll(cfg.ProjectPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}

func (g *Generator) generateFileFromTemplate(cfg *config.ProjectConfig, templateName, relativePath string) error {
	outputPath := filepath.Join(cfg.ProjectPath, relativePath)

	content, err := g.templateEngine.RenderTemplate(templateName, cfg.TemplateContext())
	if err != nil {
		return fmt.Errorf("failed to render %s template: %w", templateName, err)
	}

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", relativePath, err)
	}

	return nil
}

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
	if err := g.generateFileFromTemplate(cfg, "fmt.py.j2", filepath.Join("scripts", "fmt.py")); err != nil {
		return fmt.Errorf("failed to generate fmt.py: %w", err)
	}

	// Generate fmt_check.py
	if err := g.generateFileFromTemplate(cfg, "fmt_check.py.j2", filepath.Join("scripts", "fmt_check.py")); err != nil {
		return fmt.Errorf("failed to generate fmt_check.py: %w", err)
	}

	return nil
}
