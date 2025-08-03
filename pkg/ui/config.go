package ui

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config holds UI configuration
type Config struct {
	ShowBanner bool `config:"show_banner"`
	// Future extensions can be added here
	// EnableAnimations bool `config:"enable_animations"`
	// CurrentTheme string `config:"current_theme"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		ShowBanner: true,
	}
}

// getConfigPath returns the path to the .pyinitrc file
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, ".pyinitrc"), nil
}

// LoadConfig loads configuration from ~/.pyinitrc
func LoadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// If config doesn't exist, return error so caller can create default
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist")
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Warning: Failed to close config file: %v\n", err)
		}
	}(file)

	config := DefaultConfig()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		value = strings.Trim(value, "\"'")

		switch key {
		case "show_banner":
			if boolVal, err := strconv.ParseBool(value); err == nil {
				config.ShowBanner = boolVal
			}
			// Future config options can be added here
			// case "enable_animations":
			//     if boolVal, err := strconv.ParseBool(value); err == nil {
			//         config.EnableAnimations = boolVal
			//     }
			// case "current_theme":
			//     config.CurrentTheme = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return config, nil
}

// Save saves the configuration to ~/.pyinitrc
func (c *Config) Save() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Warning: Failed to close config file: %v\n", err)
		}
	}(file)

	w := bufio.NewWriter(file)

	// Write header comment
	_, _ = fmt.Fprintln(w, "# pyinit configuration file")
	_, _ = fmt.Fprintln(w, "# This file configures the behavior of the pyinit CLI tool")
	_, _ = fmt.Fprintln(w, "")

	// Write current config
	_, _ = fmt.Fprintf(w, "# Show ASCII banner on startup (true/false)\n")
	_, _ = fmt.Fprintf(w, "show_banner=%t\n", c.ShowBanner)

	// Placeholder for future config options
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "# Future configuration options will appear here")
	_, _ = fmt.Fprintln(w, "# enable_animations=true")
	_, _ = fmt.Fprintln(w, "# current_theme=ocean")

	if err := w.Flush(); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfigPath returns the config file path for display
func GetConfigPath() string {
	path, _ := getConfigPath()
	return path
}
