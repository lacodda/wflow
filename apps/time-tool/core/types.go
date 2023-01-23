package core

import (
	"fmt"
	"time"
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
	Id        uint          `json:"id"`
	Timestamp time.Time     `json:"timestamp"`
	Type      TimestampType `json:"type"`
}

type TimestampReq struct {
	Timestamp string        `json:"timestamp"`
	Type      TimestampType `json:"type"`
}

type TimestampRes struct {
	Data TimestampReq `json:"data"`
}

type TimestampsRes struct {
	Data      []TimestampReq `json:"data"`
	TotalTime int            `json:"totalTime"`
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

type TimestampType string

const (
	Start      TimestampType = "Start"
	End        TimestampType = "End"
	StartBreak TimestampType = "StartBreak"
	EndBreak   TimestampType = "EndBreak"
)

func (c Method) ToString() string {
	return fmt.Sprintf("%s", c)
}
