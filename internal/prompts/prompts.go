package prompts

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Pradyothsp/pyinit/internal/config"
)

// QuestionStep represents one question in our flow
type QuestionStep struct {
	ID        string                           // Unique identifier for this question
	Question  *survey.Question                 // The actual survey question
	Condition func(*config.ProjectConfig) bool // Function that returns true if we should ask this question
	Required  bool                             // Whether this question is mandatory
}

// CollectProjectInfo gathers all necessary information from the user
func CollectProjectInfo() (*config.ProjectConfig, error) {
	cfg := &config.ProjectConfig{}

	// Collect all info using a question builder pattern
	if err := collectAllDetails(cfg); err != nil {
		return nil, fmt.Errorf("failed to collect details: %w", err)
	}

	// Set a project path based on collected info
	if err := setProjectPath(cfg); err != nil {
		return nil, fmt.Errorf("failed to set project path: %w", err)
	}

	// Confirm project location
	if err := confirmProjectLocation(cfg); err != nil {
		return nil, fmt.Errorf("failed to confirm project location: %w", err)
	}

	return cfg, nil
}

func collectAllDetails(cfg *config.ProjectConfig) error {
	allQuestions := buildCompleteQuestionFlow()

	for _, step := range allQuestions {
		// Check if we should ask this question
		if step.Condition != nil && !step.Condition(cfg) {
			continue // Skip this question
		}

		// Use a string type instead of interface{}
		var answer string

		// Handle a nil validator case
		var options []survey.AskOpt
		if step.Question.Validate != nil {
			options = append(options, survey.WithValidator(step.Question.Validate))
		}

		if err := survey.AskOne(step.Question.Prompt, &answer, options...); err != nil {
			return fmt.Errorf("failed to collect %s: %w", step.ID, err)
		}

		// Immediately update config so next questions can use this info
		if err := updateConfigFromAnswer(cfg, step.ID, answer); err != nil {
			return fmt.Errorf("failed to process %s: %w", step.ID, err)
		}
	}

	// Print the collected configuration for debugging
	fmt.Printf("Collected Configuration: %+v\n", cfg)

	return nil
}

// Helper function to update config based on question ID
func updateConfigFromAnswer(cfg *config.ProjectConfig, questionID string, answer string) error {
	switch questionID {
	case "username":
		cfg.UserName = answer
	case "email":
		cfg.Email = answer
	case "projectname":
		cfg.ProjectName = answer
	case "projecttype":
		cfg.ProjectType = answer
	case "webframework":
		cfg.WebFramework = answer
	case "maindirname":
		cfg.MainDirName = answer
	case "description":
		cfg.ProjectDescription = answer
	case "pythonversion":
		cfg.PythonVersion = answer
	default:
		return fmt.Errorf("unknown question ID: %s", questionID)
	}

	return nil
}

func buildCompleteQuestionFlow() []QuestionStep {
	return []QuestionStep{
		// User Details
		{
			ID:        "username",
			Required:  true,
			Condition: nil,
			Question: &survey.Question{
				Name:     "username",
				Prompt:   &survey.Input{Message: "Enter your name:"},
				Validate: survey.Required,
			},
		},
		{
			ID:        "email",
			Required:  true,
			Condition: nil,
			Question: &survey.Question{
				Name:     "email",
				Prompt:   &survey.Input{Message: "Enter your email:"},
				Validate: validateEmail,
			},
		},

		// Project Details
		{
			ID:        "projectname",
			Required:  true,
			Condition: nil,
			Question: &survey.Question{
				Name:     "projectname",
				Prompt:   &survey.Input{Message: "Enter project name:"},
				Validate: survey.Required,
			},
		},
		{
			ID:        "projecttype",
			Required:  true,
			Condition: nil,
			Question: &survey.Question{
				Name: "projecttype",
				Prompt: &survey.Select{
					Message: "Select project type:",
					Options: config.ProjectTypes(),
					Default: "basic",
				},
			},
		},
		{
			ID:       "webframework",
			Required: false,
			Condition: func(cfg *config.ProjectConfig) bool {
				return cfg.ProjectType == "web"
			},
			Question: &survey.Question{
				Name: "webframework",
				Prompt: &survey.Select{
					Message: "Select web framework:",
					Options: config.WebFrameworks(),
					Default: "fastapi",
				},
			},
		},
		{
			ID:        "maindirname",
			Required:  true,
			Condition: nil,
			Question: &survey.Question{
				Name:     "maindirname",
				Prompt:   &survey.Input{Message: "Enter Main Directory Name:"},
				Validate: survey.Required,
			},
		},
		{
			ID:        "description",
			Required:  true,
			Condition: nil,
			Question: &survey.Question{
				Name: "description",
				Prompt: &survey.Input{
					Message: "Enter project description:",
					Default: "A Python project",
				},
			},
		},
		{
			ID:        "pythonversion",
			Required:  true,
			Condition: nil,
			Question: &survey.Question{
				Name: "pythonversion",
				Prompt: &survey.Input{
					Message: "Enter Python version (default is 3.13):",
					Default: "3.13",
				},
			},
		},
	}
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
