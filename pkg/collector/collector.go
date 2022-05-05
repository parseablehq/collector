package collector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kube-collector/pkg/client"
	"kube-collector/pkg/http"
	"kube-collector/pkg/store"
	"kube-collector/pkg/utils"

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

		if store.IsStoreEmpty(podContainerName) {

			query := fmt.Sprintf("select max(time) from %s where meta_PodName = '%s' and meta_ContainerName = '%s'", streamName, pod.GetName(), container.Name)
			createQuery := map[string]string{
				"query": query,
			}

			jQuery, err := json.Marshal(createQuery)
			if err != nil {
				return nil, err
			}
			var http http.HttpParseable = http.NewHttpRequest("POST", utils.GetParseableQueryURL(), nil, jQuery)
			resp, err := http.DoHttpRequest()
			if err != nil {
				return nil, err
			}
			respData, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			time, err := time.Parse(time.RFC3339, string(respData))
			if err != nil {
				return nil, err
			}
			store.SetLastTimestamp(podContainerName, time)
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
	return err
}
