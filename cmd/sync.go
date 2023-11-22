package cmd

import (
	"wflow/api"
	"wflow/core"
	"wflow/db"

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
		if err := syncTimestamps(); err != nil {
			core.Danger("Error: %v\n", err.Error())
			return
		}
		if err := syncTasks(); err != nil {
			core.Danger("Error: %v\n", err.Error())
			return
		}

		if !FlagSyncShow {
			db.DeleteTimestamps()
			db.DeleteTasks()
			core.Success("Your data is synced!\n")
		}
	},
}

func syncTimestamps() error {
	timestamps, _ := db.GetTimestamps()
	if len(timestamps) == 0 {
		core.Warning("No timestamps for synchronizing with the server!\n")
		return nil
	}
	err := handleTimestamps(timestamps)
	if err != nil {
		return err
	}
	core.Info("Timestamps:\n")
	printTimestamps(timestamps)

	return nil
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

func syncTasks() error {
	tasks, _ := db.GetTasks()
	if len(tasks) == 0 {
		core.Warning("No tasks for synchronizing with the server!\n")
		return nil
	}
	err := handleTasks(tasks)
	if err != nil {
		return err
	}
	core.Info("Tasks:\n")
	printTasks(tasks)

	return nil
}

func handleTasks(tasks []core.Task) error {
	for _, task := range tasks {
		if !FlagSyncShow {
			_, err := api.PushTask(core.TaskReq{
				Date:         task.Date,
				Name:         task.Name,
				Comment:      task.Comment,
				Completeness: task.Completeness,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func printTasks(tasks []core.Task) {
	for _, task := range tasks {
		core.Info("[%d] %s %s (Completeness: %v%%)\n", task.Id, task.Date.Format("2006-01-02"), task.Name, task.Completeness)
	}
}
