package prompts

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Pradyothsp/pyinit/internal/config"
)

// CollectProjectInfo gathers all necessary information from the user
func CollectProjectInfo() (*config.ProjectConfig, error) {
	cfg := &config.ProjectConfig{}

	if err := collectUserDetails(cfg); err != nil {
		return nil, fmt.Errorf("failed to collect user details: %w", err)
	}

	if err := collectProjectDetails(cfg); err != nil {
		return nil, fmt.Errorf("failed to collect project details: %w", err)
	}

	if err := setProjectPath(cfg); err != nil {
		return nil, fmt.Errorf("failed to set project path: %w", err)
	}

	return cfg, nil
}

func collectUserDetails(cfg *config.ProjectConfig) error {
	questions := []*survey.Question{
		{
			Name:   "username",
			Prompt: &survey.Input{Message: "Enter your name:"},
		},
		{
			Name:   "email",
			Prompt: &survey.Input{Message: "Enter your email:"},
		},
	}

	answers := struct {
		Username string
		Email    string
	}{}

	if err := survey.Ask(questions, &answers); err != nil {
		return err
	}

	cfg.UserName = answers.Username
	cfg.Email = answers.Email

	return nil
}

func collectProjectDetails(cfg *config.ProjectConfig) error {
	questions := []*survey.Question{
		{
			Name:   "projectname",
			Prompt: &survey.Input{Message: "Enter project name:"},
		},
		{
			Name: "projecttype",
			Prompt: &survey.Select{
				Message: "Select project type:",
				Options: config.ProjectTypes(),
				Default: "basic",
			},
		},
		{
			Name: "projectstructure",
			Prompt: &survey.Select{
				Message: "Select project structure:",
				Options: config.ProjectStructures(),
				Default: "src",
			},
		},
	}

	answers := struct {
		ProjectName      string `survey:"projectname"`
		ProjectType      string `survey:"projecttype"`
		ProjectStructure string `survey:"projectstructure"`
	}{}

	if err := survey.Ask(questions, &answers); err != nil {
		return err
	}

	cfg.ProjectName = answers.ProjectName
	cfg.ProjectType = answers.ProjectType
	cfg.ProjectStructure = answers.ProjectStructure

	return nil
}

func setProjectPath(cfg *config.ProjectConfig) error {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Create a project path
	projectDir := config.SanitizeProjectName(cfg.ProjectName)
	cfg.ProjectPath = filepath.Join(cwd, projectDir)

	return nil
}

// ConfirmDirectoryCreation asks the user for confirmation if the directory exists
func ConfirmDirectoryCreation(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		confirm := false
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Directory '%s' already exists. Continue anyway?", filepath.Base(path)),
			Default: false,
		}
		if err := survey.AskOne(prompt, &confirm); err != nil {
			return false, fmt.Errorf("failed to get confirmation: %w", err)
		}
		return confirm, nil
	}
	return true, nil
}
