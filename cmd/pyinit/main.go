package main

import (
	"fmt"
	"os"

	"github.com/Pradyothsp/pyinit/internal/generator"
	"github.com/Pradyothsp/pyinit/internal/prompts"
	"github.com/Pradyothsp/pyinit/internal/setup"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pyinit",
	Short: "Python project scaffolding tool",
	Long:  "An interactive CLI tool to create Python project scaffolds with customizable structure",
	Run:   runInteractive,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}

func runInteractive(cmd *cobra.Command, args []string) {
	// Collect user information
	cfg, err := prompts.CollectProjectInfo()
	if err != nil {
		fmt.Printf("Warning: Failed to collect project info: %v\n", err)
		return
	}

	// Generate project
	gen := generator.New()
	if err := gen.GenerateProject(cfg); err != nil {
		fmt.Printf("Warning: Failed to generate project: %v\n", err)
		return
	}

	fmt.Printf("✅ Project '%s' created successfully at: %s\n", cfg.ProjectName, cfg.ProjectPath)

	// Handle environment setup
	if err := handleEnvironmentSetup(cfg.ProjectPath); err != nil {
		fmt.Printf("Warning: Failed to setup environment: %v\n", err)
		return
	}
}

func handleEnvironmentSetup(projectPath string) error {
	// Ask a user if they want to set up an environment
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
