package prompts

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Pradyothsp/pyinit/internal/config"
)

// Test the email validation function
func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectError bool
		errorText   string
	}{
		{
			name:        "valid email",
			input:       "user@example.com",
			expectError: false,
		},
		{
			name:        "valid email with plus",
			input:       "user+tag@example.com",
			expectError: false,
		},
		{
			name:        "valid email with dots",
			input:       "first.last@example.com",
			expectError: false,
		},
		{
			name:        "valid email with numbers",
			input:       "user123@example123.com",
			expectError: false,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
			errorText:   "email is required",
		},
		{
			name:        "whitespace only",
			input:       "   ",
			expectError: true,
			errorText:   "email is required",
		},
		{
			name:        "missing @",
			input:       "userexample.com",
			expectError: true,
			errorText:   "enter valid email",
		},
		{
			name:        "missing domain",
			input:       "user@",
			expectError: true,
			errorText:   "enter valid email",
		},
		{
			name:        "missing user",
			input:       "@example.com",
			expectError: true,
			errorText:   "enter valid email",
		},
		{
			name:        "missing TLD",
			input:       "user@example",
			expectError: true,
			errorText:   "enter valid email",
		},
		{
			name:        "invalid characters",
			input:       "user@exam ple.com",
			expectError: true,
			errorText:   "enter valid email",
		},
		{
			name:        "multiple @ symbols",
			input:       "user@exam@ple.com",
			expectError: true,
			errorText:   "enter valid email",
		},
		{
			name:        "non-string input",
			input:       123,
			expectError: true,
			errorText:   "invalid input",
		},
		{
			name:        "nil input",
			input:       nil,
			expectError: true,
			errorText:   "invalid input",
		},
		{
			name:        "valid email with subdomain",
			input:       "user@mail.example.com",
			expectError: false,
		},
		{
			name:        "valid email with long TLD",
			input:       "user@example.technology",
			expectError: false,
		},
		{
			name:        "TLD too short",
			input:       "user@example.c",
			expectError: true,
			errorText:   "enter valid email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEmail(tt.input)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for input %v, but got none", tt.input)
				} else if tt.errorText != "" && !strings.Contains(err.Error(), tt.errorText) {
					t.Errorf("Expected error containing %q, got %q", tt.errorText, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for input %v, but got: %v", tt.input, err)
				}
			}
		})
	}
}

// Test the updateConfigFromAnswer function
func TestUpdateConfigFromAnswer(t *testing.T) {
	cfg := &config.ProjectConfig{}

	tests := []struct {
		name       string
		questionID string
		answer     string
		expectErr  bool
		checkFunc  func(*config.ProjectConfig) bool
	}{
		{
			name:       "update username",
			questionID: "username",
			answer:     "John Doe",
			expectErr:  false,
			checkFunc:  func(c *config.ProjectConfig) bool { return c.UserName == "John Doe" },
		},
		{
			name:       "update email",
			questionID: "email",
			answer:     "john@example.com",
			expectErr:  false,
			checkFunc:  func(c *config.ProjectConfig) bool { return c.Email == "john@example.com" },
		},
		{
			name:       "update project name",
			questionID: "projectname",
			answer:     "My Project",
			expectErr:  false,
			checkFunc:  func(c *config.ProjectConfig) bool { return c.ProjectName == "My Project" },
		},
		{
			name:       "update project type",
			questionID: "projecttype",
			answer:     "web",
			expectErr:  false,
			checkFunc:  func(c *config.ProjectConfig) bool { return c.ProjectType == "web" },
		},
		{
			name:       "update web framework",
			questionID: "webframework",
			answer:     "fastapi",
			expectErr:  false,
			checkFunc:  func(c *config.ProjectConfig) bool { return c.WebFramework == "fastapi" },
		},
		{
			name:       "update main dir name",
			questionID: "maindirname",
			answer:     "my_project",
			expectErr:  false,
			checkFunc:  func(c *config.ProjectConfig) bool { return c.MainDirName == "my_project" },
		},
		{
			name:       "update description",
			questionID: "description",
			answer:     "A test project",
			expectErr:  false,
			checkFunc:  func(c *config.ProjectConfig) bool { return c.ProjectDescription == "A test project" },
		},
		{
			name:       "update python version",
			questionID: "pythonversion",
			answer:     "3.11",
			expectErr:  false,
			checkFunc:  func(c *config.ProjectConfig) bool { return c.PythonVersion == "3.11" },
		},
		{
			name:       "unknown question ID",
			questionID: "unknown",
			answer:     "value",
			expectErr:  true,
			checkFunc:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset config for each test
			cfg = &config.ProjectConfig{}
			
			err := updateConfigFromAnswer(cfg, tt.questionID, tt.answer)
			
			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error for question ID %s, but got none", tt.questionID)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for question ID %s, but got: %v", tt.questionID, err)
				}
				if tt.checkFunc != nil && !tt.checkFunc(cfg) {
					t.Errorf("Config update check failed for question ID %s", tt.questionID)
				}
			}
		})
	}
}

