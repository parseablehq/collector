package http

import (
	"bytes"
	"net/http"
)

const METADATA_LABEL = "X-P-META-"

type HttpParseable interface {
	DoHttpRequest() (*http.Response, error)
}

// httpRequest holds all the fields needed for a HTTP request
// to parseable server.
type HttpRequest struct {
	method string
	url    string
	labels map[string]string
	body   []byte
}

func NewHttpRequest(method, url string, labels map[string]string, body []byte) *HttpRequest {
	return &HttpRequest{method: method, url: url, labels: labels, body: body}
}

func (h *HttpRequest) DoHttpRequest() (*http.Response, error) {
	req, err := http.NewRequest(h.method, h.url, bytes.NewBuffer(h.body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	if h.labels != nil {
		for key, value := range h.labels {
			req.Header.Add(METADATA_LABEL+key, value)
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
