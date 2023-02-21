package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"finlab/apps/time-tool/config"
	"finlab/apps/time-tool/core"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func GetReq(method core.Method, url string, jsonStr []byte) *http.Request {
	host := config.ReadConfig().ApiHost
	req, _ := http.NewRequest(method.ToString(), host+url, bytes.NewBuffer(jsonStr))

	req.Header = http.Header{}
	req.Header.Set("Authorization", config.GetToken())
	req.Header.Set("Content-Type", "application/json")

	return req
}

func GetBody(req *http.Request) ([]byte, *http.Response, error) {
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, errors.New("Failed to connect to finlab server")
	}
	body, _ := ioutil.ReadAll(resp.Body)

	return body, resp, nil
}

func SignIn(credentials core.Credentials) (core.AccessToken, error) {
	accessToken := core.AccessToken{}
	reqBody, _ := json.Marshal(credentials)

	jsonStr := []byte(string(reqBody))
	req := GetReq(core.Post, "/api/auth/login", jsonStr)

	body, resp, err := GetBody(req)

	if err != nil {
		return accessToken, err
	}

	if resp.StatusCode != http.StatusCreated {
		err := core.Error{}
		json.Unmarshal([]byte(body), &err)

		return accessToken, errors.New(err.Message)
	}

	json.Unmarshal([]byte(body), &accessToken)

	return accessToken, nil
}

func PushTimestamp(timestamp core.Timestamp) (core.TimestampRes, error) {
	timestampRes := core.TimestampRes{}
	reqBody, _ := json.Marshal(&core.TimestampReq{
		Timestamp: timestamp.Timestamp.Format("2006-01-02T15:04:05Z"),
		Type:      timestamp.Type,
	})

	jsonStr := []byte(string(reqBody))
	req := GetReq(core.Post, "/api/work-time/timestamp", jsonStr)
	body, resp, err := GetBody(req)

	if err != nil {
		return timestampRes, err
	}

	if resp.StatusCode != http.StatusCreated {
		err := core.Error{}
		json.Unmarshal([]byte(body), &err)

		return timestampRes, errors.New(err.Message)
	}

	json.Unmarshal([]byte(body), &timestampRes)

	return timestampRes, nil
}

func PullTimestamps(date time.Time, raw bool) (core.TimestampsRes, error) {
	timestampsRes := core.TimestampsRes{}

	jsonStr := []byte("")
	req := GetReq(core.Get, fmt.Sprintf("/api/work-time/timestamp?date=%s&raw=%v", date.Format("2006-01-02"), raw), jsonStr)
	body, resp, err := GetBody(req)

	if err != nil {
		return timestampsRes, err
	}

	if resp.StatusCode != http.StatusOK {
		err := core.Error{}
		json.Unmarshal([]byte(body), &err)

		return timestampsRes, errors.New(err.Message)
	}

	json.Unmarshal([]byte(body), &timestampsRes)

	return timestampsRes, nil
}

func PushTask(task core.TaskReq) (core.TaskRes, error) {
	taskRes := core.TaskRes{}
	reqBody, _ := json.Marshal(task)

	jsonStr := []byte(string(reqBody))
	req := GetReq(core.Post, "/api/work-time/task", jsonStr)
	body, resp, err := GetBody(req)

	if err != nil {
		return taskRes, err
	}

	if resp.StatusCode != http.StatusCreated {
		err := core.Error{}
		json.Unmarshal([]byte(body), &err)

		return taskRes, errors.New(err.Message)
	}

	json.Unmarshal([]byte(body), &taskRes)

	return taskRes, nil
}

func PullTasks(from time.Time, to time.Time, incomplete bool) (core.TasksRes, error) {
	tasksRes := core.TasksRes{}

	jsonStr := []byte("")
	req := GetReq(core.Get, fmt.Sprintf("/api/work-time/task?from=%s&to=%s&incomplete=%v", from.Format(core.DateISOTpl), to.Format(core.DateISOTpl), incomplete), jsonStr)
	body, resp, err := GetBody(req)

	if err != nil {
		return tasksRes, err
	}

	if resp.StatusCode != http.StatusOK {
		err := core.Error{}
		json.Unmarshal([]byte(body), &err)

		return tasksRes, errors.New(err.Message)
	}

	json.Unmarshal([]byte(body), &tasksRes)

	return tasksRes, nil
}

func PullSummary(from time.Time, to time.Time) (core.SummaryRes, error) {
	summaryRes := core.SummaryRes{}

	jsonStr := []byte("")
	req := GetReq(core.Get, fmt.Sprintf("/api/work-time/summary?from=%s&to=%s", from.Format(core.DateISOTpl), to.Format(core.DateISOTpl)), jsonStr)
	body, resp, err := GetBody(req)

	if err != nil {
		return summaryRes, err
	}

	if resp.StatusCode != http.StatusOK {
		err := core.Error{}
		json.Unmarshal([]byte(body), &err)

		return summaryRes, errors.New(err.Message)
	}

	json.Unmarshal([]byte(body), &summaryRes)

	return summaryRes, nil
}
