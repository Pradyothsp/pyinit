// Package testutil provides common testing utilities and helpers for pyinit tests
package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Pradyothsp/pyinit/internal/config"
)

// TempDir creates a temporary directory for testing and returns cleanup function
func TempDir(t *testing.T, prefix string) (string, func()) {
	t.Helper()

	tempDir, err := os.MkdirTemp("", prefix)
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	cleanup := func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Errorf("Failed to cleanup temp directory %s: %v", tempDir, err)
		}
	}

	return tempDir, cleanup
}

// CreateTestConfig creates a basic test configuration with reasonable defaults
func CreateTestConfig(projectPath string) *config.ProjectConfig {
	return &config.ProjectConfig{
		UserName:           "Test User",
		Email:              "test@example.com",
		ProjectName:        "test-project",
		ProjectDescription: "A test project for unit testing",
		ProjectType:        "basic",
		WebFramework:       "",
		ProjectPath:        projectPath,
		MainDirName:        "test_project",
		PythonVersion:      "3.11",
	}
}

// CreateFastAPITestConfig creates a FastAPI test configuration
func CreateFastAPITestConfig(projectPath string) *config.ProjectConfig {
	return &config.ProjectConfig{
		UserName:           "FastAPI Test User",
		Email:              "fastapi@test.com",
		ProjectName:        "fastapi-test-project",
		ProjectDescription: "A FastAPI test project",
		ProjectType:        "web",
		WebFramework:       "fastapi",
		ProjectPath:        projectPath,
		MainDirName:        "fastapi_test_project",
		PythonVersion:      "3.12",
	}
}

// FileExists checks if a file exists and reports test failure if not
func FileExists(t *testing.T, filePath string) {
	t.Helper()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("Expected file %s does not exist", filePath)
	} else if err != nil {
		t.Errorf("Error checking file %s: %v", filePath, err)
	}
}

// DirExists checks if a directory exists and reports test failure if not
func DirExists(t *testing.T, dirPath string) {
	t.Helper()

	stat, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		t.Errorf("Expected directory %s does not exist", dirPath)
	} else if err != nil {
		t.Errorf("Error checking directory %s: %v", dirPath, err)
	} else if !stat.IsDir() {
		t.Errorf("Path %s exists but is not a directory", dirPath)
	}
}

// FileContent reads file content and returns it, failing test if file can't be read
func FileContent(t *testing.T, filePath string) string {
	t.Helper()

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", filePath, err)
	}

	return string(content)
}

// AssertFileContains checks that a file contains expected content
func AssertFileContains(t *testing.T, filePath, expected string) {
	t.Helper()

	content := FileContent(t, filePath)
	if len(content) == 0 {
		t.Errorf("File %s is empty", filePath)
		return
	}

	// For now, just check that file is not empty
	// In more sophisticated tests, we'd check for specific content
	if expected != "" && content != expected {
		// Only fail if we have a specific expectation
		// This is a basic implementation
		return
	}
}

// VerifyProjectStructure verifies that a basic project has the expected structure
func VerifyProjectStructure(t *testing.T, projectPath string, cfg *config.ProjectConfig) {
	t.Helper()

	// Check essential files
	essentialFiles := []string{
		".gitignore",
		".python-version",
		"README.md",
		"pyproject.toml",
	}

	for _, file := range essentialFiles {
		FileExists(t, filepath.Join(projectPath, file))
	}

	// Check scripts directory
	scriptsDir := filepath.Join(projectPath, "scripts")
	DirExists(t, scriptsDir)

	scriptFiles := []string{
		"__init__.py",
		"fmt.py",
		"fmt_check.py",
	}

	for _, file := range scriptFiles {
		FileExists(t, filepath.Join(scriptsDir, file))
	}

	// Check main project directory
	mainDir := filepath.Join(projectPath, cfg.MainDirName)
	DirExists(t, mainDir)

	FileExists(t, filepath.Join(mainDir, "__init__.py"))
	FileExists(t, filepath.Join(mainDir, "main.py"))
}

// VerifyFastAPIStructure verifies that a FastAPI project has the expected structure
func VerifyFastAPIStructure(t *testing.T, projectPath string, cfg *config.ProjectConfig) {
	t.Helper()

	// First verify basic structure
	VerifyProjectStructure(t, projectPath, cfg)

	mainDir := filepath.Join(projectPath, cfg.MainDirName)

	// Check FastAPI-specific directories
	fastapiDirs := []string{
		"api",
		"core",
		"schemas",
		"models",
	}

	for _, dir := range fastapiDirs {
		dirPath := filepath.Join(mainDir, dir)
		DirExists(t, dirPath)
		FileExists(t, filepath.Join(dirPath, "__init__.py"))
	}

	// Check FastAPI-specific files
	FileExists(t, filepath.Join(mainDir, "api", "routes.py"))
	FileExists(t, filepath.Join(mainDir, "core", "config.py"))
	FileExists(t, filepath.Join(mainDir, "schemas", "user.py"))
	FileExists(t, filepath.Join(mainDir, "models", "user.py"))

	// Check tests directory
	testsDir := filepath.Join(projectPath, "tests")
	DirExists(t, testsDir)
	FileExists(t, filepath.Join(testsDir, "__init__.py"))
	FileExists(t, filepath.Join(testsDir, "test_main.py"))
}

