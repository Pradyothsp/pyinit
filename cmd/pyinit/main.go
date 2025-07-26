package main

import (
	"fmt"
	"os"

	_ "github.com/Pradyothsp/pyinit/internal/config"
	"github.com/Pradyothsp/pyinit/internal/generator"
	"github.com/Pradyothsp/pyinit/internal/prompts"
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
}
