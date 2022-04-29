package cmd

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type CollectorConfig struct {
	LogStreams []LogStream `yaml:"logStreams"`
}

type LogStream struct {
	Name    string  `yaml:"name"`
	LogSpec LogSpec `yaml:"logSpec"`
}

type LogSpec struct {
	AddTags     map[string]string `yaml:"tagsToAdd"`
	Interval    string            `yaml:"collectionInterval"`
	CollectFrom CollectFrom       `yaml:"collectFrom"`
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
	// TODO -- add defaults for collectInterval and tagsToAdd
	err = yaml.Unmarshal([]byte(configfile), &logConfig)
	if err != nil {
		return nil, err
	}
	return &logConfig, nil
}
