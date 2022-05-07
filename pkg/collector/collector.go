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

package collector

import (
	"fmt"
	"kube-collector/pkg/client"
	"kube-collector/pkg/parseable"
	"kube-collector/pkg/store"

	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
)

// logMessage represents a single log message entry
type logMessage struct {
	Timestamp string      `json:"time"`
	Log       string      `json:"log"`
	LogMeta   logMetadata `json:"meta"`
}

type logMetadata struct {
	Host           string
	Source         string
	ContainerName  string
	ContainerImage string
	PodName        string
	Namespace      string
}

func GetPodLogs(pod corev1.Pod, streamName string) ([]logMessage, error) {

	for _, container := range pod.Spec.Containers {
		podContainerName := pod.GetName() + "/" + container.Name
		podLogOpts := corev1.PodLogOptions{
			Timestamps: true,
			Container:  container.Name,
		}

		if store.IsEmpty() == true {
			mtq, err := parseable.LastLogTime(streamName, pod.Name, container.Name)
			if err != nil {
				return nil, err
			}

			if mtq != nil {
				time, err := time.Parse(time.RFC3339, mtq[0].MAXSystemsTime)
				if err != nil {
					return nil, err
				}
				store.SetLastTimestamp(podContainerName, time)
			}
		}

		// use a combination of pod and container name to store the last
		// time stamp. This ensure we can uniquely fetch a container's log
		lastLogTime := store.LastTimestamp(podContainerName)
		if lastLogTime != (time.Time{}) {
			secsSinceLastLog := int64(time.Now().Sub(lastLogTime).Seconds())
			podLogOpts.SinceSeconds = &secsSinceLastLog
		}

		podLogs, err := client.KubeClient.GetPodLogs(pod, podLogOpts)
		if err != nil {
			return nil, err
		}

		if len(podLogs) > 1 {
			// last line of the log
			if err := putTimeStamp(podContainerName, podLogs); err != nil {
				return nil, err
			}
			var logMessages []logMessage
			for _, lm := range podLogs {
				newLogMessage := strings.Fields(lm)
				if len(newLogMessage) > 1 {
					log := logMessage{
						Timestamp: newLogMessage[0],
						Log:       strings.Join(newLogMessage[1:], " "),
						LogMeta: logMetadata{
							Namespace:      pod.GetNamespace(),
							Host:           pod.Status.HostIP,
							Source:         pod.Status.PodIP,
							ContainerName:  container.Name,
							ContainerImage: container.Image,
							PodName:        pod.GetName(),
						},
					}
					logMessages = append(logMessages, log)
				}
			}
			return logMessages, nil
		}
	}
	return []logMessage{}, nil
}

func putTimeStamp(podName string, podLogs []string) error {
	lastLinelog := podLogs[len(podLogs)-2]
	// seperate on space for last line
	spacedLogs := strings.Fields(lastLinelog)
	// get the timestamp appended to log by k8s
	getTimeStamp, err := time.Parse(time.RFC3339, spacedLogs[0])
	if err != nil {
		return err
	}
	// put poName to TimeStamp
	store.SetLastTimestamp(podName, getTimeStamp)
	fmt.Println(store.PoNameTime)
	return err
}
