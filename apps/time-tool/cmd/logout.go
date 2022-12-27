package cmd

import (
	"finlab/apps/time-tool/config"
	"finlab/apps/time-tool/core"

	"github.com/spf13/cobra"
)

// LogoutCmd logs the user out by cleaning the local state so the user needs to login again.
var LogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logs out the user by removing the user's session from local state.",
	Run: func(cmd *cobra.Command, args []string) {
		config.RemoveToken()

		core.Success("You've been logged out!\n")
	},
}
