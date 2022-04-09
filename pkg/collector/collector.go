package collector

import (
	"kube-collector/pkg/k8s"
	"kube-collector/pkg/store"

	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
)

// logMessage is the log type.
type logMessage struct {
	Timestamp string `json:"time"`
	Log       string `json:"log"`
}

func GetPodLogs(pod corev1.Pod) ([]logMessage, error) {

	// poLogOptions
	var podLogOpts corev1.PodLogOptions

	// getTime on store for current pod
	if store.GetTime(pod.GetName()) != (time.Time{}) {
		var newLogTime int64
		newLogTime = int64(time.Now().Sub(store.GetTime(pod.GetName())).Seconds())
		podLogOpts = corev1.PodLogOptions{
			SinceSeconds: &newLogTime,
			Timestamps:   true,
		}
	} else {
		podLogOpts = corev1.PodLogOptions{
			Timestamps: true,
		}
	}

	// getPodLogs
	podLogs, err := k8s.K8s.GetPodLogs(pod, podLogOpts)
	if err != nil {
		return []logMessage{}, err
	}

	if len(podLogs) > 1 {
		// last line of the log
		err := putTimeStamp(pod.GetName(), podLogs)
		if err != nil {
			return []logMessage{}, err
		}

		var logMessages []logMessage

		if len(podLogs) > 1 {
			for _, lm := range podLogs {
				newLogMessage := strings.Fields(lm)
				if len(newLogMessage) > 1 {
					log := logMessage{
						Timestamp: newLogMessage[0],
						Log:       strings.Join(newLogMessage[1:], " "),
					}
					logMessages = append(logMessages, log)
				}
			}
		}

		return logMessages, nil
	} else {
		return []logMessage{}, nil
	}
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
	store.PutPoNameTime(podName, getTimeStamp)
	return err
}
