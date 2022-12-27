package cmd

import (
	"encoding/json"
	"time"

	"github.com/spf13/cobra"
)

var TimestampCmd = &cobra.Command{
	Use:   "timestamp",
	Short: "Write timestamp and event type to database",
	Run: func(cmd *cobra.Command, args []string) {
		reqBody, _ := json.Marshal(&Timestamp{
			Timestamp: time.Now().Format("2006-01-02T15:04:05Z"),
			Type:      "StartBreak",
		})

		jsonStr := []byte(string(reqBody))
		req := GetReq(Post, "/api/work-time/timestamp", jsonStr)
		body, _ := GetBody(req)
		result := TimestampRes{}
		json.Unmarshal([]byte(body), &result)
		Info("Timestamp: %s\n", result.Data.Timestamp)
	},
}