// Test setProjectPath function
func TestSetProjectPath(t *testing.T) {
	// Save original working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	
	// Create a temporary directory and change to it
	tempDir, err := os.MkdirTemp("", "pyinit-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	
	// Restore original directory after test
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Errorf("Failed to restore working directory: %v", err)
		}
	}()

	tests := []struct {
		name           string
		projectName    string
		expectedSuffix string
	}{
		{
			name:           "simple project name",
			projectName:    "myproject",
			expectedSuffix: "myproject",
		},
		{
			name:           "project name with spaces",
			projectName:    "My Project",
			expectedSuffix: "my-project",
		},
		{
			name:           "project name with special chars",
			projectName:    "My-Project_123",
			expectedSuffix: "my-project_123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.ProjectConfig{
				ProjectName: tt.projectName,
			}

			err := setProjectPath(cfg)
			if err != nil {
				t.Fatalf("setProjectPath failed: %v", err)
			}

			// Get current working directory to compare against
			currentWd, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get current working directory: %v", err)
			}

			expectedPath := filepath.Join(currentWd, tt.expectedSuffix)
			if cfg.ProjectPath != expectedPath {
				t.Errorf("ProjectPath = %q, want %q", cfg.ProjectPath, expectedPath)
			}

			// Also check that the path ends with the expected suffix
			if !strings.HasSuffix(cfg.ProjectPath, tt.expectedSuffix) {
				t.Errorf("ProjectPath %q does not end with expected suffix %q", cfg.ProjectPath, tt.expectedSuffix)
			}
		})
	}
}

// Test buildCompleteQuestionFlow structure
func TestBuildCompleteQuestionFlow(t *testing.T) {
	questions := buildCompleteQuestionFlow()

	// Check that we have the expected number of questions
	expectedQuestions := []string{
		"username",
		"email", 
		"projectname",
		"projecttype",
		"webframework",
		"maindirname",
		"description",
		"pythonversion",
	}

	if len(questions) != len(expectedQuestions) {
		t.Errorf("Expected %d questions, got %d", len(expectedQuestions), len(questions))
	}

	// Check that all expected questions are present
	questionIDs := make([]string, len(questions))
	for i, q := range questions {
		questionIDs[i] = q.ID
	}

	for _, expectedID := range expectedQuestions {
		found := false
		for _, id := range questionIDs {
			if id == expectedID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected question ID %s not found", expectedID)
		}
	}

	// Check specific question properties
	for _, q := range questions {
		switch q.ID {
		case "username":
			if !q.Required {
				t.Error("Username question should be required")
			}
			if q.Condition != nil {
				t.Error("Username question should have no condition")
			}
		case "email":
			if !q.Required {
				t.Error("Email question should be required")
			}
			if q.Question.Validate == nil {
				t.Error("Email question should have validation")
			}
		case "webframework":
			if q.Required {
				t.Error("Web framework question should not be required")
			}
			if q.Condition == nil {
				t.Error("Web framework question should have condition")
			}
			
			// Test the condition function
			webCfg := &config.ProjectConfig{ProjectType: "web"}
			nonWebCfg := &config.ProjectConfig{ProjectType: "basic"}
			
			if !q.Condition(webCfg) {
				t.Error("Web framework condition should return true for web projects")
			}
			if q.Condition(nonWebCfg) {
				t.Error("Web framework condition should return false for non-web projects")
			}
		}
	}
}

