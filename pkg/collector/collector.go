package collector

import (
	"bytes"
	"context"
	"io"
	"kube-collector/pkg/k8s"
	"kube-collector/pkg/store"

	corev1 "k8s.io/api/core/v1"

	"strings"
	"time"
)

// logMessage is the log type.
type logMessage struct {
	timestamp time.Time
	log       []string
}

func GetPodLogs(pod corev1.Pod) (logMessage, error) {

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
	req := k8s.K8s.GetPodLogs(pod, podLogOpts)

	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return logMessage{}, err
	}

	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return logMessage{}, err
	}

	logs := buf.String()

	// split logs on new line
	newLogs := strings.Split(logs, "\n")

	if len(newLogs) > 1 {
		nlog := newLogs[len(newLogs)-2]
		// seperate on space
		spacedLogs := strings.Fields(nlog)

		getTimeStamp, err := time.Parse(time.RFC3339, spacedLogs[0])
		if err != nil {
			return logMessage{}, err
		}
		// put poName to TimeStamp
		store.PutPoNameTime(pod.GetName(), getTimeStamp)

		var lm logMessage
		// lm.log = make([]string, 0)
		// lm.timestamp = getTimeStamp
		// lm.log = newLogs[1:]

		// a, _ := json.Marshal(lm)

		//fmt.Println(string(a))

		return lm, nil
	} else {
		return logMessage{}, nil
	}
}
