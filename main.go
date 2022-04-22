package main

import (
	"flag"
	"kube-collector/cmd"

	log "github.com/sirupsen/logrus"

	"os"
	"sync"
	"time"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "config file for kube-collector")
	flag.Parse()
	if len(configPath) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

}

func main() {

	config, err := cmd.ReadConfig(&configPath)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	var wg sync.WaitGroup

	for _, stream := range config.LogStreams {
		wg.Add(1)
		go func(name string, logSpec cmd.LogSpec) {
			defer wg.Done()
			runKubeCollector(name, &logSpec)
		}(stream.Name, stream.LogSpec)
	}
	wg.Wait()
}

func runKubeCollector(name string, logSpec *cmd.LogSpec) {
	interval, err := time.ParseDuration(logSpec.Interval)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for range ticker.C {
		cmd.KubeCollector(name, logSpec)
	}
}
