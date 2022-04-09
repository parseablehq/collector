package http

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
)

// A HTTP interface
type HTTP interface {
	HTTP() (bodyBytes []byte, err error)
}

// Generic struct for client HTTP requests
type ClientHTTPRequests struct {
	url    string
	method string
	body   []byte
	// basic auth
	username string
	password string
}

//   Constructor to ClientHTTPRequests Struct
func NewClientHTTPRequests(methods string, url string, body []byte) *ClientHTTPRequests {
	return &ClientHTTPRequests{
		method:   methods,
		body:     body,
		username: os.Getenv("USERNAME"),
		password: os.Getenv("PASSWORD"),
	}
}

func (c ClientHTTPRequests) HTTP() (bodyBytes []byte, err error) {
	req, err := http.NewRequest(c.method, c.url, bytes.NewBuffer(c.body))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respData, nil

}
