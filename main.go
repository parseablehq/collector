package main

import (
	"fmt"
	"kube-collector/cmd"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := cmd.ReadConfig("config.yaml")
	if err != nil {
		log.Error(err)
	}
	for _, logStream := range config.LogStreams {
		runKubeCollector(time.Duration(logStream.CollectInterval), config)
	}
}

func runKubeCollector(interval time.Duration, config *cmd.LogConfig) {
	ticker := time.NewTicker(interval * time.Second)

	for t := range ticker.C {
		fmt.Println("Invoked at ", t)
		cmd.KubeCollector(config)
	}
}
