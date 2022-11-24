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
	"collector/pkg/client"
	"collector/pkg/parseable"
	"collector/pkg/store"

	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
)

// logMessage represents a single log message entry
type logMessage struct {
	Timestamp string `json:"time"`
	Log       string `json:"log"`
}

func GetPodLogs(pod corev1.Pod, url, user, pwd, streamName string) ([]logMessage, map[string]string, error) {

	for _, container := range pod.Spec.Containers {
		podContainerName := pod.GetName() + "/" + container.Name
		podLogOpts := corev1.PodLogOptions{
			Timestamps: true,
			Container:  container.Name,
		}
		// use a combination of pod and container name to store the last
		// time stamp. This ensure we can uniquely fetch a container's log
		lastLogTime, ok := store.LastTimestamp(podContainerName)
		if lastLogTime == (time.Time{}) || !ok {
			mtq, _ := parseable.LastLogTime(url, user, pwd, streamName, pod.Name, container.Name)
			// if err != nil {
			// 	//return nil, nil, err
			// }
			if len(mtq) > 3 {
				time, err := time.Parse(time.RFC3339, mtq[0].MAXSystemsTime)
				if err != nil {
					return nil, nil, err
				}
				store.SetLastTimestamp(podContainerName, time)
				lastLogTime = time
			}
		}

		if lastLogTime != (time.Time{}) {
			secsSinceLastLog := int64(time.Since(lastLogTime).Seconds())
			podLogOpts.SinceSeconds = &secsSinceLastLog
		}

		podLogs, err := client.KubeClient.GetPodLogs(pod, podLogOpts)
		if err != nil {
			return nil, nil, err
		}
		if len(podLogs) > 1 {
			// last line of the log
			if err := putTimeStamp(podContainerName, podLogs); err != nil {
				return nil, nil, err
			}
			var logMessages []logMessage
			var LogMeta map[string]string
			for _, lm := range podLogs {
				newLogMessage := strings.Fields(lm)
				if len(newLogMessage) > 1 {
					log := logMessage{
						Timestamp: newLogMessage[0],
						Log:       strings.Join(newLogMessage[1:], " "),
					}
					logMessages = append(logMessages, log)
					LogMeta = map[string]string{
						"Namespace":      pod.GetNamespace(),
						"Host":           pod.Status.HostIP,
						"Source":         pod.Status.PodIP,
						"ContainerName":  container.Name,
						"ContainerImage": container.Image,
						"PodName":        pod.GetName(),
						"PodLabels":      map2string(pod.GetLabels()),
					}
				}
			}
			return logMessages, LogMeta, nil
		}
	}
	return []logMessage{}, nil, nil
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
	return err
}

func map2string(m map[string]string) string {
	var labels []string
	for key, value := range m {
		labels = append(labels, key+"="+value)
	}
	return strings.Join(labels, ",")
}
