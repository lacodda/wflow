package cmd

import (
	"wflow/api"
	"wflow/config"
	"wflow/core"
	"wflow/excel"
	"wflow/mail"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	FlagReportDate     = ""
	FlagReportSend     = false
	FlagReportTestSend = false
)

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
		timestampsRes, err := api.PullTimestamps(date, false)
		if err != nil {
			core.Danger("Error: %v\n", err.Error())
			return
		}

		from, to := core.DayRange(date)
		tasksRes, err := api.PullTasks(from, to, false)
		if err != nil {
			core.Danger("Error: %v\n", err.Error())
			return
		}

		core.Info("Date: %s\n\n", date.Format(core.DateDotTpl))
		printTimestampsRes(timestampsRes)
		printTaskRes(tasksRes.Data)

		fileName, err := excel.SeveXLSXReport(date, timestampsRes, tasksRes)
		if err != nil {
			core.Danger("Error: %s\n", err.Error())
		}

		if FlagReportSend {
			if len(timestampsRes.Data) == 0 {
				core.Danger("Message not sent. There is no timestamps!\n")
				return
			}
			if len(tasksRes.Data) == 0 {
				core.Danger("Message not sent. There are no tasks!\n")
				return
			}

			var mailObj = config.ReadConfig().Mail
			mailObj.Subject = fmt.Sprintf(config.ReadConfig().Mail.Subject, date.Format(core.DateDotTpl), core.MinutesToTimeStr(timestampsRes.TotalTime))
			if FlagReportTestSend {
				mailObj = config.ReadConfig().TestMail
				mailObj.Subject = fmt.Sprintf(config.ReadConfig().TestMail.Subject, date.Format(core.DateDotTpl), core.MinutesToTimeStr(timestampsRes.TotalTime))
			}
			mailObj.Attachments = fileName
			msg := mail.CreateMail(mailObj)
			mail.SendMail(msg)
			core.Success("Your report %s has been sent to %s!\n", fileName, mailObj.To)
		}
	},
}
