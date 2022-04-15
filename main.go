package main

import (
	"flag"
	"kube-collector/cmd"

	log "github.com/sirupsen/logrus"

	"os"
	"sync"
	"time"
)

var config string

func init() {
	flag.StringVar(&config, "config", "", "config file for kube-collector")
	flag.Parse()
	if len(config) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

}

func main() {

	config, err := cmd.ReadConfig(&config)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	var wg sync.WaitGroup

	for _, logStream := range config.LogStreams {
		wg.Add(1)
		go func(logS cmd.LogStream) {
			defer wg.Done()
			runKubeCollector(&logS)
		}(logStream)

	}
	wg.Wait()
}

func runKubeCollector(logStream *cmd.LogStream) {
	ticker := time.NewTicker(time.Duration(logStream.CollectInterval) * time.Second)
	for range ticker.C {
		cmd.KubeCollector(logStream)
	}
}
