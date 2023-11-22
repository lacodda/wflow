package cmd

import (
	"wflow/api"
	"wflow/core"
	"wflow/excel"
	"time"

	"github.com/spf13/cobra"
)

var (
	FlagMonthDate = ""
)

var MonthlyReportCmd = &cobra.Command{
	Use:   "month",
	Short: "Prepare a monthly report",
	Run: func(cmd *cobra.Command, args []string) {
		var date = time.Now()
		if len(FlagMonthDate) > 0 {
			var err error
			date, err = time.Parse("2006-01", FlagMonthDate)
			if err != nil {
				core.Danger("Error: %v\n", err.Error())
				return
			}
		}
		month := date.Month()
		year := date.Year()
		calendarRes, err := api.PullCalendar(month, year)
		if err != nil {
			core.Danger("Error: %v\n", err.Error())
			return
		}
		excel.SeveXLSXMonthlyReport(month, year, calendarRes)
	},
}
