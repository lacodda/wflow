package cmd

import (
	"finlab/apps/time-tool/api"
	"finlab/apps/time-tool/core"
	"time"

	"github.com/spf13/cobra"
)

var TimestampCmd = &cobra.Command{
	Use:   "timestamp",
	Short: "Write timestamp and event type to database",
	Run: func(cmd *cobra.Command, args []string) {
		timestamp := core.Timestamp{
			Timestamp: time.Now(),
			Type:      core.EndBreak,
		}
		timestampRes, err := api.PushTimestamp(timestamp)
		if err != nil {
			core.Danger("Error: %v\n", err.Error())
			api.SetTimestamp(timestamp)
			return
		}

		core.Info("Timestamp: %s\n", timestampRes.Data.Timestamp)
	},
}
