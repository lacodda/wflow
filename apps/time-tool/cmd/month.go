package cmd

import (
	"finlab/apps/time-tool/api"
	"finlab/apps/time-tool/core"
	"finlab/apps/time-tool/excel"
	"time"

	"github.com/spf13/cobra"
)

var MonthlyReportCmd = &cobra.Command{
	Use:   "month",
	Short: "Prepare a monthly report",
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
		excel.SeveXLSXMonthlyReport(from, to, summaryRes)
	},
}
