package config

import (
	"strings"
)

// ProjectConfig holds all the configuration for a Python project
type ProjectConfig struct {
	UserName         string
	Email            string
	ProjectName      string
	ProjectType      string
	ProjectStructure string
	ProjectPath      string
}

// ProjectTypes returns available project types
func ProjectTypes() []string {
	return []string{"cli", "web", "library", "data-science"}
}

// ProjectStructures returns available project structures
func ProjectStructures() []string {
	return []string{"src", "direct"}
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
	return map[string]interface{}{
		"project_name":      pc.ProjectName,
		"project_type":      pc.ProjectType,
		"project_structure": pc.ProjectStructure,
		"user_name":         pc.UserName,
		"email":             pc.Email,
	}
}
