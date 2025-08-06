package setup

import (
	"fmt"
	"os"
	"os/exec"
)

// DevDependencies sets up the Python development environment using uv
func DevDependencies(projectPath string) error {
	fmt.Println("🔧 Setting up development environment...")

	// Check if uv is installed
	if err := checkUvInstalled(); err != nil {
		ShowManualInstructions(projectPath)
		return err
	}

	// Run uv add --dev ruff pyright
	cmd := exec.Command("uv", "add", "--dev", "ruff", "pyright")
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add development dependencies: %w", err)
	}

	fmt.Println("✅ Development environment setup complete!")
	return nil
}

// checkUvInstalled verifies that uv is available in the system PATH
func checkUvInstalled() error {
	_, err := exec.LookPath("uv")
	if err != nil {
		return fmt.Errorf("uv is not installed. Please install uv first: https://docs.astral.sh/uv/getting-started/installation/")
	}
	return nil
}

// ShowManualInstructions displays instructions for manual environment setup
func ShowManualInstructions(projectPath string) {
	fmt.Println("💡 You can set up the development environment later by running:")
	fmt.Println("   cd", projectPath)
	fmt.Println("   uv add --dev ruff pyright")
}

// FastAPIDependencies installs selected FastAPI dependencies using uv
func FastAPIDependencies(projectPath string, selectedDeps []string) error {
	fmt.Println("🚀 Installing FastAPI dependencies...")

	// Check if uv is installed
	if err := checkUvInstalled(); err != nil {
		ShowManualFastAPIInstructions(projectPath, selectedDeps)
		return err
	}

	// Install selected dependencies
	if len(selectedDeps) > 0 {
		args := append([]string{"add"}, selectedDeps...)
		cmd := exec.Command("uv", args...)
		cmd.Dir = projectPath
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add FastAPI dependencies: %w", err)
		}
	}

	// Sync dependencies including dev dependencies
	syncCmd := exec.Command("uv", "sync", "--dev")
	syncCmd.Dir = projectPath
	syncCmd.Stdout = os.Stdout
	syncCmd.Stderr = os.Stderr

	if err := syncCmd.Run(); err != nil {
		return fmt.Errorf("failed to sync dependencies: %w", err)
	}

	fmt.Println("✅ FastAPI dependencies installed successfully!")
	return nil
}

// ShowManualFastAPIInstructions displays instructions for manual FastAPI dependency setup
func ShowManualFastAPIInstructions(projectPath string, selectedDeps []string) {
	fmt.Println("💡 You can install FastAPI dependencies later by running:")
	fmt.Println("   cd", projectPath)
	if len(selectedDeps) > 0 {
		fmt.Printf("   uv add")
		for _, dep := range selectedDeps {
			fmt.Printf(" %s", dep)
		}
		fmt.Println()
	}
	fmt.Println("   uv sync --dev")
}