// Test conditional logic for web framework question
func TestWebFrameworkCondition(t *testing.T) {
	questions := buildCompleteQuestionFlow()
	
	var webFrameworkQuestion *QuestionStep
	for _, q := range questions {
		if q.ID == "webframework" {
			webFrameworkQuestion = &q
			break
		}
	}

	if webFrameworkQuestion == nil {
		t.Fatal("Web framework question not found")
	}

	tests := []struct {
		name         string
		projectType  string
		shouldAsk    bool
	}{
		{
			name:         "web project should ask framework",
			projectType:  "web",
			shouldAsk:    true,
		},
		{
			name:         "basic project should not ask framework",
			projectType:  "basic",
			shouldAsk:    false,
		},
		{
			name:         "cli project should not ask framework",
			projectType:  "cli",
			shouldAsk:    false,
		},
		{
			name:         "library project should not ask framework",
			projectType:  "library",
			shouldAsk:    false,
		},
		{
			name:         "data-science project should not ask framework",
			projectType:  "data-science",
			shouldAsk:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.ProjectConfig{
				ProjectType: tt.projectType,
			}

			result := webFrameworkQuestion.Condition(cfg)
			if result != tt.shouldAsk {
				t.Errorf("Condition for project type %s = %v, want %v", tt.projectType, result, tt.shouldAsk)
			}
		})
	}
}

// Test that question structure is valid
func TestQuestionStructureValidity(t *testing.T) {
	questions := buildCompleteQuestionFlow()

	for _, q := range questions {
		t.Run(fmt.Sprintf("question_%s", q.ID), func(t *testing.T) {
			if q.ID == "" {
				t.Error("Question ID should not be empty")
			}

			if q.Question == nil {
				t.Error("Question should not be nil")
			}

			if q.Question != nil && q.Question.Name == "" {
				t.Error("Question name should not be empty")
			}

			if q.Question != nil && q.Question.Prompt == nil {
				t.Error("Question prompt should not be nil")
			}
		})
	}
}

// Test the interaction between questions (that web framework is asked only for web projects)
func TestQuestionFlow(t *testing.T) {
	questions := buildCompleteQuestionFlow()
	
	// Simulate answering questions for a web project
	webConfig := &config.ProjectConfig{}
	
	for _, step := range questions {
		// Skip if condition doesn't match
		if step.Condition != nil && !step.Condition(webConfig) {
			continue
		}

		// Simulate answers based on question ID
		var answer string
		switch step.ID {
		case "username":
			answer = "Test User"
		case "email":
			answer = "test@example.com"
		case "projectname":
			answer = "Test Project"
		case "projecttype":
			answer = "web"
		case "webframework":
			answer = "fastapi"
		case "maindirname":
			answer = "test_project"
		case "description":
			answer = "Test description"
		case "pythonversion":
			answer = "3.11"
		}

		// Update config
		if err := updateConfigFromAnswer(webConfig, step.ID, answer); err != nil {
			t.Fatalf("Failed to update config for %s: %v", step.ID, err)
		}
	}

	// Verify that all expected fields are set
	if webConfig.UserName != "Test User" {
		t.Error("UserName not set correctly")
	}
	if webConfig.ProjectType != "web" {
		t.Error("ProjectType not set correctly")
	}
	if webConfig.WebFramework != "fastapi" {
		t.Error("WebFramework not set correctly")
	}
}

// Test askForCustomPath function indirectly by testing the logic
func TestCustomPathLogic(t *testing.T) {
	// This tests the path joining logic that would be used in askForCustomPath
	projectName := "My Test Project"
	customDir := "/path/to/custom"
	
	projectDir := config.SanitizeProjectName(projectName)
	expectedPath := filepath.Join(customDir, projectDir)
	
	if projectDir != "my-test-project" {
		t.Errorf("Expected sanitized name 'my-test-project', got %s", projectDir)
	}
	
	expectedFull := "/path/to/custom/my-test-project"
	if expectedPath != expectedFull {
		t.Errorf("Expected path %s, got %s", expectedFull, expectedPath)
	}
}