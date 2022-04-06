package main

import (
	"fmt"
	"io/ioutil"
	"kube-collector/pkg/collector"
	"kube-collector/pkg/k8s"
	"kube-collector/pkg/store"
	"time"

	"gopkg.in/yaml.v3"
)

// logMessage is the CRI internal log type.
type logMessage struct {
	timestamp time.Time
	log       []string
}

type LogConfig struct {
	LogStreams []struct {
		Name        string `yaml:"name"`
		CollectFrom []struct {
			Namespace   string            `yaml:"namespace"`
			PodSelector map[string]string `yaml:"podSelector"`
		} `yaml:"collectFrom"`
	} `yaml:"logStreams"`
}

func main() {
	configfile, _ := ioutil.ReadFile("config.yaml")

	var logConfig LogConfig
	yaml.Unmarshal([]byte(configfile), &logConfig)

	// 1. read yaml file
	// 2. namespace and label selectors

	configs := logConfig.LogStreams

	for _, v := range configs {
		for _, vv := range v.CollectFrom {
			pods, _ := k8s.K8s.ListPods(vv.Namespace, "namespace=operator")

			ticker := time.NewTicker(5 * time.Second)
			for t := range ticker.C {
				for _, p := range pods.Items {
					fmt.Println(store.GetTime(p.Name))

					fmt.Println(p.GetName())
					fmt.Println("Invoked at ", t)

					collector.GetPodLogs(p)
					///fmt.Println(a)
				}
			}
		}

	}

}
