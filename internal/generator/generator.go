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
	if err := g.generateReadme(cfg); err != nil {
		return fmt.Errorf("failed to generate README.md: %w", err)
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

func (g *Generator) generateReadme(cfg *config.ProjectConfig) error {
	templatePath := "README.md.j2"
	outputPath := filepath.Join(cfg.ProjectPath, "README.md")

	content, err := g.templateEngine.RenderTemplate(templatePath, cfg.TemplateContext())
	if err != nil {
		return fmt.Errorf("failed to render README template: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write README.md: %w", err)
	}

	return nil
}
