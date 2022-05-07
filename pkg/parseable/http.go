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
	tags   map[string]string
	body   []byte
}

func newRequest(method, url string, tags map[string]string, body []byte) *httpRequest {
	return &httpRequest{method: method, url: url, tags: tags, body: body}
}

func (h *httpRequest) Do() (*http.Response, error) {
	req, err := http.NewRequest(h.method, h.url, bytes.NewBuffer(h.body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	if h.tags != nil {
		for key, value := range h.tags {
			req.Header.Add(METADATA_LABEL+key, value)
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
