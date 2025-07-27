package template

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/Pradyothsp/pyinit"
	"io"
	"path/filepath"

	"github.com/flosch/pongo2/v6"
)

// Engine handles template rendering using Pongo2 (Gonja-compatible)
type Engine struct {
	templateDir string
	loader      pongo2.TemplateLoader
}

// NewEngine creates a new template engine with embedded templates
func NewEngine() *Engine {
	// Create a loader that uses embedded files
	loader := &EmbeddedLoader{fs: pyinit.EmbeddedTemplates}

	return &Engine{
		templateDir: "embedded",
		loader:      loader,
	}
}

// EmbeddedLoader implements pongo2.TemplateLoader for embedded files
type EmbeddedLoader struct {
	fs embed.FS
}

func (e *EmbeddedLoader) Abs(base, name string) string {
	return filepath.Join("templates", name)
}

func (e *EmbeddedLoader) Get(path string) (io.Reader, error) {
	content, err := e.fs.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(content), nil
}

// RenderTemplate renders a template file with the given context
func (e *Engine) RenderTemplate(templateFile string, context map[string]interface{}) (string, error) {
	// Create a template set with our loader
	set := pongo2.NewSet("pyinit", e.loader)

	// Get the template
	template, err := set.FromFile(templateFile)
	if err != nil {
		return "", fmt.Errorf("failed to load template %s: %w", templateFile, err)
	}

	// Render the template
	output, err := template.Execute(pongo2.Context(context))
	if err != nil {
		return "", fmt.Errorf("failed to render template %s: %w", templateFile, err)
	}

	return output, nil
}

// SetTemplateDir sets a custom template directory
func (e *Engine) SetTemplateDir(dir string) error {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for template directory: %w", err)
	}

	e.templateDir = absDir
	e.loader = pongo2.MustNewLocalFileSystemLoader(absDir)

	return nil
}
