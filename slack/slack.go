package slack

import (
	"beagle/logger"
	"bytes"
	"encoding/json"
	"net/http"
)

var log = logger.New("slack")

// Notify sends Slack notification to defined channel
var Notify = func(url string, msg map[string]any) {

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Error(err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()

}
