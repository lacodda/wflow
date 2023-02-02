package cmd

import (
	"errors"
	"finlab/apps/time-tool/api"
	"finlab/apps/time-tool/core"
	"strconv"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	FlagTaskShow = false
	FlagTaskDate = ""
)

var TaskCmd = &cobra.Command{
	Use:   "task",
	Short: "Logs out the user by removing the user's session from local state.",
	Run: func(cmd *cobra.Command, args []string) {
		var date = time.Now()
		if len(FlagTaskDate) > 0 {
			var err error
			date, err = time.Parse(core.DateTpl, FlagTaskDate)
			if err != nil {
				core.Danger("Error: %v\n", err.Error())
				return
			}
		}

		if FlagTaskShow {
			from, to := core.DayRange(date)
			tasksRes, err := api.PullTasks(from, to)
			if err != nil {
				core.Danger("Error: %v\n", err.Error())
				return
			}
			core.Info("Date: %s\n", date.Format(core.DateTpl))
			printTaskRes(tasksRes.Data)
			return
		}

		namePrompt := promptui.Prompt{
			Label:    "Task name",
			Validate: validateTaskName,
		}

		commentPrompt := promptui.Prompt{
			Label: "Task comment",
		}

		completenessPrompt := promptui.Prompt{
			Label:    "Task completeness",
			Validate: validateNumber,
		}

		name, errN := namePrompt.Run()
		comment, errCt := commentPrompt.Run()
		completeness, errCs := completenessPrompt.Run()
		promptErr := core.NotNil(errN, errCt, errCs)

		if promptErr != nil {
			core.Danger("Prompt failed: %v\n", promptErr)
			return
		}

		completenessFloat, _ := strconv.ParseFloat(completeness, 64)

		taskRes, err := api.PushTask(core.Task{
			Date:         date,
			Name:         name,
			Comment:      comment,
			Completeness: completenessFloat,
		})

		if err != nil {
			core.Danger("Error: %v\n", err.Error())
			return
		}

		core.Info("Task: %s (Completeness: %v%%)\n", taskRes.Data.Name, taskRes.Data.Completeness)
	},
}

func validateTaskName(input string) error {
	if len(input) < 3 {
		return errors.New("Task name must have more than 3 characters")
	}
	return nil
}

func validateNumber(input string) error {
	_, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return errors.New("Invalid number")
	}
	return nil
}

func printTaskRes(tasks []core.Task) {
	for key, task := range tasks {
		core.Info("[%d] %s (Completeness: %v%%)\n", key+1, task.Name, task.Completeness)
	}
}
