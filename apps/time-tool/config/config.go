package config

import (
	"encoding/json"
	"finlab/apps/time-tool/core"
	"finlab/apps/time-tool/gitlab"
	"finlab/apps/time-tool/mail"
	"path/filepath"

	"github.com/shibukawa/configdir"
)

type Config struct {
	Login    string              `json:"login"`
	ApiHost  string              `json:"apiHost"`
	Mail     mail.Mail           `json:"mail"`
	TestMail mail.Mail           `json:"testMail"`
	GitLab   gitlab.GitLabConfig `json:"gitLab"`
}

var DefaultConfig = Config{
	Login:   "",
	ApiHost: "http://localhost:3333",
}

const (
	ConfigFile      string = "config.json"
	AccessTokenFile string = "access-token.json"
	DbFile          string = "data.db"
)

var configDirs = configdir.New("lacodda", "time-tool")

func ReadConfig() Config {
	var config Config

	folder := configDirs.QueryFolderContainsFile(ConfigFile)
	if folder != nil {
		data, _ := folder.ReadFile(ConfigFile)
		json.Unmarshal(data, &config)
	} else {
		config = DefaultConfig
	}

	return config
}

func WriteConfig(config Config) {
	data, _ := json.Marshal(&config)
	folders := configDirs.QueryFolders(configdir.Global)
	folders[0].WriteFile(ConfigFile, data)
}

func InitConfig() {
	WriteConfig(DefaultConfig)
}

func SaveToken(token core.AccessToken) {
	data, _ := json.Marshal(&token)
	folders := configDirs.QueryFolders(configdir.Global)
	folders[0].WriteFile(AccessTokenFile, data)
}

func GetToken() string {
	var accessToken core.AccessToken

	folder := configDirs.QueryFolderContainsFile(AccessTokenFile)
	if folder != nil {
		data, _ := folder.ReadFile(AccessTokenFile)
		json.Unmarshal(data, &accessToken)
	} else {
		accessToken = core.AccessToken{AccessToken: ""}
	}

	return "Bearer " + accessToken.AccessToken
}

func RemoveToken() {
	folder := configDirs.QueryFolderContainsFile(AccessTokenFile)
	if folder != nil {
		core.Remove(filepath.Join(folder.Path, AccessTokenFile))
	}
}

func DbPath() string {
	folders := configDirs.QueryFolders(configdir.Global)
	if !folders[0].Exists(DbFile) {
		folders[0].MkdirAll()
	}

	return filepath.Join(folders[0].Path, DbFile)
}
