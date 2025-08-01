package prompts

import (
	"fmt"
	"regexp"
	"strings"
)

func validateEmail(ans interface{}) error {
	str, ok := ans.(string)
	if !ok {
		return fmt.Errorf("invalid input")
	}

	email := strings.TrimSpace(str)

	if email == "" {
		return fmt.Errorf("email is required")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !emailRegex.MatchString(email) {
		return fmt.Errorf("enter valid email")
	}

	return nil
}
