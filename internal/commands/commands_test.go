package commands

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewCommands(t *testing.T) {
	commands := NewCommands()

	if commands == nil {
		t.Fatal("NewCommands() returned nil")
	}

	if commands.rootCmd == nil {
		t.Fatal("Root command is nil")
	}
}

func TestRootCommandStructure(t *testing.T) {
	commands := NewCommands()
	rootCmd := commands.rootCmd

	// Test basic command properties
	if rootCmd.Use != "pyinit" {
		t.Errorf("Root command Use = %q, want %q", rootCmd.Use, "pyinit")
	}

	if rootCmd.Short == "" {
		t.Error("Root command Short description is empty")
	}

	if rootCmd.Long == "" {
		t.Error("Root command Long description is empty")
	}

	if rootCmd.Run == nil {
		t.Error("Root command Run function is nil")
	}
}

func TestVersionFlag(t *testing.T) {
	commands := NewCommands()
	rootCmd := commands.rootCmd

	// Check that version flag exists
	versionFlag := rootCmd.Flags().Lookup("version")
	if versionFlag == nil {
		t.Fatal("Version flag not found")
	}

	if versionFlag.Shorthand != "v" {
		t.Errorf("Version flag shorthand = %q, want %q", versionFlag.Shorthand, "v")
	}

	if versionFlag.Usage == "" {
		t.Error("Version flag usage is empty")
	}
}

func TestConfigCommandStructure(t *testing.T) {
	commands := NewCommands()
	rootCmd := commands.rootCmd

	// Find config command
	var configCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "config" {
			configCmd = cmd
			break
		}
	}

	if configCmd == nil {
		t.Fatal("Config command not found")
	}

	// Test config command properties
	if configCmd.Short == "" {
		t.Error("Config command Short description is empty")
	}

	if configCmd.Long == "" {
		t.Error("Config command Long description is empty")
	}

	// Test config subcommands
	subcommands := configCmd.Commands()
	expectedSubcommands := []string{"show", "reset", "banner"}
	
	if len(subcommands) < len(expectedSubcommands) {
		t.Errorf("Expected at least %d config subcommands, got %d", len(expectedSubcommands), len(subcommands))
	}

	// Check that expected subcommands exist
	foundSubcommands := make(map[string]bool)
	for _, subcmd := range subcommands {
		foundSubcommands[subcmd.Use] = true
	}

	for _, expected := range expectedSubcommands {
		if !foundSubcommands[expected] {
			t.Errorf("Expected config subcommand %q not found", expected)
		}
	}
}

func TestConfigShowCommand(t *testing.T) {
	commands := NewCommands()
	
	showCmd := commands.createConfigShowCommand()
	if showCmd == nil {
		t.Fatal("Config show command is nil")
	}

	if showCmd.Use != "show" {
		t.Errorf("Config show command Use = %q, want %q", showCmd.Use, "show")
	}

	if showCmd.Short == "" {
		t.Error("Config show command Short description is empty")
	}

	if showCmd.Run == nil {
		t.Error("Config show command Run function is nil")
	}
}

func TestConfigResetCommand(t *testing.T) {
	commands := NewCommands()
	
	resetCmd := commands.createConfigResetCommand()
	if resetCmd == nil {
		t.Fatal("Config reset command is nil")
	}

	if resetCmd.Use != "reset" {
		t.Errorf("Config reset command Use = %q, want %q", resetCmd.Use, "reset")
	}

	if resetCmd.Short == "" {
		t.Error("Config reset command Short description is empty")
	}

	if resetCmd.Run == nil {
		t.Error("Config reset command Run function is nil")
	}
}

func TestCommandHelpText(t *testing.T) {
	commands := NewCommands()
	rootCmd := commands.rootCmd

	// Capture help output
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"--help"})

	// Execute help command
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Help command failed: %v", err)
	}

	helpOutput := buf.String()

	// Check that help contains expected content
	expectedContent := []string{
		"pyinit",
		"interactive CLI tool",
		"Usage:",
		"Available Commands:",
		"config",
		"Flags:",
		"--version",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(helpOutput, expected) {
			t.Errorf("Help output does not contain %q\nOutput:\n%s", expected, helpOutput)
		}
	}
}

