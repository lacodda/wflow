package db

import (
	"encoding/json"
	"finlab/apps/time-tool/core"
)

const (
	schemaTaskSql = `CREATE TABLE IF NOT EXISTS task (
		id INTEGER NOT NULL PRIMARY KEY,
        date DATETIME NOT NULL,
        name TEXT NOT NULL,
        comment TEXT,
        completeness INT);`
	insertTaskSql = `INSERT INTO task (date, name, comment, completeness) VALUES (?, ?, ?, ?);`
	selectTaskSql = `SELECT * FROM task`
	deleteTaskSql = `DELETE FROM task`
)

func GetTasksRepo() (repo *SQLiteRepository, err error) {
	repo, err = Repo(schemaTaskSql)
	if err != nil {
		return nil, err
	}
	return
}

func SetTask(task core.TaskReq) error {
	repo, err := GetTasksRepo()
	if err != nil {
		return err
	}

	_, err = repo.Create(insertTaskSql, task.Date, task.Name, task.Comment, task.Completeness)
	return err
}

func GetTasks(ids ...[]int) ([]core.Task, error) {
	repo, err := GetTasksRepo()
	if err != nil {
		return nil, err
	}

	var tasksJson []string

	if len(ids) > 0 && ids[0] != nil {
		tasksJson, err = repo.Select(selectTaskSql, ids[0])
	} else {
		tasksJson, err = repo.Select(selectTaskSql)
	}

	if err != nil {
		return nil, err
	}

	var tasks []core.Task

	for _, taskJson := range tasksJson {
		var task core.Task
		json.Unmarshal([]byte(taskJson), &task)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func DeleteTasksByIds(ids []int) ([]core.Task, error) {
	repo, err := GetTasksRepo()
	if err != nil {
		return nil, err
	}

	tasks, err := GetTasks(ids)
	if err != nil {
		return nil, err
	}

	err = repo.Delete(deleteTaskSql, ids)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func DeleteTasks() error {
	repo, err := GetTasksRepo()
	if err != nil {
		return err
	}

	err = repo.Delete(deleteTaskSql)
	if err != nil {
		return err
	}

	return nil
}
