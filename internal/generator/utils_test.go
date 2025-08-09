package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Pradyothsp/pyinit/internal/config"
)

// TemplateEngine interface for mocking
type MockTemplateEngineInterface interface {
	RenderTemplate(templateFile string, context map[string]interface{}) (string, error)
}

// SimpleMockTemplateEngine for testing utils
type SimpleMockTemplateEngine struct {
	templates map[string]string
	rendered  []string
}

func NewSimpleMockTemplateEngine(templates map[string]string) *SimpleMockTemplateEngine {
	return &SimpleMockTemplateEngine{
		templates: templates,
		rendered:  make([]string, 0),
	}
}

func (m *SimpleMockTemplateEngine) RenderTemplate(templateFile string, context map[string]interface{}) (string, error) {
	m.rendered = append(m.rendered, templateFile)
	
	if content, exists := m.templates[templateFile]; exists {
		// Simple variable substitution for testing
		result := content
		for key, value := range context {
			placeholder := fmt.Sprintf("{{ %s }}", key)
			result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
		}
		return result, nil
	}
	return "", fmt.Errorf("template not found: %s", templateFile)
}

// Test helper functions
func createTempTestDir(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "pyinit-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return tempDir
}

func cleanupTestDir(t *testing.T, dir string) {
	if err := os.RemoveAll(dir); err != nil {
		t.Errorf("Failed to cleanup temp dir %s: %v", dir, err)
	}
}

func createBasicTestConfig(projectPath string) *config.ProjectConfig {
	return &config.ProjectConfig{
		UserName:           "Test User",
		Email:              "test@example.com",
		ProjectName:        "test-project",
		ProjectDescription: "Test project description",
		ProjectType:        "basic",
		WebFramework:       "",
		ProjectPath:        projectPath,
		MainDirName:        "test_project",
		PythonVersion:      "3.11",
	}
}

// Mock Generator for testing utility functions
type MockGenerator struct {
	templateEngine *SimpleMockTemplateEngine
}

func NewMockGenerator(templates map[string]string) *MockGenerator {
	return &MockGenerator{
		templateEngine: NewSimpleMockTemplateEngine(templates),
	}
}

func (mg *MockGenerator) generateFileFromTemplate(cfg *config.ProjectConfig, templateName, relativePath string) error {
	outputPath := filepath.Join(cfg.ProjectPath, relativePath)

	// Ensure directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	content, err := mg.templateEngine.RenderTemplate(templateName, cfg.TemplateContext())
	if err != nil {
		return fmt.Errorf("failed to render %s template: %w", templateName, err)
	}

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", relativePath, err)
	}

	return nil
}

func TestGenerateFileFromTemplate(t *testing.T) {
	tempDir := createTempTestDir(t)
	defer cleanupTestDir(t, tempDir)

	templates := map[string]string{
		"test.j2": "Hello {{ project_name }}!\nUser: {{ user_name }}",
	}

	generator := NewMockGenerator(templates)
	cfg := createBasicTestConfig(tempDir)

	err := generator.generateFileFromTemplate(cfg, "test.j2", "test.txt")
	if err != nil {
		t.Fatalf("generateFileFromTemplate failed: %v", err)
	}

	// Verify file was created
	filePath := filepath.Join(tempDir, "test.txt")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	expectedContent := "Hello test-project!\nUser: Test User"
	if string(content) != expectedContent {
		t.Errorf("File content = %q, want %q", string(content), expectedContent)
	}

	// Verify template was rendered
	if len(generator.templateEngine.rendered) != 1 {
		t.Errorf("Expected 1 template to be rendered, got %d", len(generator.templateEngine.rendered))
	}

	if generator.templateEngine.rendered[0] != "test.j2" {
		t.Errorf("Expected template 'test.j2' to be rendered, got %q", generator.templateEngine.rendered[0])
	}
}

func TestGenerateFileFromTemplate_NestedPath(t *testing.T) {
	tempDir := createTempTestDir(t)
	defer cleanupTestDir(t, tempDir)

	templates := map[string]string{
		"nested.j2": "Content for {{ project_name }}",
	}

	generator := NewMockGenerator(templates)
	cfg := createBasicTestConfig(tempDir)

	err := generator.generateFileFromTemplate(cfg, "nested.j2", "subdir/nested/test.txt")
	if err != nil {
		t.Fatalf("generateFileFromTemplate failed: %v", err)
	}

	// Verify nested directories were created
	filePath := filepath.Join(tempDir, "subdir", "nested", "test.txt")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read nested file: %v", err)
	}

	expectedContent := "Content for test-project"
	if string(content) != expectedContent {
		t.Errorf("File content = %q, want %q", string(content), expectedContent)
	}
}

