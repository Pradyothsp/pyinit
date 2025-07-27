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
