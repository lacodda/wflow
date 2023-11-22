package cmd

import (
	"wflow/api"
	"wflow/config"
	"wflow/core"
	"wflow/db"
	"wflow/gitlab"
	"wflow/validator"
	"fmt"
	"regexp"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var (
	FlagTaskShow   = false
	FlagTaskFind   = false
	FlagTaskDate   = ""
	FlagTaskDelete = false
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
			show(date)
			return
		} else if FlagTaskFind {
			findAndPush(date)
			findAndPushGitLab(date)
			return
		} else if FlagTaskDelete {
			findAndDelete(date)
			return
		}

		taskReq := core.TaskReq{}
		taskReq.Date = date

		err := survey.Ask(getTaskQuestions(core.TaskReq{Completeness: 100}), &taskReq)
		if err != nil {
			core.Danger("Prompt failed: %v\n", err.Error())
			return
		}
		taskRes, err := api.PushTask(taskReq)

		if err != nil {
			core.Danger("Error: %v\n", err.Error())
			err := db.SetTask(taskReq)
			if err != nil {
				core.Danger("Error: %v\n", err.Error())
			}
			core.Success("Task saved in local database\n")
			core.Info("Task: %s (Completeness: %v%%)\n", taskReq.Name, taskReq.Completeness)
			return
		}

		core.Info("Task: %s (Completeness: %v%%)\n", taskRes.Data.Name, taskRes.Data.Completeness)
	},
}

func getTaskQuestions(taskReq core.TaskReq) []*survey.Question {
	return []*survey.Question{
		{
			Name:      "name",
			Prompt:    &survey.Input{Message: "Task name", Default: taskReq.Name},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
		{
			Name:   "comment",
			Prompt: &survey.Input{Message: "Task comment", Default: taskReq.Comment},
		},
		{
			Name:     "completeness",
			Prompt:   &survey.Input{Message: "Task completeness", Default: fmt.Sprintf("%v", taskReq.Completeness)},
			Validate: validator.IsNumber,
		},
	}
}

func getSelectTasks(taskNames []string) []*survey.Question {
	return []*survey.Question{
		{
			Name:   "name",
			Prompt: &survey.MultiSelect{Message: "Choose a tasks", Options: taskNames},
		},
	}
}

func show(date time.Time) {
	from, to := core.DayRange(date)
	tasksRes, err := api.PullTasks(from, to, false)
	if err != nil {
		core.Danger("Error: %v\n", err.Error())
		return
	}
	core.Info("Date: %s\n\n", date.Format(core.DateDotTpl))
	printTaskRes(tasksRes.Data)
}

func getTaskNames(date time.Time) ([]string, error) {
	from, to := core.DayRange(date)
	tasksRes, err := api.PullTasks(from, to, false)
	if err != nil {
		return nil, err
	}
	taskNames := []string{}
	for _, task := range tasksRes.Data {
		taskNames = append(taskNames, task.Name)
	}

	return taskNames, nil
}

func findAndPush(date time.Time) {
	from, to := core.LastWeekRange()
	tasksRes, err := api.PullTasks(from, to, true)
	if err != nil {
		core.Danger("Error: %v\n", err.Error())
		return
	}

	selectedTaskNames := []string{}
	taskNames := make([]string, len(tasksRes.Data))

	for key, task := range tasksRes.Data {
		taskNames[key] = fmt.Sprintf("%s (Completeness: %v%%)", task.Name, task.Completeness)
	}

	survey.Ask(getSelectTasks(taskNames), &selectedTaskNames)
	re := regexp.MustCompile(`(.+)\s\(\w+:\s.+\)`)

	for _, n := range selectedTaskNames {
		task := tasksRes.FindByName(re.FindStringSubmatch(n)[1])
		taskReq := core.TaskReq{
			TaskId: task.TaskId,
			Date:   date,
		}

		err := survey.Ask(getTaskQuestions(core.TaskReq{
			Name:         task.Name,
			Comment:      task.Comment,
			Completeness: task.Completeness + 1,
		}), &taskReq)

		if err != nil {
			core.Danger("Prompt failed: %v\n", err.Error())
			return
		}

		taskRes, err := api.PushTask(taskReq)

		if err != nil {
			core.Danger("Error: %v\n", err.Error())
			return
		}

		core.Info("Task: %s (Completeness: %v%%)\n", taskRes.Data.Name, taskRes.Data.Completeness)
	}

	core.Info("Selected tasks saved with date: %s\n", date.Format(core.DateTpl))
}

func findAndPushGitLab(date time.Time) {
	gitLabConfig := config.ReadConfig().GitLab
	taskNames, err := getTaskNames(date)
	if err != nil {
		core.Danger("Error: %v\n", err.Error())
		return
	}
	commits := gitlab.GetCommitsByDate(gitLabConfig, date, taskNames)
	if len(commits) > 0 {
		core.Info("GitLab Commits:\n")
	}
	for _, name := range commits {
		taskReq := core.TaskReq{
			Date:         date,
			Name:         name,
			Completeness: 100,
		}
		taskRes, err := api.PushTask(taskReq)

		if err != nil {
			core.Danger("Error: %v\n", err.Error())
			return
		}

		core.Info("Task: %s (Completeness: %v%%)\n", taskRes.Data.Name, taskRes.Data.Completeness)
	}
	if len(commits) > 0 {
		core.Info("=========================================================\n")
	}
}

func findAndDelete(date time.Time) {
	from, to := core.DayRange(date)
	tasksRes, err := api.PullTasks(from, to, false)
	if err != nil {
		core.Danger("Error: %v\n", err.Error())
		return
	}

	selectedTaskNames := []string{}
	taskNames := make([]string, len(tasksRes.Data))

	for key, task := range tasksRes.Data {
		taskNames[key] = fmt.Sprintf("%s (Completeness: %v%%)", task.Name, task.Completeness)
	}

	survey.Ask(getSelectTasks(taskNames), &selectedTaskNames)
	re := regexp.MustCompile(`(.+)\s\(\w+:\s(\d+)%\)`)

	for _, n := range selectedTaskNames {
		task := tasksRes.FindByName(re.FindStringSubmatch(n)[1])
		err := api.DeleteTask(task.Id)

		if err != nil {
			core.Danger("Error: %v\n", err.Error())
			return
		}

		core.Danger("Task: %s (Completeness: %v%%) has been deleted!\n", task.Name, task.Completeness)
	}
}

func printTaskRes(tasks []core.Task) {
	if len(tasks) > 0 {
		core.Info("Tasks:\n")
		core.Info("=========================================================\n")
	}
	for key, task := range tasks {
		core.Info("[%d] %s (Completeness: %v%%)\n", key+1, task.Name, task.Completeness)
	}
	if len(tasks) > 0 {
		core.Info("\n")
	}
}
