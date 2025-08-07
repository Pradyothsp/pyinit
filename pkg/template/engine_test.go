package template

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"strings"
	"testing"
)

//go:embed  test_templates/*
var testTemplatesFS embed.FS

// Test templates for isolated testing
const (
	simpleTemplate  = `Hello {{ name }}!`
	complexTemplate = `# {{ project_name }}

A {{ project_type }} project by {{ user_name }}.

{% if project_description %}
## Description
{{ project_description }}
{% endif %}

{% for feature in features %}
- {{ feature }}
{% endfor %}`
)

// MockLoader implements pongo2.TemplateLoader for testing
type MockLoader struct {
	templates map[string]string
}

func NewMockLoader(templates map[string]string) *MockLoader {
	return &MockLoader{templates: templates}
}

func (m *MockLoader) Abs(base, name string) string {
	return name
}

func (m *MockLoader) Get(path string) (io.Reader, error) {
	if content, exists := m.templates[path]; exists {
		return bytes.NewReader([]byte(content)), nil
	}
	return nil, fmt.Errorf("template not found: %s", path)
}

func TestNewEngine(t *testing.T) {
	engine := NewEngine()

	if engine == nil {
		t.Fatal("NewEngine() returned nil")
	}

	if engine.templateDir != "embedded" {
		t.Errorf("Expected templateDir to be 'embedded', got %s", engine.templateDir)
	}

	if engine.loader == nil {
		t.Fatal("Engine loader is nil")
	}
}

func TestEmbeddedLoader_Abs(t *testing.T) {
	loader := &EmbeddedLoader{fs: testTemplatesFS}

	tests := []struct {
		base     string
		name     string
		expected string
	}{
		{"", "basic/README.md.j2", "templates/basic/README.md.j2"},
		{"templates", "core/python-version.j2", "templates/core/python-version.j2"},
		{"/some/path", "web/fastapi/main.py.j2", "templates/web/fastapi/main.py.j2"},
	}

	for _, tt := range tests {
		result := loader.Abs(tt.base, tt.name)
		if result != tt.expected {
			t.Errorf("Abs(%q, %q) = %q, want %q", tt.base, tt.name, result, tt.expected)
		}
	}
}

