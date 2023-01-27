package cmd

import (
	"finlab/apps/time-tool/api"
	"finlab/apps/time-tool/core"
	"finlab/apps/time-tool/excel"
	"time"

	"github.com/spf13/cobra"
)

var FlagReportDate = ""

var ReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Prepare a report",
	Run: func(cmd *cobra.Command, args []string) {
		var date = time.Now()

		if len(FlagReportDate) > 0 {
			var err error
			date, err = time.Parse(core.DateTpl, FlagReportDate)
			if err != nil {
				core.Danger("Error: %v\n", err.Error())
				return
			}
		} else {
			timestamp := core.Timestamp{
				Timestamp: time.Now(),
				Type:      core.TimestampType(End),
			}
			_, err := api.PushTimestamp(timestamp)
			if err != nil {
				core.Danger("Error: %v\n", err.Error())
			}
		}
		timestampsRes, err := api.PullTimestamps(date)
		if err != nil {
			core.Danger("Error: %v\n", err.Error())
			return
		}

		core.Info("Date: %s\n", date.Format(core.DateTpl))
		printTimestampsRes(timestampsRes.Data)
		core.Info("Total time: %s\n", core.MinutesToTimeStr(timestampsRes.TotalTime))

		if err := excel.SeveXlsx(date, timestampsRes); err != nil {
			core.Danger("Error: %s\n", err.Error())
		}

		core.Success("You've been logged out!\n")
	},
}
