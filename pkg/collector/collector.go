package collector

import (
	"kube-collector/pkg/client"
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
}

func GetPodLogs(pod corev1.Pod) ([]logMessage, error) {

	for _, container := range pod.Spec.Containers {
		podLogOpts := corev1.PodLogOptions{
			Timestamps: true,
			Container:  container.Name,
		}

		lastLogTime := store.LastTimestamp(pod.GetName())
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
			if err := putTimeStamp(pod.GetName(), podLogs); err != nil {
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
