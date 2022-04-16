package cmd

import (
	"bytes"
	"encoding/json"
	"kube-collector/pkg/collector"
	"kube-collector/pkg/k8s"
	"net/http"

	"os"
	"time"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

// logMessage is the CRI internal log type.
type logMessage struct {
	timestamp time.Time
	log       []string
}

func KubeCollector(configs *LogStream) {

	for _, collectFrom := range configs.CollectFrom {
		var podsList []*v1.PodList
		for k, v := range collectFrom.PodSelector {
			pods, err := k8s.K8s.ListPods(collectFrom.Namespace, k+"="+v)
			if err != nil {
				log.Error(err)
				return
			}
			podsList = append(podsList, pods)
		}
		for _, po := range podsList {
			for _, p := range po.Items {
				logs, err := collector.GetPodLogs(p)
				if len(logs) > 0 {
					if err != nil {
						log.Error(err)
						return
					} else {
						log.Infof("Successfully collected log from [%s] in [%s] namespace", p.GetName(), p.Namespace)
					}
					jLogs, err := json.Marshal(logs)
					if err != nil {
						return
					}

					err = httpPost(jLogs, configs.AddLabels, os.Getenv("PARSEABLE_URL")+"/api/v1/stream/"+configs.Name)
					if err != nil {
						log.Error(err)
						return
					} else {
						log.Infof("Successfully sent log from [%s] in [%s] namespace to server [%s]", p.GetName(), p.GetNamespace(), os.Getenv("PARSEABLE_URL"))
					}
				}
			}
		}

	}

}

func httpPost(logs []byte, labels map[string]string, url string) error {

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(logs))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	for key, value := range labels {
		req.Header.Add(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
