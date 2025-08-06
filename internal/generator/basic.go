package generator

import (
	"fmt"

	"github.com/Pradyothsp/pyinit/internal/config"
)

// GeneratorBasicProject creates a basic Python project structure
func (g *Generator) GeneratorBasicProject(cfg *config.ProjectConfig) error {
	// Generate README.md
	if err := g.generateFileFromTemplate(cfg, "basic/README.md.j2", "README.md"); err != nil {
		return fmt.Errorf("failed to generate README.md: %w", err)
	}

	// Create the main project structure
	if err := g.createMainDirectory(cfg); err != nil {
		return fmt.Errorf("failed to create main project directory: %w", err)
	}

	// Create pyproject.toml
	if err := g.generateFileFromTemplate(cfg, "basic/pyproject.toml.j2", "pyproject.toml"); err != nil {
		return fmt.Errorf("failed to generate pyproject.toml: %w", err)
	}

	return nil
}