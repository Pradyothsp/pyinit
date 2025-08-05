package generator

import (
	"fmt"
	"github.com/Pradyothsp/pyinit/internal/config"
	"os"
	"path/filepath"
)

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
