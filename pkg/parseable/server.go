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
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type MaxTimeQuery []struct {
	MAXSystemsTime string `json:"MAX(systems.time)"`
}

func CreateStream(streamName string) error {
	req := newRequest("PUT", streamURL(streamName), nil, nil)
	if resp, err := req.Do(); err != nil {
		return err
	} else if resp.StatusCode == 400 {
		// Server retruns 400 if stream already exists
		// we ignore that error and return nil
		return nil
	} else if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d while creating stream: %s", resp.StatusCode, streamName)
	}
	return nil
}

func PostLogs(streamName string, logs []byte, tags map[string]string) error {
	req := newRequest("POST", streamURL(streamName), tags, logs)
	if resp, err := req.Do(); err != nil {
		return err
	} else if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d while posting log data to stream: %s", resp.StatusCode, streamName)
	}
	return nil
}

func LastLogTime(streamName, podName, containerName string) (MaxTimeQuery, error) {
	query := map[string]string{
		"query": fmt.Sprintf("select max(time) from %s where meta_PodName = '%s' and meta_ContainerName = '%s'", streamName, podName, containerName),
	}

	queryJson, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	req := newRequest("GET", queryURL(), nil, queryJson)
	resp, err := req.Do()
	if err != nil {
		return nil, err
	} else if resp.StatusCode == 500 {
		// This is the case where the log stream is empty
		respData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if string(respData) == "Error during planning: No data file found" {
			return nil, nil
		}
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d while querying log data timestamp in stream: %s", resp.StatusCode, streamName)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var mtq MaxTimeQuery
	err = json.Unmarshal(respData, &mtq)
	if err != nil {
		return nil, err
	}

	return mtq, nil
}
