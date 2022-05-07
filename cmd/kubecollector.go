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

package cmd

import (
	"encoding/json"
	"kube-collector/pkg/client"
	"kube-collector/pkg/collector"
	"kube-collector/pkg/parseable"

	"time"

	"os"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

func RunKubeCollector(streamName string, logSpec *LogSpec) {
	if err := parseable.CreateStream(streamName); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	log.Infof("Successfully created Log Stream [%s] on server [%s]", streamName, os.Getenv("PARSEABLE_URL"))
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
				if err := parseable.PostLogs(streamName, jLogs, logSpec.Tags); err != nil {
					log.Error(err)
				}
				log.Infof("Successfully sent log from [%s] in [%s] namespace to server [%s]", p.GetName(), p.GetNamespace(), os.Getenv("PARSEABLE_URL"))
			}
		}
	}
}
