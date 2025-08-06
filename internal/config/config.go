package config

import (
	"strings"
)

// ProjectConfig holds all the configuration for a Python project
type ProjectConfig struct {
	UserName           string
	Email              string
	ProjectName        string
	ProjectDescription string
	ProjectType        string
	WebFramework       string
	ProjectPath        string
	MainDirName        string
	PythonVersion      string
}

// ProjectTypes returns available project types
func ProjectTypes() []string {
	return []string{"basic", "cli", "web", "library", "data-science"}
}

func WebFrameworks() []string {
	return []string{"fastapi", "flask", "django"}
}

// SanitizeProjectName converts the project name to a valid directory name
func SanitizeProjectName(name string) string {
	// Replace spaces with hyphens and convert to lowercase
	sanitized := strings.ToLower(strings.ReplaceAll(name, " ", "-"))

	// Remove any characters that aren't alphanumeric, hyphens, or underscores
	result := ""
	for _, char := range sanitized {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' || char == '_' {
			result += string(char)
		}
	}

	return result
}

// TemplateContext returns a map for template rendering
func (pc *ProjectConfig) TemplateContext() map[string]interface{} {
	// Create ruff-compatible Python version (e.g., "3.13" -> "py313")
	pythonVersionForRuff := "py" + strings.ReplaceAll(pc.PythonVersion, ".", "")
	
	return map[string]interface{}{
		"project_name":              pc.ProjectName,
		"project_type":              pc.ProjectType,
		"project_description":       pc.ProjectDescription,
		"user_name":                 pc.UserName,
		"email":                     pc.Email,
		"main_dir_name":             pc.MainDirName,
		"python_version":            pc.PythonVersion,
		"python_version_for_ruff":   pythonVersionForRuff,
	}
}
