package cmd

import (
	"finlab/apps/time-tool/core"
	"finlab/apps/time-tool/excel"
	"time"

	"github.com/spf13/cobra"
)

var ReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Prepare a report",
	Run: func(cmd *cobra.Command, args []string) {
		if err := excel.SeveXlsx(time.Now()); err != nil {
			core.Danger("Error: %s\n", err.Error())
		}

		core.Success("You've been logged out!\n")
	},
}