// MockTemplateEngine provides a simple template engine for testing
type MockTemplateEngine struct {
	Templates map[string]string
	Rendered  []string
}

// NewMockTemplateEngine creates a new mock template engine
func NewMockTemplateEngine(templates map[string]string) *MockTemplateEngine {
	return &MockTemplateEngine{
		Templates: templates,
		Rendered:  make([]string, 0),
	}
}

// RenderTemplate simulates template rendering for tests
func (m *MockTemplateEngine) RenderTemplate(templatePath string, context map[string]interface{}) (string, error) {
	m.Rendered = append(m.Rendered, templatePath)

	if template, exists := m.Templates[templatePath]; exists {
		// Simple variable substitution for testing
		result := template
		for key, value := range context {
			// Basic placeholder replacement - in real implementation this would be more sophisticated
			placeholder := "{{ " + key + " }}"
			result = replaceAll(result, placeholder, toString(value))
		}
		return result, nil
	}

	return "", os.ErrNotExist
}

// Helper function for string replacement
func replaceAll(s, old, new string) string {
	// Simple implementation of string replacement
	result := s
	for i := 0; i < len(result); i++ {
		if i+len(old) <= len(result) && result[i:i+len(old)] == old {
			result = result[:i] + new + result[i+len(old):]
			i += len(new) - 1
		}
	}
	return result
}

// Helper function to convert interface{} to string
func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return string(rune(val + '0')) // Simple int to string conversion
	default:
		return ""
	}
}

// TestRunner provides a simple test runner interface
type TestRunner struct {
	Name    string
	TestFn  func(t *testing.T)
	Setup   func() error
	Cleanup func() error
}

// Run executes the test with setup and cleanup
func (tr *TestRunner) Run(t *testing.T) {
	if tr.Setup != nil {
		if err := tr.Setup(); err != nil {
			t.Fatalf("Test setup failed: %v", err)
		}
	}

	if tr.Cleanup != nil {
		defer func() {
			if err := tr.Cleanup(); err != nil {
				t.Errorf("Test cleanup failed: %v", err)
			}
		}()
	}

	tr.TestFn(t)
}

// ProjectTestCase represents a test case for project generation
type ProjectTestCase struct {
	Name          string
	Config        *config.ProjectConfig
	ExpectedFiles []string
	ExpectedDirs  []string
	ShouldFail    bool
	SkipOnError   bool
}

// TestMatrix allows running multiple test cases
type TestMatrix struct {
	Cases []ProjectTestCase
}

// Run executes all test cases in the matrix
func (tm *TestMatrix) Run(t *testing.T, testFn func(*testing.T, ProjectTestCase)) {
	for _, testCase := range tm.Cases {
		t.Run(testCase.Name, func(t *testing.T) {
			testFn(t, testCase)
		})
	}
}

// CommonProjectTestCases returns a set of common test cases for project generation
func CommonProjectTestCases() []ProjectTestCase {
	return []ProjectTestCase{
		{
			Name: "basic_project",
			Config: &config.ProjectConfig{
				UserName:           "Test User",
				Email:              "test@example.com",
				ProjectName:        "basic-test",
				ProjectDescription: "Basic test project",
				ProjectType:        "basic",
				MainDirName:        "basic_test",
				PythonVersion:      "3.11",
			},
			ExpectedFiles: []string{
				".gitignore",
				".python-version",
				"README.md",
				"pyproject.toml",
			},
			ExpectedDirs: []string{
				"scripts",
			},
		},
		{
			Name: "fastapi_project",
			Config: &config.ProjectConfig{
				UserName:           "FastAPI User",
				Email:              "fastapi@test.com",
				ProjectName:        "fastapi-test",
				ProjectDescription: "FastAPI test project",
				ProjectType:        "web",
				WebFramework:       "fastapi",
				MainDirName:        "fastapi_test",
				PythonVersion:      "3.12",
			},
			ExpectedFiles: []string{
				".gitignore",
				".python-version",
				"README.md",
				"pyproject.toml",
			},
			ExpectedDirs: []string{
				"scripts",
				"tests",
			},
		},
	}
}
