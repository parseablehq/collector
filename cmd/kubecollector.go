package cmd

import (
	"encoding/json"
	"kube-collector/pkg/client"
	"kube-collector/pkg/collector"
	"kube-collector/pkg/http"
	"kube-collector/pkg/utils"

	"time"

	"os"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

func RunKubeCollector(streamName string, logSpec *LogSpec) {
	// Create stream

	var http http.HttpParseable = http.NewHttpRequest("PUT", utils.GetParseableStreamURL(streamName), nil, nil)

	_, err := http.DoHttpRequest()
	if err != nil {
		// TODO: Make sure to ignore the error if the stream already exists
		log.Error("Failed to create Log Stream due to error: ", err.Error())
		return
	} else {
		log.Infof("Successfully created Log Stream [%s] on server [%s]", streamName, os.Getenv("PARSEABLE_URL"))
	}

	interval, err := time.ParseDuration(logSpec.CollectInterval)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	ticker := time.NewTicker(interval)
	for range ticker.C {
		kubeCollector(streamName, logSpec)
	}
}

func kubeCollector(streamName string, logSpec *LogSpec) {

	collectFrom := logSpec.CollectFrom
	var podsList []*v1.PodList
	for k, v := range collectFrom.PodSelector {
		pods, err := client.KubeClient.ListPods(collectFrom.Namespace, k+"="+v)
		if err != nil {
			log.Error(err)
			return
		}
		podsList = append(podsList, pods)
	}
	for _, po := range podsList {
		for _, p := range po.Items {
			logs, err := collector.GetPodLogs(p, streamName)
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
				var http http.HttpParseable = http.NewHttpRequest("POST", utils.GetParseableStreamURL(streamName), logSpec.AddTags, jLogs)

				_, err = http.DoHttpRequest()
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
