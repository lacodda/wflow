package cmd

import (
	"errors"
	"finlab/apps/time-tool/api"
	"finlab/apps/time-tool/core"
	"finlab/apps/time-tool/db"
	"time"

	"github.com/spf13/cobra"
)

type Value interface {
	String() string
	Set(string) error
	Type() string
}

type TimestampType core.TimestampType

const (
	Start      TimestampType = "Start"
	End        TimestampType = "End"
	StartBreak TimestampType = "StartBreak"
	EndBreak   TimestampType = "EndBreak"
)

// String is used both by fmt.Print and by Cobra in help text
func (e *TimestampType) String() string {
	return string(*e)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *TimestampType) Set(v string) error {
	switch v {
	case "Start", "End", "StartBreak", "EndBreak":
		*e = TimestampType(v)
		return nil
	default:
		return errors.New(`Must be one of "Start", "End", "StartBreak", or "EndBreak"`)
	}
}

// Type is only used in help text
func (e *TimestampType) Type() string {
	return "type"
}

var (
	FlagTimestampType = StartBreak
	FlagTimestampShow = false
	FlagTimestampRaw  = false
	FlagTimestampDate = ""
)

var TimestampCmd = &cobra.Command{
	Use:   "timestamp",
	Short: "Write timestamp and event type to database",
	Run: func(cmd *cobra.Command, args []string) {
		if FlagTimestampShow {
			var date = time.Now()
			if len(FlagTimestampDate) > 0 {
				var err error
				date, err = time.Parse(core.DateTpl, FlagTimestampDate)
				if err != nil {
					core.Danger("Error: %v\n", err.Error())
					return
				}
			}
			timestampsRes, err := api.PullTimestamps(date, FlagTimestampRaw)
			if err != nil {
				core.Danger("Error: %v\n", err.Error())
				return
			}
			core.Info("Date: %s\n\n", date.Format(core.DateDotTpl))
			printTimestampsRes(timestampsRes)
			return
		}
		timestamp := core.Timestamp{
			Timestamp: time.Now(),
			Type:      core.TimestampType(FlagTimestampType),
		}
		timestampRes, err := api.PushTimestamp(timestamp)
		if err != nil {
			core.Danger("Error: %v\n", err.Error())
			err := db.SetTimestamp(timestamp)
			if err != nil {
				core.Danger("Error: %v\n", err.Error())
			}
			core.Success("Timestamp saved in local database\n")
			core.Info("Timestamp (%s): %s\n", timestamp.Type, timestamp.Timestamp.Format(core.DateISOTpl))
			return
		}

		core.Info("Timestamp (%s): %s\n", timestampRes.Data.Type, timestampRes.Data.Timestamp)
	},
}

func printTimestampsRes(timestamps core.TimestampsRes) {
	if len(timestamps.Data) > 0 {
		core.Info("Timestamps:\n")
		core.Info("=========================================================\n")
	}
	for key, timestamp := range timestamps.Data {
		time, _ := time.Parse(core.DateISOTpl, timestamp.Timestamp)
		core.Info("[%d] %s (%s)\n", key+1, time.Format(core.TimeTpl), timestamp.Type)
	}
	if timestamps.TotalTime > 0 {
		core.Info("=========================================================\n")
		core.Info("Total time: %s\n", core.MinutesToTimeStr(timestamps.TotalTime))
		core.Info("\n")
	}
}
