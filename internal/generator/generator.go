package generator

import (
	"fmt"

	"github.com/Pradyothsp/pyinit/internal/config"
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
	// Common steps
	if err := g.GenerateCommonProject(cfg); err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	if cfg.ProjectType == "basic" {
		if err := g.GeneratorBasicProject(cfg); err != nil {
			return fmt.Errorf("failed to create basic project structure %w", err)
		}
	}

	if cfg.ProjectType == "web" {
		if cfg.WebFramework == "fastapi" {
			if err := g.GenerateFastAPIProject(cfg); err != nil {
				return fmt.Errorf("failed to create fastapi project %w", err)
			}
		}
	}

	return nil
}
