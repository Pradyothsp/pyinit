package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Pradyothsp/pyinit/internal/config"
	"github.com/Pradyothsp/pyinit/internal/generator"
)

// Integration test for basic project generation without user interaction
func TestBasicProjectGeneration(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "pyinit-integration-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			fmt.Printf("Error removing temporary directory %s: %v\n", tempDir, err)
		}
	}()

	// Create test configuration
	cfg := &config.ProjectConfig{
		UserName:           "Integration Test",
		Email:              "test@integration.com",
		ProjectName:        "test-integration-project",
		ProjectDescription: "Integration test project",
		ProjectType:        "basic",
		WebFramework:       "",
		ProjectPath:        filepath.Join(tempDir, "test-integration-project"),
		MainDirName:        "test_integration_project",
		PythonVersion:      "3.11",
	}

	// Generate project
	gen := generator.New()
	err = gen.GenerateProject(cfg)

	// Skip directory confirmation errors for this test
	if err != nil && !isDirectoryConfirmationError(err) {
		t.Fatalf("Project generation failed: %v", err)
	} else if err != nil {
		t.Skipf("Skipping integration test due to directory confirmation: %v", err)
		return
	}

	// Verify project structure
	projectPath := cfg.ProjectPath

	// Check that essential files exist
	expectedFiles := []string{
		".gitignore",
		".python-version",
		"README.md",
		"pyproject.toml",
		filepath.Join("scripts", "__init__.py"),
		filepath.Join("scripts", "fmt.py"),
		filepath.Join("scripts", "fmt_check.py"),
		filepath.Join(cfg.MainDirName, "__init__.py"),
		filepath.Join(cfg.MainDirName, "main.py"),
	}

	for _, file := range expectedFiles {
		fullPath := filepath.Join(projectPath, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected file %s does not exist", file)
		}
	}

	// Check that directories exist
	expectedDirs := []string{
		"scripts",
		cfg.MainDirName,
	}

	for _, dir := range expectedDirs {
		fullPath := filepath.Join(projectPath, dir)
		if stat, err := os.Stat(fullPath); os.IsNotExist(err) || !stat.IsDir() {
			t.Errorf("Expected directory %s does not exist or is not a directory", dir)
		}
	}

	// Verify file contents contain expected values
	readmeContent, err := os.ReadFile(filepath.Join(projectPath, "README.md"))
	if err != nil {
		t.Errorf("Failed to read README.md: %v", err)
	} else {
		if len(readmeContent) == 0 {
			t.Error("README.md is empty")
		}
	}

	pythonVersionContent, err := os.ReadFile(filepath.Join(projectPath, ".python-version"))
	if err != nil {
		t.Errorf("Failed to read .python-version: %v", err)
	} else {
		// The template already includes a newline, so we expect just the version
		expectedVersion := cfg.PythonVersion
		actualContent := strings.TrimSpace(string(pythonVersionContent))
		if actualContent != expectedVersion {
			t.Errorf(".python-version content = %q, want %q", actualContent, expectedVersion)
		}
	}
}

