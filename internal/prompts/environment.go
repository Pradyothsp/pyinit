package prompts

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

// AskForEnvironmentSetup prompts the user whether they want to set up the development environment
func AskForEnvironmentSetup() (bool, error) {
	setupEnv := false
	prompt := &survey.Confirm{
		Message: "Do you want to set up the development environment now?",
		Default: true,
	}

	if err := survey.AskOne(prompt, &setupEnv); err != nil {
		return false, fmt.Errorf("failed to get environment setup confirmation: %w", err)
	}

	return setupEnv, nil
}

// AskForFastAPIDependencies prompts the user to select FastAPI dependencies they want to install
func AskForFastAPIDependencies() ([]string, error) {
	var selectedDeps []string
	
	availableDeps := []string{
		"fastapi",
		"uvicorn[standard]",
		"sqlalchemy",
		"pydantic-settings",
		"python-jose[cryptography]",
		"passlib[bcrypt]",
		"alembic",
		"httpx", // for testing
	}

	prompt := &survey.MultiSelect{
		Message: "Select FastAPI dependencies to install:",
		Options: availableDeps,
		Default: []string{"fastapi", "uvicorn[standard]"}, // Core dependencies selected by default
		Help:    "Use space to select/deselect, Enter to confirm",
	}

	if err := survey.AskOne(prompt, &selectedDeps); err != nil {
		return nil, fmt.Errorf("failed to get FastAPI dependencies selection: %w", err)
	}

	return selectedDeps, nil
}
