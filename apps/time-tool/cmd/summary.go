package cmd

import (
	"finlab/apps/time-tool/api"
	"finlab/apps/time-tool/core"
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
		core.Info("Date: %s\n", date.Format(core.DateTpl))
		printSummaryRes(summaryRes.Data)
		if len(summaryRes.Data) > 0 {
			core.Info("Average time: %s\n", core.MinutesToTimeStr(summaryRes.TotalTime/len(summaryRes.Data)))
		}
		core.Info("Total time: %s\n", core.MinutesToTimeStr(summaryRes.TotalTime))
	},
}

func printSummaryRes(summary []core.Summary) {
	for key, sum := range summary {
		date, _ := time.Parse(core.DateISOTpl, sum.Date)
		core.Info("[%d] %s (%s)\n", key+1, date.Format(core.DateDotTpl), core.MinutesToTimeStr(sum.Time))
	}
}
