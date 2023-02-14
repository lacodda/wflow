package cmd

import (
	"finlab/apps/time-tool/api"
	"finlab/apps/time-tool/core"
	"finlab/apps/time-tool/db"

	"github.com/spf13/cobra"
)

var FlagSyncShow = false
var FlagSyncDelete []int

// Synchronizing local storage with the server
var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronizing local storage with the server.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(FlagSyncDelete) > 0 {
			deletedTimestamps, _ := db.DeleteTimestampsByIds(FlagSyncDelete)
			if len(deletedTimestamps) == 0 {
				core.Danger("No records for deleting!\n")
				return
			}
			core.Success("Timestamps are deleted:\n")
			printTimestamps(deletedTimestamps)
			return
		}
		timestamps, _ := db.GetTimestamps()
		if len(timestamps) == 0 {
			core.Warning("No records for synchronizing with the server!\n")
			return
		}
		err := handleTimestamps(timestamps)
		if err != nil {
			core.Danger("Error: %v\n", err.Error())
			return
		}
		core.Info("Timestamps:\n")
		printTimestamps(timestamps)

		if !FlagSyncShow {
			db.DeleteTimestamps()
			core.Success("Your data is synced!\n")
		}
	},
}

func handleTimestamps(timestamps []core.Timestamp) error {
	for _, timestamp := range timestamps {
		if !FlagSyncShow {
			_, err := api.PushTimestamp(timestamp)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func printTimestamps(timestamps []core.Timestamp) {
	for _, timestamp := range timestamps {
		core.Info("[%d] %s (%s)\n", timestamp.Id, timestamp.Timestamp.Format("2006-01-02 15:04"), timestamp.Type)
	}
}
