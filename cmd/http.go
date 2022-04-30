package cmd

import (
	"bytes"
	"net/http"
	"os"
)

func parseableStreamURL(streamName string) string {
	return os.Getenv("PARSEABLE_URL") + "/api/v1/stream/" + streamName
}

func httpPost(logs []byte, labels map[string]string, url string) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(logs))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for key, value := range labels {
		req.Header.Add(key, value)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func httpPut(url string) error {
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
