package cmd

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	DEFAULT_LOG_COLLECT_INTERVAL = "1m"
)

type CollectorConfig struct {
	LogStreams []LogStream `yaml:"logStreams"`
}

type LogStream struct {
	Name    string  `yaml:"name"`
	LogSpec LogSpec `yaml:"logSpec"`
}

type LogSpec struct {
	AddTags         map[string]string `yaml:"tagsToAdd"`
	CollectInterval string            `yaml:"collectInterval"`
	CollectFrom     CollectFrom       `yaml:"collectFrom"`
}

type CollectFrom struct {
	Namespace   string            `yaml:"namespace"`
	PodSelector map[string]string `yaml:"podSelector"`
}

func ReadConfig(path *string) (*CollectorConfig, error) {
	configfile, err := ioutil.ReadFile(*path)
	if err != nil {
		return nil, err
	}

	var logConfig CollectorConfig

	err = yaml.Unmarshal([]byte(configfile), &logConfig)
	if err != nil {
		return nil, err
	}
	logConfig.ensureDefaults()
	return &logConfig, nil
}

func (logConfig *CollectorConfig) ensureDefaults() {
	for _, logStream := range logConfig.LogStreams {
		if logStream.LogSpec.CollectInterval == "" {
			logStream.LogSpec.CollectInterval = DEFAULT_LOG_COLLECT_INTERVAL
		}
	}
}
