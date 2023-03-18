package gitlab

import (
	"encoding/json"
	"finlab/apps/time-tool/core"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

type GitLabConfig struct {
	Host   string `json:"host"`
	Token  string `json:"token"`
	UserId int16  `json:"userId"`
}

type PushData struct {
	CommitTitle string `json:"commit_title"`
	CommitTo    string `json:"commit_to"`
}

type Event struct {
	PushData  PushData `json:"push_data"`
	Author    string   `json:"author_username"`
	Date      string   `json:"created_at"`
	ProjectId int      `json:"project_id"`
}

type Commit struct {
	Message string `json:"message"`
}

type ProjectCommit struct {
	ProjectId int
	CommitTo  string
}

type UserID struct {
	UserId int16
}

const dateFormat string = "2006-01-02"
const urlUserEvents string = "/users/{{.UserId}}/events"
const urlCommit string = "/projects/{{.ProjectId}}/repository/commits/{{.CommitTo}}"

func GetCommitsByDate(gitLabConfig GitLabConfig, date time.Time, names []string) []string {
	url := core.GetUrl(urlUserEvents, UserID{UserId: gitLabConfig.UserId})
	req := getReq(gitLabConfig, url)
	query := req.URL.Query()
	query.Add("after", date.AddDate(0, 0, -1).Format(dateFormat))
	query.Add("before", date.AddDate(0, 0, 1).Format(dateFormat))
	req.URL.RawQuery = query.Encode()

	body := getBody(req)
	data := getData(body)
	result := make([]string, 0)

	for _, event := range uniqueCommitTitle(data) {
		projectCommit := ProjectCommit{
			ProjectId: event.ProjectId,
			CommitTo:  event.PushData.CommitTo,
		}
		commitMsg := getCommitMessage(gitLabConfig, projectCommit)
		if !core.Contains(names, commitMsg) {
			result = append(result, commitMsg)
		}
	}

	return result
}

func getCommitMessage(gitLabConfig GitLabConfig, projectCommit ProjectCommit) string {
	url := core.GetUrl(urlCommit, projectCommit)
	req := getReq(gitLabConfig, url)
	body := getBody(req)
	result := Commit{}
	json.Unmarshal([]byte(body), &result)
	messageSlice := strings.Split(strings.ReplaceAll(result.Message, "\r\n", "\n"), "\n")
	return strings.TrimSuffix(messageSlice[0], "\n")
}

func getReq(gitLabConfig GitLabConfig, url string) *http.Request {
	req, _ := http.NewRequest("GET", gitLabConfig.Host+url, nil)

	req.Header = http.Header{}
	req.Header.Set("Authorization", gitLabConfig.Token)
	req.Header.Set("Content-Type", "application/json")

	return req
}

func getBody(req *http.Request) []byte {
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

func getData(body []byte) []Event {
	result := []Event{}
	json.Unmarshal([]byte(body), &result)

	return result
}

func uniqueCommitTitle(events []Event) []Event {
	return core.Filter(events, func(result []Event, event Event) bool {
		index := slices.IndexFunc(result, func(e Event) bool {
			return event.PushData.CommitTitle == e.PushData.CommitTitle
		})

		return index == -1 && len(event.PushData.CommitTitle) > 0
	})
}
