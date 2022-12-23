package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
}

type Error struct {
	StatusCode int8   `json:"statusCode"`
	Message    string `json:"message"`
	Error      string `json:"error"`
}

type Timestamp struct {
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
}

type TimestampRes struct {
	Data Timestamp `json:"data"`
}

type Method string

const (
	Get     Method = "GET"
	Head           = "HEAD"
	Post           = "POST"
	Put            = "PUT"
	Patch          = "PATCH" // RFC 5789
	Delete         = "DELETE"
	Connect        = "CONNECT"
	Options        = "OPTIONS"
	Trace          = "TRACE"
)

const (
	Host            string = "http://localhost:3333"
	AccessTokenFile string = "access-token.json"
)

func GetReq(method Method, url string, jsonStr []byte) *http.Request {
	req, _ := http.NewRequest(method.toString(), Host+url, bytes.NewBuffer(jsonStr))

	req.Header = http.Header{}
	req.Header.Set("Authorization", getToken())
	req.Header.Set("Content-Type", "application/json")

	return req
}

func GetBody(req *http.Request) []byte {
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

func getToken() string {
	var accessToken AccessToken = AccessToken{}
	err := Load(AccessTokenFile, &accessToken)
	if err != nil {
		log.Fatalln(err)
	}

	return "Bearer " + accessToken.AccessToken
}

func (c Method) toString() string {
	return fmt.Sprintf("%s", c)
}
