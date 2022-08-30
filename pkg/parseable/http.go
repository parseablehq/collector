// Copyright (C) 2022 Parseable, Inc.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package parseable

import (
	"bytes"
	"net/http"
	"runtime"
	"strings"
)

const METADATA_LABEL = "X-P-META-"

type HttpParseable interface {
	Do() (*http.Response, error)
}

// httpRequest holds all the fields needed for a HTTP request
// to parseable server.
type httpRequest struct {
	method string
	url    string
	labels map[string]string
	body   []byte
}

func newRequest(method, url string, labels map[string]string, body []byte) *httpRequest {
	return &httpRequest{method: method, url: url, labels: labels, body: body}
}

func (h *httpRequest) Do(user, pwd string) (*http.Response, error) {
	req, err := http.NewRequest(h.method, h.url, bytes.NewBuffer(h.body))
	if err != nil {
		return nil, err
	}
	if user != "" && pwd != "" {
		req.SetBasicAuth(user, pwd)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", getUserAgent())

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

func getUserAgent() string {
	userAgentParts := []string{}
	uaAppend := func(p, q string) {
		userAgentParts = append(userAgentParts, p, q)
	}

	uaAppend("Parseable collector (", runtime.GOOS)
	uaAppend("; ", runtime.GOARCH)
	uaAppend("; ", runtime.Version())

	return strings.Join(userAgentParts, ")")
}