func TestGenerateFileFromTemplate_TemplateNotFound(t *testing.T) {
	tempDir := createTempTestDir(t)
	defer cleanupTestDir(t, tempDir)

	generator := NewMockGenerator(map[string]string{})
	cfg := createBasicTestConfig(tempDir)

	err := generator.generateFileFromTemplate(cfg, "nonexistent.j2", "test.txt")
	if err == nil {
		t.Fatal("Expected error for nonexistent template, got nil")
	}

	if !strings.Contains(err.Error(), "failed to render nonexistent.j2 template") {
		t.Errorf("Expected template render error, got: %v", err)
	}
}

func TestGenerateFileFromTemplate_InvalidPath(t *testing.T) {
	tempDir := createTempTestDir(t)
	defer cleanupTestDir(t, tempDir)

	// Create a file where we want to create a directory
	blockingFile := filepath.Join(tempDir, "blocking")
	if err := os.WriteFile(blockingFile, []byte("blocking"), 0644); err != nil {
		t.Fatalf("Failed to create blocking file: %v", err)
	}

	templates := map[string]string{
		"test.j2": "content",
	}

	generator := NewMockGenerator(templates)
	cfg := createBasicTestConfig(tempDir)

	// Try to create a file in a path that requires creating a directory 
	// where a file already exists
	err := generator.generateFileFromTemplate(cfg, "test.j2", "blocking/test.txt")
	if err == nil {
		t.Fatal("Expected error when trying to create directory where file exists, got nil")
	}
}

func TestGenerateFileFromTemplate_ContextVariables(t *testing.T) {
	tempDir := createTempTestDir(t)
	defer cleanupTestDir(t, tempDir)

	templates := map[string]string{
		"context.j2": `Project: {{ project_name }}
Type: {{ project_type }}
Description: {{ project_description }}
User: {{ user_name }}
Email: {{ email }}
Main Dir: {{ main_dir_name }}
Python: {{ python_version }}
Ruff: {{ python_version_for_ruff }}`,
	}

	generator := NewMockGenerator(templates)
	cfg := createBasicTestConfig(tempDir)

	err := generator.generateFileFromTemplate(cfg, "context.j2", "context.txt")
	if err != nil {
		t.Fatalf("generateFileFromTemplate failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tempDir, "context.txt"))
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	expectedLines := []string{
		"Project: test-project",
		"Type: basic",
		"Description: Test project description",
		"User: Test User",
		"Email: test@example.com",
		"Main Dir: test_project",
		"Python: 3.11",
		"Ruff: py311",
	}

	contentStr := string(content)
	for _, expected := range expectedLines {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Expected content to contain %q, but it didn't.\nContent:\n%s", expected, contentStr)
		}
	}
}

func TestGenerateFileFromTemplate_EmptyTemplate(t *testing.T) {
	tempDir := createTempTestDir(t)
	defer cleanupTestDir(t, tempDir)

	templates := map[string]string{
		"empty.j2": "",
	}

	generator := NewMockGenerator(templates)
	cfg := createBasicTestConfig(tempDir)

	err := generator.generateFileFromTemplate(cfg, "empty.j2", "empty.txt")
	if err != nil {
		t.Fatalf("generateFileFromTemplate failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tempDir, "empty.txt"))
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if len(content) != 0 {
		t.Errorf("Expected empty file, got %d bytes: %q", len(content), string(content))
	}
}

// Benchmark the file generation utility
func BenchmarkGenerateFileFromTemplate(b *testing.B) {
	tempDir, _ := os.MkdirTemp("", "pyinit-bench-*")
	defer os.RemoveAll(tempDir)

	templates := map[string]string{
		"bench.j2": "Project: {{ project_name }}\nType: {{ project_type }}\nUser: {{ user_name }}\nEmail: {{ email }}",
	}

	generator := NewMockGenerator(templates)
	cfg := createBasicTestConfig(tempDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filename := fmt.Sprintf("bench_%d.txt", i)
		err := generator.generateFileFromTemplate(cfg, "bench.j2", filename)
		if err != nil {
			b.Fatalf("generateFileFromTemplate failed: %v", err)
		}
	}
}

func BenchmarkGenerateFileFromTemplate_NestedPath(b *testing.B) {
	tempDir, _ := os.MkdirTemp("", "pyinit-bench-*")
	defer os.RemoveAll(tempDir)

	templates := map[string]string{
		"bench.j2": "Nested content for {{ project_name }}",
	}

	generator := NewMockGenerator(templates)
	cfg := createBasicTestConfig(tempDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filename := fmt.Sprintf("level1/level2/level3/bench_%d.txt", i)
		err := generator.generateFileFromTemplate(cfg, "bench.j2", filename)
		if err != nil {
			b.Fatalf("generateFileFromTemplate failed: %v", err)
		}
	}
}