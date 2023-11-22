package cmd

import (
	"wflow/api"
	"wflow/core"
	"time"

	"github.com/spf13/cobra"
)

var (
	FlagSummaryDate        = ""
	FlagSummaryRecalculate = false
)

var SummaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Get summary",
	Run: func(cmd *cobra.Command, args []string) {
		var date = time.Now()
		if len(FlagSummaryDate) > 0 {
			var err error
			date, err = time.Parse("2006-01", FlagSummaryDate)
			if err != nil {
				core.Danger("Error: %v\n", err.Error())
				return
			}
		}
		from, to := core.MonthRange(date)
		summaryRes, err := api.PullSummary(from, to, FlagSummaryRecalculate)
		if err != nil {
			core.Danger("Error: %v\n", err.Error())
			return
		}
		core.Info("Date: %s\n\n", date.Format(core.DateDotTpl))
		printSummaryRes(summaryRes)
	},
}

func printSummaryRes(summary core.SummaryRes) {
	if len(summary.Data) > 0 {
		core.Info("Summary:\n")
		core.Info("=========================================================\n")
	}
	for key, summary := range summary.Data {
		date, _ := time.Parse(core.DateISOTpl, summary.Date)
		core.Info("[%d] %s (%s)\n", key+1, date.Format(core.DateDotTpl), core.MinutesToTimeStr(summary.Time))
	}
	if len(summary.Data) > 0 {
		core.Info("=========================================================\n")
		core.Info("Average time: %s\n", core.MinutesToTimeStr(summary.TotalTime/len(summary.Data)))
		core.Info("Total time: %s\n", core.MinutesToTimeStr(summary.TotalTime))
	}
}
