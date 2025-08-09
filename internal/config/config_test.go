package config

import (
	"reflect"
	"testing"
)

func TestProjectTypes(t *testing.T) {
	types := ProjectTypes()
	expected := []string{"basic", "cli", "web", "library", "data-science"}
	
	if !reflect.DeepEqual(types, expected) {
		t.Errorf("ProjectTypes() = %v, want %v", types, expected)
	}
}

func TestWebFrameworks(t *testing.T) {
	frameworks := WebFrameworks()
	expected := []string{"fastapi", "flask", "django"}
	
	if !reflect.DeepEqual(frameworks, expected) {
		t.Errorf("WebFrameworks() = %v, want %v", frameworks, expected)
	}
}

func TestSanitizeProjectName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple lowercase",
			input:    "myproject",
			expected: "myproject",
		},
		{
			name:     "uppercase to lowercase",
			input:    "MyProject",
			expected: "myproject",
		},
		{
			name:     "spaces to hyphens",
			input:    "My Project",
			expected: "my-project",
		},
		{
			name:     "multiple spaces",
			input:    "My   Project   Name",
			expected: "my---project---name",
		},
		{
			name:     "special characters removed",
			input:    "my-project@123!#$",
			expected: "my-project123",
		},
		{
			name:     "underscores preserved",
			input:    "my_project_name",
			expected: "my_project_name",
		},
		{
			name:     "hyphens preserved",
			input:    "my-project-name",
			expected: "my-project-name",
		},
		{
			name:     "numbers preserved",
			input:    "project123",
			expected: "project123",
		},
		{
			name:     "mixed case with special chars",
			input:    "My-Super_Project@2024!",
			expected: "my-super_project2024",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only special characters",
			input:    "!@#$%^&*()",
			expected: "",
		},
		{
			name:     "unicode characters removed",
			input:    "project🚀test",
			expected: "projecttest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeProjectName(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeProjectName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestProjectConfig_TemplateContext(t *testing.T) {
	config := &ProjectConfig{
		UserName:           "John Doe",
		Email:              "john@example.com",
		ProjectName:        "Test Project",
		ProjectDescription: "A test project for unit testing",
		ProjectType:        "web",
		WebFramework:       "fastapi",
		ProjectPath:        "/path/to/project",
		MainDirName:        "test_project",
		PythonVersion:      "3.11",
	}

	context := config.TemplateContext()

	expectedKeys := []string{
		"project_name",
		"project_type", 
		"project_description",
		"user_name",
		"email",
		"main_dir_name",
		"python_version",
		"python_version_for_ruff",
	}

	// Check that all expected keys exist
	for _, key := range expectedKeys {
		if _, exists := context[key]; !exists {
			t.Errorf("TemplateContext() missing expected key: %s", key)
		}
	}

	// Check specific values
	if context["project_name"] != config.ProjectName {
		t.Errorf("TemplateContext()[\"project_name\"] = %v, want %v", context["project_name"], config.ProjectName)
	}

	if context["user_name"] != config.UserName {
		t.Errorf("TemplateContext()[\"user_name\"] = %v, want %v", context["user_name"], config.UserName)
	}

	if context["email"] != config.Email {
		t.Errorf("TemplateContext()[\"email\"] = %v, want %v", context["email"], config.Email)
	}

	if context["python_version"] != config.PythonVersion {
		t.Errorf("TemplateContext()[\"python_version\"] = %v, want %v", context["python_version"], config.PythonVersion)
	}

	// Check python_version_for_ruff formatting
	expectedRuffVersion := "py311"
	if context["python_version_for_ruff"] != expectedRuffVersion {
		t.Errorf("TemplateContext()[\"python_version_for_ruff\"] = %v, want %v", context["python_version_for_ruff"], expectedRuffVersion)
	}
}

func TestTemplateContext_PythonVersionForRuff(t *testing.T) {
	tests := []struct {
		name           string
		pythonVersion  string
		expectedRuff   string
	}{
		{
			name:          "Python 3.11",
			pythonVersion: "3.11",
			expectedRuff:  "py311",
		},
		{
			name:          "Python 3.12",
			pythonVersion: "3.12",
			expectedRuff:  "py312",
		},
		{
			name:          "Python 3.13",
			pythonVersion: "3.13",
			expectedRuff:  "py313",
		},
		{
			name:          "Python 3.8",
			pythonVersion: "3.8",
			expectedRuff:  "py38",
		},
		{
			name:          "Python 3.10",
			pythonVersion: "3.10",
			expectedRuff:  "py310",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &ProjectConfig{
				PythonVersion: tt.pythonVersion,
			}

			context := config.TemplateContext()
			ruffVersion := context["python_version_for_ruff"]

			if ruffVersion != tt.expectedRuff {
				t.Errorf("PythonVersion %s -> python_version_for_ruff = %v, want %v", 
					tt.pythonVersion, ruffVersion, tt.expectedRuff)
			}
		})
	}
}

func TestProjectConfig_DefaultValues(t *testing.T) {
	config := &ProjectConfig{}
	context := config.TemplateContext()

	// Test that empty values are handled correctly
	if context["project_name"] != "" {
		t.Errorf("Expected empty project_name for default config, got %v", context["project_name"])
	}

	if context["python_version_for_ruff"] != "py" {
		t.Errorf("Expected 'py' for empty python version, got %v", context["python_version_for_ruff"])
	}
}

func TestProjectConfig_ContextIntegrity(t *testing.T) {
	config := &ProjectConfig{
		UserName:           "Test User",
		Email:              "test@test.com",
		ProjectName:        "integration-test",
		ProjectDescription: "Integration test project",
		ProjectType:        "basic",
		ProjectPath:        "/tmp/test",
		MainDirName:        "integration_test",
		PythonVersion:      "3.12",
	}

	context := config.TemplateContext()

	// Verify context contains exactly what we expect
	expectedContext := map[string]interface{}{
		"project_name":              "integration-test",
		"project_type":              "basic", 
		"project_description":       "Integration test project",
		"user_name":                 "Test User",
		"email":                     "test@test.com",
		"main_dir_name":             "integration_test",
		"python_version":            "3.12",
		"python_version_for_ruff":   "py312",
	}

	if !reflect.DeepEqual(context, expectedContext) {
		t.Errorf("TemplateContext() = %+v, want %+v", context, expectedContext)
	}
}

// Test edge cases and potential security issues
func TestSanitizeProjectName_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "very long string",
			input:    "this-is-a-very-long-project-name-that-might-cause-issues-in-some-systems-but-should-be-handled-gracefully",
			expected: "this-is-a-very-long-project-name-that-might-cause-issues-in-some-systems-but-should-be-handled-gracefully",
		},
		{
			name:     "path injection attempt",
			input:    "../../../etc/passwd",
			expected: "etcpasswd",
		},
		{
			name:     "shell injection attempt", 
			input:    "; rm -rf /; echo 'pwned'",
			expected: "-rm--rf--echo-pwned",
		},
		{
			name:     "sql injection attempt",
			input:    "'; DROP TABLE projects; --",
			expected: "-drop-table-projects---",
		},
		{
			name:     "whitespace only",
			input:    "   \t\n   ",
			expected: "------",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeProjectName(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeProjectName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}