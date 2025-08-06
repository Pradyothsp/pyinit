package commands

import (
	"fmt"

	"github.com/Pradyothsp/pyinit/internal/config"
	"github.com/Pradyothsp/pyinit/internal/generator"
	"github.com/Pradyothsp/pyinit/internal/prompts"
	"github.com/Pradyothsp/pyinit/internal/setup"
	"github.com/Pradyothsp/pyinit/internal/version"
	"github.com/Pradyothsp/pyinit/pkg/ui"
	"github.com/spf13/cobra"
)

// runInteractive handles the main interactive project creation
func (c *Commands) runInteractive(cmd *cobra.Command, args []string) {
	// Check for version flag
	if versionFlag, _ := cmd.Flags().GetBool("version"); versionFlag {
		buildInfo := version.GetBuildInfo()
		fmt.Println(buildInfo.String())
		return
	}

	// Show banner if enabled
	if err := c.showBannerIfEnabled(); err != nil {
		fmt.Printf("Warning: Banner display failed: %v\n", err)
	}

	// Collect user information
	cfg, err := prompts.CollectProjectInfo()
	if err != nil {
		fmt.Printf("Error: Failed to collect project info: %v\n", err)
		return
	}

	// Generate project
	gen := generator.New()
	if err := gen.GenerateProject(cfg); err != nil {
		fmt.Printf("Error: Failed to generate project: %v\n", err)
		return
	}

	fmt.Printf("✅ Project '%s' created successfully at: %s\n", cfg.ProjectName, cfg.ProjectPath)

	// Handle FastAPI dependencies setup (only for FastAPI projects)
	if err := c.handleFastAPIDependencies(cfg); err != nil {
		fmt.Printf("Warning: Failed to setup FastAPI dependencies: %v\n", err)
	}

	// Handle environment setup
	if err := c.handleEnvironmentSetup(cfg.ProjectPath); err != nil {
		fmt.Printf("Warning: Failed to setup environment: %v\n", err)
		return
	}
}

// showBannerIfEnabled displays banner if enabled in config
func (c *Commands) showBannerIfEnabled() error {
	banner, err := ui.NewBanner()
	if err != nil {
		// Continue without banner rather than fail
		return nil
	}

	banner.Show()
	return nil
}

// handleFastAPIDependencies manages FastAPI dependency installation
func (c *Commands) handleFastAPIDependencies(cfg *config.ProjectConfig) error {
	// Only handle FastAPI dependencies for FastAPI projects
	if cfg.ProjectType != "web" || cfg.WebFramework != "fastapi" {
		return nil
	}

	// Ask user to select FastAPI dependencies
	selectedDeps, err := prompts.AskForFastAPIDependencies()
	if err != nil {
		return fmt.Errorf("failed to prompt for FastAPI dependencies: %w", err)
	}

	if len(selectedDeps) == 0 {
		fmt.Println("No dependencies selected, skipping installation.")
		return nil
	}

	// Install selected dependencies
	return setup.FastAPIDependencies(cfg.ProjectPath, selectedDeps)
}

// handleEnvironmentSetup manages the development environment setup
func (c *Commands) handleEnvironmentSetup(projectPath string) error {
	// Ask user if they want to set up environment
	setupEnv, err := prompts.AskForEnvironmentSetup()
	if err != nil {
		return fmt.Errorf("failed to prompt for environment setup: %w", err)
	}

	if !setupEnv {
		setup.ShowManualInstructions(projectPath)
		return nil
	}

	// Set up the development environment
	return setup.DevDependencies(projectPath)
}
