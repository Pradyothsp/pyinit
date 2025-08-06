package commands

import (
	"github.com/spf13/cobra"
)

// Commands holds all command dependencies
type Commands struct {
	rootCmd *cobra.Command
}

// NewCommands creates a new Commands instance with all subcommands
func NewCommands() *Commands {
	cmd := &Commands{}
	cmd.setupRootCommand()
	cmd.setupConfigCommands()
	return cmd
}

// Execute runs the root command
func (c *Commands) Execute() error {
	return c.rootCmd.Execute()
}

// setupRootCommand initializes the root command
func (c *Commands) setupRootCommand() {
	c.rootCmd = &cobra.Command{
		Use:   "pyinit",
		Short: "Interactive Python project scaffolding tool",
		Long:  "An interactive CLI tool to create Python project scaffolds with customizable structure",
		Run:   c.runInteractive,
	}
	
	// Add version flag
	c.rootCmd.Flags().BoolP("version", "v", false, "Show version information")
}

// setupConfigCommands adds all config-related commands
func (c *Commands) setupConfigCommands() {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage pyinit configuration",
		Long:  "Configure pyinit behavior via ~/.pyinitrc",
	}

	// Add config subcommands
	configCmd.AddCommand(c.createConfigShowCommand())
	configCmd.AddCommand(c.createConfigResetCommand())
	configCmd.AddCommand(c.createBannerCommands())

	c.rootCmd.AddCommand(configCmd)
}

// createConfigShowCommand creates the config show command
func (c *Commands) createConfigShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		Run:   c.showConfig,
	}
}

// createConfigResetCommand creates the config reset command
func (c *Commands) createConfigResetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Reset configuration to defaults",
		Run:   c.resetConfig,
	}
}

// showConfig displays current configuration
func (c *Commands) showConfig(cmd *cobra.Command, args []string) {
	// Implementation moved to config.go
	c.handleShowConfig()
}

// resetConfig resets configuration to defaults
func (c *Commands) resetConfig(cmd *cobra.Command, args []string) {
	// Implementation moved to config.go
	c.handleResetConfig()
}
