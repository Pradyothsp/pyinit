package commands

import (
	"fmt"

	"github.com/Pradyothsp/pyinit/pkg/ui"
)

// handleShowConfig displays current configuration
func (c *Commands) handleShowConfig() {
	fmt.Println("Current pyinit configuration:")
	fmt.Printf("  Config file: %s\n", ui.GetConfigPath())

	// Read banner status fresh each time
	banner, err := ui.NewBanner()
	if err != nil {
		fmt.Printf("  Banner enabled: unknown (error reading config: %v)\n", err)
	} else {
		fmt.Printf("  Banner enabled: %t\n", banner.IsEnabled())
	}

	fmt.Println("\nTo modify configuration:")
	fmt.Println("  pyinit config banner enable    # Enable banner")
	fmt.Println("  pyinit config banner disable   # Disable banner")
	fmt.Println("  pyinit config reset             # Reset to defaults")
	fmt.Printf("  edit %s           # Edit config manually\n", ui.GetConfigPath())
}

// handleResetConfig resets configuration to defaults
func (c *Commands) handleResetConfig() {
	// Create default config
	defaultConfig := ui.DefaultConfig()

	if err := defaultConfig.Save(); err != nil {
		fmt.Printf("Error: Failed to reset configuration: %v\n", err)
		return
	}

	fmt.Println("✅ Configuration reset to defaults")
	fmt.Printf("Config file: %s\n", ui.GetConfigPath())
}
