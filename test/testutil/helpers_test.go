package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTempDir(t *testing.T) {
	tempDir, cleanup := TempDir(t, "test-helper-*")
	defer cleanup()
	
	// Verify directory exists
	if stat, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Errorf("Temporary directory %s was not created", tempDir)
	} else if !stat.IsDir() {
		t.Errorf("Temporary path %s is not a directory", tempDir)
	}
	
	// Create a test file in the directory
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Errorf("Failed to create test file: %v", err)
	}
	
	// Cleanup should remove the directory
	cleanup()
	
	// Directory should no longer exist
	if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
		t.Errorf("Temporary directory %s was not cleaned up", tempDir)
	}
}

func TestCreateTestConfig(t *testing.T) {
	tempDir, cleanup := TempDir(t, "config-test-*")
	defer cleanup()
	
	cfg := CreateTestConfig(tempDir)
	
	// Verify required fields are set
	if cfg.UserName == "" {
		t.Error("Test config UserName is empty")
	}
	
	if cfg.Email == "" {
		t.Error("Test config Email is empty")
	}
	
	if cfg.ProjectName == "" {
		t.Error("Test config ProjectName is empty")
	}
	
	if cfg.ProjectType == "" {
		t.Error("Test config ProjectType is empty")
	}
	
	if cfg.ProjectPath != tempDir {
		t.Errorf("Test config ProjectPath = %s, want %s", cfg.ProjectPath, tempDir)
	}
	
	if cfg.PythonVersion == "" {
		t.Error("Test config PythonVersion is empty")
	}
}

func TestCreateFastAPITestConfig(t *testing.T) {
	tempDir, cleanup := TempDir(t, "fastapi-config-test-*")
	defer cleanup()
	
	cfg := CreateFastAPITestConfig(tempDir)
	
	// Verify it's a web project with FastAPI
	if cfg.ProjectType != "web" {
		t.Errorf("FastAPI config ProjectType = %s, want %s", cfg.ProjectType, "web")
	}
	
	if cfg.WebFramework != "fastapi" {
		t.Errorf("FastAPI config WebFramework = %s, want %s", cfg.WebFramework, "fastapi")
	}
	
	if cfg.ProjectPath != tempDir {
		t.Errorf("FastAPI config ProjectPath = %s, want %s", cfg.ProjectPath, tempDir)
	}
}

func TestFileExists(t *testing.T) {
	tempDir, cleanup := TempDir(t, "file-exists-test-*")
	defer cleanup()
	
	// Create a test file
	testFile := filepath.Join(tempDir, "exists.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// This should pass (file exists)
	FileExists(t, testFile)
	
	// This would fail if run (file doesn't exist)
	// FileExists(t, filepath.Join(tempDir, "nonexistent.txt"))
}

func TestDirExists(t *testing.T) {
	tempDir, cleanup := TempDir(t, "dir-exists-test-*")
	defer cleanup()
	
	// Create a test directory
	testDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	
	// This should pass (directory exists)
	DirExists(t, testDir)
	
	// This would fail if run (directory doesn't exist)
	// DirExists(t, filepath.Join(tempDir, "nonexistent"))
}

func TestFileContent(t *testing.T) {
	tempDir, cleanup := TempDir(t, "file-content-test-*")
	defer cleanup()
	
	// Create a test file with known content
	expectedContent := "Hello, test world!"
	testFile := filepath.Join(tempDir, "content.txt")
	if err := os.WriteFile(testFile, []byte(expectedContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Read content using helper
	actualContent := FileContent(t, testFile)
	
	if actualContent != expectedContent {
		t.Errorf("FileContent() = %q, want %q", actualContent, expectedContent)
	}
}

func TestMockTemplateEngine(t *testing.T) {
	templates := map[string]string{
		"test.j2": "Hello {{ name }}!",
		"complex.j2": "Project: {{ project_name }}, User: {{ user_name }}",
	}
	
	engine := NewMockTemplateEngine(templates)
	
	// Test simple template
	context := map[string]interface{}{
		"name": "World",
	}
	
	result, err := engine.RenderTemplate("test.j2", context)
	if err != nil {
		t.Fatalf("RenderTemplate failed: %v", err)
	}
	
	expected := "Hello World!"
	if result != expected {
		t.Errorf("RenderTemplate() = %q, want %q", result, expected)
	}
	
	// Verify template was recorded as rendered
	if len(engine.Rendered) != 1 {
		t.Errorf("Expected 1 rendered template, got %d", len(engine.Rendered))
	}
	
	if engine.Rendered[0] != "test.j2" {
		t.Errorf("Expected rendered template 'test.j2', got %q", engine.Rendered[0])
	}
	
	// Test nonexistent template
	_, err = engine.RenderTemplate("nonexistent.j2", context)
	if err == nil {
		t.Error("Expected error for nonexistent template")
	}
}

func TestProjectTestCases(t *testing.T) {
	cases := CommonProjectTestCases()
	
	if len(cases) == 0 {
		t.Error("CommonProjectTestCases returned no test cases")
	}
	
	// Verify each test case has required fields
	for i, testCase := range cases {
		if testCase.Name == "" {
			t.Errorf("Test case %d has empty name", i)
		}
		
		if testCase.Config == nil {
			t.Errorf("Test case %s has nil config", testCase.Name)
		}
		
		if len(testCase.ExpectedFiles) == 0 {
			t.Errorf("Test case %s has no expected files", testCase.Name)
		}
	}
	
	// Check that we have both basic and FastAPI test cases
	hasBasic := false
	hasFastAPI := false
	
	for _, testCase := range cases {
		if testCase.Config.ProjectType == "basic" {
			hasBasic = true
		}
		if testCase.Config.ProjectType == "web" && testCase.Config.WebFramework == "fastapi" {
			hasFastAPI = true
		}
	}
	
	if !hasBasic {
		t.Error("CommonProjectTestCases missing basic project test case")
	}
	
	if !hasFastAPI {
		t.Error("CommonProjectTestCases missing FastAPI project test case")
	}
}

func TestTestMatrix(t *testing.T) {
	cases := []ProjectTestCase{
		{
			Name: "test_case_1",
			Config: CreateTestConfig("/tmp/test1"),
		},
		{
			Name: "test_case_2", 
			Config: CreateTestConfig("/tmp/test2"),
		},
	}
	
	matrix := &TestMatrix{Cases: cases}
	
	runCount := 0
	matrix.Run(t, func(t *testing.T, testCase ProjectTestCase) {
		runCount++
		if testCase.Config == nil {
			t.Error("Test case config is nil")
		}
	})
	
	if runCount != len(cases) {
		t.Errorf("Expected %d test runs, got %d", len(cases), runCount)
	}
}

// Test the helper string functions
func TestStringHelpers(t *testing.T) {
	// Test replaceAll
	result := replaceAll("Hello {{ name }}", "{{ name }}", "World")
	expected := "Hello World"
	if result != expected {
		t.Errorf("replaceAll() = %q, want %q", result, expected)
	}
	
	// Test toString
	testCases := []struct {
		input    interface{}
		expected string
	}{
		{"hello", "hello"},
		{123, string(rune(123 + '0'))}, // This is a simple conversion for testing
		{nil, ""},
	}
	
	for _, tc := range testCases {
		result := toString(tc.input)
		if result != tc.expected {
			t.Errorf("toString(%v) = %q, want %q", tc.input, result, tc.expected)
		}
	}
}