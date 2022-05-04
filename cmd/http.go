package cmd

import (
	"bytes"
	"net/http"
	"os"
)

const METADATA_LABEL = "X-P-META-"

func parseableStreamURL(streamName string) string {
	return os.Getenv("PARSEABLE_URL") + "/api/v1/stream/" + streamName
}

type httpParseable interface {
	doHttpRequest() error
}

// httpRequest holds all the fields needed for a HTTP request
// to parseable server.
type httpRequest struct {
	method string
	url    string
	labels map[string]string
	body   []byte
}

func newHttpRequest(method, url string, labels map[string]string, body []byte) *httpRequest {
	return &httpRequest{method: method, url: url, labels: labels, body: body}
}

func (h *httpRequest) doHttpRequest() error {
	req, err := http.NewRequest(h.method, h.url, bytes.NewBuffer(h.body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	if h.labels != nil {
		for key, value := range h.labels {
			req.Header.Add(METADATA_LABEL+key, value)
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
