package cmd

import (
	"encoding/json"
	"path/filepath"

	"github.com/shibukawa/configdir"
)

type Config struct {
	Login   string `json:"login"`
	ApiHost string `json:"apiHost"`
}

var DefaultConfig = Config{
	Login:   "",
	ApiHost: "http://localhost:3333",
}

const (
	ConfigFile      string = "config.json"
	AccessTokenFile string = "access-token.json"
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

func SaveToken(token AccessToken) {
	data, _ := json.Marshal(&token)
	folders := configDirs.QueryFolders(configdir.Global)
	folders[0].WriteFile(AccessTokenFile, data)
}

func GetToken() string {
	var accessToken AccessToken

	folder := configDirs.QueryFolderContainsFile(AccessTokenFile)
	if folder != nil {
		data, _ := folder.ReadFile(AccessTokenFile)
		json.Unmarshal(data, &accessToken)
	} else {
		accessToken = AccessToken{AccessToken: ""}
	}

	return "Bearer " + accessToken.AccessToken
}

func RemoveToken() {
	folder := configDirs.QueryFolderContainsFile(AccessTokenFile)
	if folder != nil {
		Remove(filepath.Join(folder.Path, AccessTokenFile))
	}
}
