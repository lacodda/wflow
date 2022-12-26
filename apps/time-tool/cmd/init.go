package cmd

import (
	"github.com/spf13/cobra"
)

// Init
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Configuration initialization",
	Run: func(cmd *cobra.Command, args []string) {
		InitConfig()
		Success.Printf("Your configuration is initialized\n")
	},
}
