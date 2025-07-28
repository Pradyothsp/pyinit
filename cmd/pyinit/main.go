package main

import (
	"fmt"
	"os"

	"github.com/Pradyothsp/pyinit/internal/generator"
	"github.com/Pradyothsp/pyinit/internal/prompts"
	"github.com/Pradyothsp/pyinit/internal/setup"
	"github.com/Pradyothsp/pyinit/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	banner *ui.Banner

	rootCmd = &cobra.Command{
		Use:   "pyinit",
		Short: "Interactive Python project scaffolding tool",
		Long:  "An interactive CLI tool to create Python project scaffolds with customizable structure",
		Run:   runInteractive,
	}

	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Manage pyinit configuration",
		Long:  "Configure pyinit behavior via ~/.pyinitrc",
	}

	bannerCmd = &cobra.Command{
		Use:   "banner",
		Short: "Manage banner display",
		Long:  "Enable or disable the ASCII banner",
	}

	bannerEnableCmd = &cobra.Command{
		Use:   "enable",
		Short: "Enable ASCII banner",
		Run:   enableBanner,
	}

	bannerDisableCmd = &cobra.Command{
		Use:   "disable",
		Short: "Disable ASCII banner",
		Run:   disableBanner,
	}

	configShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		Run:   showConfig,
	}
)

func init() {
	// Initialize banner
	var err error
	banner, err = ui.NewBanner()
	if err != nil {
		fmt.Printf("Warning: Failed to initialize banner: %v\n", err)
		// Continue without a banner rather than exit
	}

	// Add config subcommands
	bannerCmd.AddCommand(bannerEnableCmd, bannerDisableCmd)
	configCmd.AddCommand(bannerCmd, configShowCmd)
	rootCmd.AddCommand(configCmd)
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
	// Show banner if enabled
	if banner != nil {
		banner.Show()
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

// Config command handlers
func enableBanner(cmd *cobra.Command, args []string) {
	if banner == nil {
		fmt.Println("Error: Banner system not initialized")
		return
	}

	if err := banner.Enable(); err != nil {
		fmt.Printf("Error: Failed to enable banner: %v\n", err)
		return
	}

	fmt.Println("✅ Banner enabled")
}

func disableBanner(cmd *cobra.Command, args []string) {
	if banner == nil {
		fmt.Println("Error: Banner system not initialized")
		return
	}

	if err := banner.Disable(); err != nil {
		fmt.Printf("Error: Failed to disable banner: %v\n", err)
		return
	}

	fmt.Println("✅ Banner disabled")
}

func showConfig(cmd *cobra.Command, args []string) {
	fmt.Println("Current pyinit configuration:")
	fmt.Printf("  Config file: %s\n", ui.GetConfigPath())

	if banner != nil {
		fmt.Printf("  Banner enabled: %t\n", banner.IsEnabled())
	}

	fmt.Println("\nTo modify configuration:")
	fmt.Println("  pyinit config banner enable    # Enable banner")
	fmt.Println("  pyinit config banner disable   # Disable banner")
	fmt.Printf("  edit %s           # Edit config manually\n", ui.GetConfigPath())
}
