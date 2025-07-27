package prompts

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	// Confirm project location
	if err := confirmProjectLocation(cfg); err != nil {
		return nil, fmt.Errorf("failed to confirm project location: %w", err)
	}

	// Set the main project directory
	if err := SetMainDirName(cfg); err != nil {
		return nil, fmt.Errorf("failed to set main directory name: %w", err)
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
				Default: "cli",
			},
		},
		{
			Name: "projectstructure",
			Prompt: &survey.Select{
				Message: "Select project structure:",
				Options: config.ProjectStructures(),
				Default: "direct",
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

func confirmProjectLocation(cfg *config.ProjectConfig) error {
	// Get an absolute path for display
	absPath, err := filepath.Abs(cfg.ProjectPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Ask if the user wants to create a project in the current directory
	createHere := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Create project in current directory?\nProject will be created at: %s", absPath),
		Default: true,
	}

	if err := survey.AskOne(prompt, &createHere); err != nil {
		return fmt.Errorf("failed to get directory confirmation: %w", err)
	}

	// If a user doesn't want to create here, ask for a custom path
	if !createHere {
		if err := askForCustomPath(cfg); err != nil {
			return fmt.Errorf("failed to get custom path: %w", err)
		}
	}

	return nil
}

func askForCustomPath(cfg *config.ProjectConfig) error {
	// Get the current working directory as default
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	var customDir string
	prompt := &survey.Input{
		Message: "Enter parent directory path:",
		Default: cwd,
		Help:    fmt.Sprintf("Project '%s' will be created inside this directory", config.SanitizeProjectName(cfg.ProjectName)),
	}

	if err := survey.AskOne(prompt, &customDir); err != nil {
		return fmt.Errorf("failed to get custom directory: %w", err)
	}

	// Create the full project path by combining custom directory + project name
	projectDir := config.SanitizeProjectName(cfg.ProjectName)
	cfg.ProjectPath = filepath.Join(customDir, projectDir)

	return nil
}

func SetMainDirName(cfg *config.ProjectConfig) error {
	// Set the main directory name based on the project structure
	if cfg.ProjectStructure == "src" {
		cfg.MainDirName = "src"
	} else {
		cfg.MainDirName = strings.ReplaceAll(cfg.ProjectName, "-", "_")
	}

	return nil
}