func TestEngine_RenderTemplate_Simple(t *testing.T) {
	engine := &Engine{
		templateDir: "test",
		loader: NewMockLoader(map[string]string{
			"simple.j2": simpleTemplate,
		}),
	}

	context := map[string]interface{}{
		"name": "World",
	}

	result, err := engine.RenderTemplate("simple.j2", context)
	if err != nil {
		t.Fatalf("RenderTemplate failed: %v", err)
	}

	expected := "Hello World!"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestEngine_RenderTemplate_Complex(t *testing.T) {
	engine := &Engine{
		templateDir: "test",
		loader: NewMockLoader(map[string]string{
			"complex.j2": complexTemplate,
		}),
	}

	context := map[string]interface{}{
		"project_name":        "MyProject",
		"project_type":        "web",
		"user_name":           "John Doe",
		"project_description": "An awesome project",
		"features":            []string{"Feature 1", "Feature 2", "Feature 3"},
	}

	result, err := engine.RenderTemplate("complex.j2", context)
	if err != nil {
		t.Fatalf("RenderTemplate failed: %v", err)
	}

	// Check that all expected content is present
	expectedParts := []string{
		"# MyProject",
		"A web project by John Doe.",
		"## Description",
		"An awesome project",
		"- Feature 1",
		"- Feature 2",
		"- Feature 3",
	}

	for _, part := range expectedParts {
		if !strings.Contains(result, part) {
			t.Errorf("Expected result to contain %q, but it didn't.\nResult:\n%s", part, result)
		}
	}
}

func TestEngine_RenderTemplate_ConditionalContent(t *testing.T) {
	engine := &Engine{
		templateDir: "test",
		loader: NewMockLoader(map[string]string{
			"conditional.j2": complexTemplate,
		}),
	}

	// Test with description
	contextWithDescription := map[string]interface{}{
		"project_name":        "TestProject",
		"project_type":        "basic",
		"user_name":           "Jane Doe",
		"project_description": "Test description",
		"features":            []string{},
	}

	result, err := engine.RenderTemplate("conditional.j2", contextWithDescription)
	if err != nil {
		t.Fatalf("RenderTemplate failed: %v", err)
	}

	if !strings.Contains(result, "## Description") {
		t.Error("Expected description section when project_description is provided")
	}

	// Test without description
	contextWithoutDescription := map[string]interface{}{
		"project_name": "TestProject",
		"project_type": "basic",
		"user_name":    "Jane Doe",
		"features":     []string{},
	}

	result, err = engine.RenderTemplate("conditional.j2", contextWithoutDescription)
	if err != nil {
		t.Fatalf("RenderTemplate failed: %v", err)
	}

	if strings.Contains(result, "## Description") {
		t.Error("Did not expect description section when project_description is not provided")
	}
}

func TestEngine_RenderTemplate_MissingVariable(t *testing.T) {
	engine := &Engine{
		templateDir: "test",
		loader: NewMockLoader(map[string]string{
			"missing.j2": "Hello {{ missing_variable }}!",
		}),
	}

	context := map[string]interface{}{
		"name": "World",
	}

	result, err := engine.RenderTemplate("missing.j2", context)
	if err != nil {
		t.Fatalf("RenderTemplate failed: %v", err)
	}

	// Pongo2 renders missing variables as empty strings
	expected := "Hello !"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestEngine_RenderTemplate_TemplateNotFound(t *testing.T) {
	engine := &Engine{
		templateDir: "test",
		loader:      NewMockLoader(map[string]string{}),
	}

	context := map[string]interface{}{}

	_, err := engine.RenderTemplate("nonexistent.j2", context)
	if err == nil {
		t.Fatal("Expected error when template not found, but got none")
	}

	if !strings.Contains(err.Error(), "failed to load template") {
		t.Errorf("Expected error message about loading template, got: %v", err)
	}
}

func TestEngine_RenderTemplate_EmptyContext(t *testing.T) {
	engine := &Engine{
		templateDir: "test",
		loader: NewMockLoader(map[string]string{
			"empty.j2": "Static content only",
		}),
	}

	result, err := engine.RenderTemplate("empty.j2", nil)
	if err != nil {
		t.Fatalf("RenderTemplate failed: %v", err)
	}

	expected := "Static content only"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestEngine_RenderTemplate_SpecialCharacters(t *testing.T) {
	template := `Project: {{ project_name }}
Special chars: !@#$%^&*()
Unicode: 🚀 {{ emoji }}`

	engine := &Engine{
		templateDir: "test",
		loader: NewMockLoader(map[string]string{
			"special.j2": template,
		}),
	}

	context := map[string]interface{}{
		"project_name": "Test-Project_123",
		"emoji":        "✨",
	}

	result, err := engine.RenderTemplate("special.j2", context)
	if err != nil {
		t.Fatalf("RenderTemplate failed: %v", err)
	}

	expectedParts := []string{
		"Project: Test-Project_123",
		"Special chars: !@#$%^&*()",
		"Unicode: 🚀 ✨",
	}

	for _, part := range expectedParts {
		if !strings.Contains(result, part) {
			t.Errorf("Expected result to contain %q, but it didn't.\nResult:\n%s", part, result)
		}
	}
}

func TestEngine_SetTemplateDir(t *testing.T) {
	engine := NewEngine()

	// Use current directory since it exists
	err := engine.SetTemplateDir(".")
	if err != nil {
		t.Fatalf("SetTemplateDir failed: %v", err)
	}

	// The path should be converted to absolute and should contain current directory
	if engine.templateDir == "" {
		t.Error("Expected templateDir to be set to absolute path")
	}

	// Verify the loader was updated
	if engine.loader == nil {
		t.Error("Expected loader to be set after SetTemplateDir")
	}
}

// Integration test with real embedded templates
func TestEngine_RenderTemplate_RealTemplates(t *testing.T) {
	engine := NewEngine()

	context := map[string]interface{}{
		"python_version": "3.11",
	}

	result, err := engine.RenderTemplate("core/python-version.j2", context)
	if err != nil {
		t.Fatalf("RenderTemplate with real template failed: %v", err)
	}

	expected := "3.11"
	if strings.TrimSpace(result) != expected {
		t.Errorf("Expected %q, got %q", expected, strings.TrimSpace(result))
	}
}

// Benchmark tests
func BenchmarkEngine_RenderTemplate_Simple(b *testing.B) {
	engine := &Engine{
		templateDir: "test",
		loader: NewMockLoader(map[string]string{
			"simple.j2": simpleTemplate,
		}),
	}

	context := map[string]interface{}{
		"name": "World",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.RenderTemplate("simple.j2", context)
		if err != nil {
			b.Fatalf("RenderTemplate failed: %v", err)
		}
	}
}

func BenchmarkEngine_RenderTemplate_Complex(b *testing.B) {
	engine := &Engine{
		templateDir: "test",
		loader: NewMockLoader(map[string]string{
			"complex.j2": complexTemplate,
		}),
	}

	context := map[string]interface{}{
		"project_name":        "BenchmarkProject",
		"project_type":        "web",
		"user_name":           "Benchmark User",
		"project_description": "A project for benchmarking template rendering performance",
		"features":            []string{"Fast", "Reliable", "Scalable", "Maintainable", "Secure"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.RenderTemplate("complex.j2", context)
		if err != nil {
			b.Fatalf("RenderTemplate failed: %v", err)
		}
	}
}

func BenchmarkEngine_RenderTemplate_RealTemplate(b *testing.B) {
	engine := NewEngine()

	context := map[string]interface{}{
		"project_name":        "BenchmarkProject",
		"project_type":        "basic",
		"user_name":           "Benchmark User",
		"email":               "bench@example.com",
		"project_description": "A project for benchmarking real template performance",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.RenderTemplate("basic/README.md.j2", context)
		if err != nil {
			b.Fatalf("RenderTemplate failed: %v", err)
		}
	}
}

