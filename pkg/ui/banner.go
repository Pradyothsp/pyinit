package ui

import (
	"fmt"
	"github.com/Pradyothsp/pyinit/internal/version"
)

// Banner handles ASCII banner display
type Banner struct {
	config *Config
}

// NewBanner creates a new banner instance
func NewBanner() (*Banner, error) {
	config, err := LoadConfig()
	if err != nil {
		// If config doesn't exist, create default
		config = DefaultConfig()
		if err := config.Save(); err != nil {
			// If we can't save config, continue with defaults
			fmt.Printf("Warning: Could not save config to ~/.pyinitrc: %v\n", err)
		}
	}

	return &Banner{
		config: config,
	}, nil
}

// Show displays the banner if enabled in config
func (b *Banner) Show() {
	if !b.config.ShowBanner {
		return
	}

	fmt.Print(b.generateBanner())
}

// generateBanner creates the ASCII art banner
func (b *Banner) generateBanner() string {
	// ASCII art for "pyinit"
	asciiArt := `
██████╗ ██╗   ██╗██╗███╗   ██╗██╗████████╗
██╔══██╗╚██╗ ██╔╝██║████╗  ██║██║╚══██╔══╝
██████╔╝ ╚████╔╝ ██║██╔██╗ ██║██║   ██║   
██╔═══╝   ╚██╔╝  ██║██║╚██╗██║██║   ██║   
██║        ██║   ██║██║ ╚████║██║   ██║   
╚═╝        ╚═╝   ╚═╝╚═╝  ╚═══╝╚═╝   ╚═╝   
`

	tagline := "🚀 Interactive Python Project Scaffolding Tool"
	versionInfo := version.GetVersion()

	return fmt.Sprintf("%s\n%s\n%s\n\n", asciiArt, tagline, versionInfo)
}

// IsEnabled returns whether the banner is enabled
func (b *Banner) IsEnabled() bool {
	return b.config.ShowBanner
}

// Enable enables the banner and saves config
func (b *Banner) Enable() error {
	b.config.ShowBanner = true
	return b.config.Save()
}

// Disable disables the banner and saves config
func (b *Banner) Disable() error {
	b.config.ShowBanner = false
	return b.config.Save()
}
