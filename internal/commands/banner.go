package commands

import (
	"fmt"

	"github.com/Pradyothsp/pyinit/pkg/ui"
	"github.com/spf13/cobra"
)

// createBannerCommands creates banner-related commands
func (c *Commands) createBannerCommands() *cobra.Command {
	bannerCmd := &cobra.Command{
		Use:   "banner",
		Short: "Manage banner display",
		Long:  "Enable or disable the ASCII banner",
	}

	bannerEnableCmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable ASCII banner",
		Run:   c.enableBanner,
	}

	bannerDisableCmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable ASCII banner",
		Run:   c.disableBanner,
	}

	bannerCmd.AddCommand(bannerEnableCmd, bannerDisableCmd)
	return bannerCmd
}

// enableBanner enables the banner display
func (c *Commands) enableBanner(cmd *cobra.Command, args []string) {
	banner, err := ui.NewBanner()
	if err != nil {
		fmt.Printf("Error: Failed to initialize banner system: %v\n", err)
		return
	}

	if err := banner.Enable(); err != nil {
		fmt.Printf("Error: Failed to enable banner: %v\n", err)
		return
	}

	fmt.Println("✅ Banner enabled")
}

// disableBanner disables the banner display
func (c *Commands) disableBanner(cmd *cobra.Command, args []string) {
	banner, err := ui.NewBanner()
	if err != nil {
		fmt.Printf("Error: Failed to initialize banner system: %v\n", err)
		return
	}

	if err := banner.Disable(); err != nil {
		fmt.Printf("Error: Failed to disable banner: %v\n", err)
		return
	}

	fmt.Println("✅ Banner disabled")
}