func TestConfigCommandHelp(t *testing.T) {
	commands := NewCommands()
	rootCmd := commands.rootCmd

	// Capture help output for config command
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"config", "--help"})

	// Execute config help command
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Config help command failed: %v", err)
	}

	helpOutput := buf.String()

	// Check that config help contains expected content
	expectedContent := []string{
		"config",
		"Configure pyinit behavior",
		"Usage:",
		"Available Commands:",
		"show",
		"reset",
		"banner",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(helpOutput, expected) {
			t.Errorf("Config help output does not contain %q\nOutput:\n%s", expected, helpOutput)
		}
	}
}

func TestVersionFlagHandling(t *testing.T) {
	commands := NewCommands()
	rootCmd := commands.rootCmd

	// Test that version flag is present and can be parsed
	rootCmd.SetArgs([]string{"--version"})
	
	// Parse flags to set the version flag value
	err := rootCmd.ParseFlags([]string{"--version"})
	if err != nil {
		t.Fatalf("Failed to parse version flag: %v", err)
	}

	// Check that version flag is set
	versionFlag, err := rootCmd.Flags().GetBool("version")
	if err != nil {
		t.Fatalf("Failed to get version flag value: %v", err)
	}

	if !versionFlag {
		t.Error("Version flag should be true when --version is passed")
	}
}

func TestVersionFlagShorthand(t *testing.T) {
	commands := NewCommands()
	rootCmd := commands.rootCmd

	// Test that short version flag can be parsed
	err := rootCmd.ParseFlags([]string{"-v"})
	if err != nil {
		t.Fatalf("Failed to parse short version flag: %v", err)
	}

	// Check that version flag is set
	versionFlag, err := rootCmd.Flags().GetBool("version")
	if err != nil {
		t.Fatalf("Failed to get version flag value: %v", err)
	}

	if !versionFlag {
		t.Error("Version flag should be true when -v is passed")
	}
}

func TestCommandValidation(t *testing.T) {
	commands := NewCommands()
	rootCmd := commands.rootCmd

	// Test that commands have proper validation
	rootCmd.SetArgs([]string{"nonexistent"})
	
	// Capture error output
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	err := rootCmd.Execute()
	if err == nil {
		t.Error("Expected error for nonexistent command, got nil")
	}

	// Should contain error message
	output := buf.String()
	if !strings.Contains(output, "unknown command") || !strings.Contains(output, "nonexistent") {
		t.Errorf("Expected error message about unknown command, got: %s", output)
	}
}

func TestExecuteMethod(t *testing.T) {
	commands := NewCommands()
	
	// Test that Execute method exists and can be called
	// We'll use version flag to avoid interactive prompts
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	os.Args = []string{"pyinit", "--version"}
	
	err := commands.Execute()
	if err != nil {
		t.Errorf("Commands.Execute() failed: %v", err)
	}
}

// Test command structure completeness
func TestCommandCompleteness(t *testing.T) {
	commands := NewCommands()
	rootCmd := commands.rootCmd

	// Check that all expected commands are present
	allCommands := make(map[string]*cobra.Command)
	
	// Add root command
	allCommands[rootCmd.Use] = rootCmd
	
	// Add all subcommands recursively
	var addSubcommands func(*cobra.Command)
	addSubcommands = func(cmd *cobra.Command) {
		for _, subcmd := range cmd.Commands() {
			allCommands[subcmd.Use] = subcmd
			addSubcommands(subcmd)
		}
	}
	addSubcommands(rootCmd)

	// Verify expected commands exist
	expectedCommands := []string{"pyinit", "config", "show", "reset", "banner"}
	for _, expected := range expectedCommands {
		if _, exists := allCommands[expected]; !exists {
			t.Errorf("Expected command %q not found in command tree", expected)
		}
	}

	// Verify each command has proper structure
	for name, cmd := range allCommands {
		if cmd.Use == "" {
			t.Errorf("Command %q has empty Use field", name)
		}
		if cmd.Short == "" && name != "banner" { // banner might be composite
			t.Errorf("Command %q has empty Short description", name)
		}
	}
}

func TestCommandChaining(t *testing.T) {
	// Test that commands can be chained properly
	commands := NewCommands()
	
	// This should not panic
	result := NewCommands()
	if result == nil {
		t.Error("Command chaining failed - NewCommands() returned nil")
	}
	
	if commands.rootCmd.Use != result.rootCmd.Use {
		t.Error("Command chaining created inconsistent commands")
	}
}

// Benchmark command creation
func BenchmarkNewCommands(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewCommands()
	}
}

func BenchmarkCommandExecution(b *testing.B) {
	commands := NewCommands()
	rootCmd := commands.rootCmd
	
	// Use a buffer to avoid output during benchmarking
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"--version"})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		rootCmd.Execute()
	}
}