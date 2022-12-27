package cmd

import (
	"finlab/apps/time-tool/config"
	"finlab/apps/time-tool/core"

	"github.com/spf13/cobra"
)

// Init
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Configuration initialization",
	Run: func(cmd *cobra.Command, args []string) {
		config.InitConfig()
		core.Success("Your configuration is initialized\n")
	},
}
