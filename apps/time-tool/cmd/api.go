package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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
	Head    Method = "HEAD"
	Post    Method = "POST"
	Put     Method = "PUT"
	Patch   Method = "PATCH" // RFC 5789
	Delete  Method = "DELETE"
	Connect Method = "CONNECT"
	Options Method = "OPTIONS"
	Trace   Method = "TRACE"
)

func GetReq(method Method, url string, jsonStr []byte) *http.Request {
	host := ReadConfig().ApiHost
	req, _ := http.NewRequest(method.toString(), host+url, bytes.NewBuffer(jsonStr))

	req.Header = http.Header{}
	req.Header.Set("Authorization", GetToken())
	req.Header.Set("Content-Type", "application/json")

	return req
}

func SignIn(email string, password string) (AccessToken, error) {
	accessToken := AccessToken{}
	reqBody, _ := json.Marshal(&Credentials{
		Email:    email,
		Password: password,
	})

	jsonStr := []byte(string(reqBody))
	req := GetReq(Post, "/api/auth/login", jsonStr)

	body, resp := GetBody(req)

	if resp.StatusCode != http.StatusCreated {
		err := Error{}
		json.Unmarshal([]byte(body), &err)

		return accessToken, errors.New(err.Message)
	}

	json.Unmarshal([]byte(body), &accessToken)

	return accessToken, nil
}

func GetBody(req *http.Request) ([]byte, *http.Response) {
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)

	return body, resp
}

func (c Method) toString() string {
	return fmt.Sprintf("%s", c)
}
