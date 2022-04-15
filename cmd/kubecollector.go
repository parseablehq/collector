package cmd

import (
	"encoding/json"
	"kube-collector/pkg/collector"
	"kube-collector/pkg/http"
	"kube-collector/pkg/k8s"
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

				err = post2Server(jLogs, os.Getenv("PARSEABLE_URL")+"/api/v1/stream/"+configs.Name)
				if err != nil {
					log.Error(err)
					return
				} else {
					log.Infof("Successfully sent log from [%s] in [%s] namespace to server [%s]", p.GetName(), p.GetNamespace(), p.Namespace)
				}
			}
		}

	}

}

func post2Server(logs interface{}, url string) error {
	jLogs, err := json.Marshal(logs)
	if err != nil {
		return err
	}

	var http http.HTTP = http.NewClientHTTPRequests("POST", url, jLogs)

	_, err = http.HTTP()
	if err != nil {
		return err
	}
	return nil
}
