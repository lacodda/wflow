package cmd

import (
	"finlab/apps/time-tool/api"
	"finlab/apps/time-tool/core"

	"github.com/spf13/cobra"
)

// Synchronizing local storage with the server
var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronizing local storage with the server.",
	Run: func(cmd *cobra.Command, args []string) {
		timestamps, _ := api.GetTimestamps()
		for _, timestamp := range timestamps {
			timestampRes, err := api.PushTimestamp(timestamp)
			if err != nil {
				core.Danger("Error: %v\n", err.Error())
				return
			}
			core.Info("Timestamp (%s): %s\n", timestampRes.Data.Type, timestampRes.Data.Timestamp)
		}
		api.DeleteTimestamps()
		core.Success("Your data is synced!\n")
	},
}