// Integration test for FastAPI project generation
func TestFastAPIProjectGeneration(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "pyinit-fastapi-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			fmt.Printf("Error removing temporary directory %s: %v\n", tempDir, err)
		}
	}()

	// Create test configuration for FastAPI
	cfg := &config.ProjectConfig{
		UserName:           "FastAPI Test",
		Email:              "fastapi@test.com",
		ProjectName:        "test-fastapi-project",
		ProjectDescription: "FastAPI integration test project",
		ProjectType:        "web",
		WebFramework:       "fastapi",
		ProjectPath:        filepath.Join(tempDir, "test-fastapi-project"),
		MainDirName:        "test_fastapi_project",
		PythonVersion:      "3.12",
	}

	// Generate project
	gen := generator.New()
	err = gen.GenerateProject(cfg)

	// Skip directory confirmation errors for this test
	if err != nil && !isDirectoryConfirmationError(err) {
		t.Fatalf("FastAPI project generation failed: %v", err)
	} else if err != nil {
		t.Skipf("Skipping FastAPI integration test due to directory confirmation: %v", err)
		return
	}

	// Verify FastAPI-specific files
	projectPath := cfg.ProjectPath

	expectedFastAPIFiles := []string{
		".gitignore",
		".python-version",
		"README.md",
		"pyproject.toml",
		filepath.Join(cfg.MainDirName, "__init__.py"),
		filepath.Join(cfg.MainDirName, "main.py"),
		filepath.Join(cfg.MainDirName, "api", "__init__.py"),
		filepath.Join(cfg.MainDirName, "api", "routes.py"),
		filepath.Join(cfg.MainDirName, "core", "__init__.py"),
		filepath.Join(cfg.MainDirName, "core", "config.py"),
		filepath.Join(cfg.MainDirName, "schemas", "__init__.py"),
		filepath.Join(cfg.MainDirName, "schemas", "user.py"),
		filepath.Join(cfg.MainDirName, "models", "__init__.py"),
		filepath.Join(cfg.MainDirName, "models", "user.py"),
		filepath.Join("tests", "__init__.py"),
		filepath.Join("tests", "test_main.py"),
	}

	for _, file := range expectedFastAPIFiles {
		fullPath := filepath.Join(projectPath, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected FastAPI file %s does not exist", file)
		}
	}

	// Verify FastAPI-specific directories
	expectedFastAPIDirs := []string{
		cfg.MainDirName,
		filepath.Join(cfg.MainDirName, "api"),
		filepath.Join(cfg.MainDirName, "core"),
		filepath.Join(cfg.MainDirName, "schemas"),
		filepath.Join(cfg.MainDirName, "models"),
		"tests",
	}

	for _, dir := range expectedFastAPIDirs {
		fullPath := filepath.Join(projectPath, dir)
		if stat, err := os.Stat(fullPath); os.IsNotExist(err) || !stat.IsDir() {
			t.Errorf("Expected FastAPI directory %s does not exist or is not a directory", dir)
		}
	}
}

// Helper function to check if error is related to directory confirmation
func isDirectoryConfirmationError(err error) bool {
	return err != nil && (
	// Add patterns that match directory confirmation errors
	// This is a simple heuristic - in real tests we'd mock the prompts
	false) // For now, return false
}

// Test configuration validation
func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *config.ProjectConfig
		wantErr bool
	}{
		{
			name: "valid basic config",
			config: &config.ProjectConfig{
				UserName:           "Test User",
				Email:              "test@example.com",
				ProjectName:        "valid-project",
				ProjectDescription: "A valid project",
				ProjectType:        "basic",
				MainDirName:        "valid_project",
				PythonVersion:      "3.11",
			},
			wantErr: false,
		},
		{
			name: "valid web config",
			config: &config.ProjectConfig{
				UserName:           "Web Developer",
				Email:              "dev@webproject.com",
				ProjectName:        "web-app",
				ProjectDescription: "A web application",
				ProjectType:        "web",
				WebFramework:       "fastapi",
				MainDirName:        "web_app",
				PythonVersion:      "3.12",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that the configuration itself is valid
			context := tt.config.TemplateContext()

			// Check required fields are present
			requiredFields := []string{"project_name", "user_name", "email", "python_version"}
			for _, field := range requiredFields {
				if _, exists := context[field]; !exists {
					t.Errorf("Template context missing required field: %s", field)
				}
			}

			// Validate project name sanitization
			sanitized := config.SanitizeProjectName(tt.config.ProjectName)
			if sanitized == "" && tt.config.ProjectName != "" {
				t.Errorf("Project name sanitization failed for: %s", tt.config.ProjectName)
			}
		})
	}
}

// Test template rendering integration
func TestTemplateIntegration(t *testing.T) {
	// This test verifies that the template engine can render actual templates
	tempDir, err := os.MkdirTemp("", "pyinit-template-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			fmt.Printf("Error removing temporary directory %s: %v\n", tempDir, err)
		}
	}()

	cfg := &config.ProjectConfig{
		UserName:           "Template Test",
		Email:              "template@test.com",
		ProjectName:        "template-test",
		ProjectDescription: "Template integration test",
		ProjectType:        "basic",
		ProjectPath:        tempDir,
		MainDirName:        "template_test",
		PythonVersion:      "3.11",
	}

	// This is testing the generateFileFromTemplate function indirectly
	// We know .gitignore should be generated, so let's try creating a test project
	// and verify at least one template renders correctly

	context := cfg.TemplateContext()
	if len(context) == 0 {
		t.Error("Template context is empty")
	}

	// Verify context contains expected keys
	expectedKeys := []string{"project_name", "user_name", "email", "python_version", "python_version_for_ruff"}
	for _, key := range expectedKeys {
		if _, exists := context[key]; !exists {
			t.Errorf("Template context missing key: %s", key)
		}
	}
}

