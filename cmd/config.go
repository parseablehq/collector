package cmd

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type LogConfig struct {
	LogStreams []LogStream `yaml:"logStream"`
}

type LogStream struct {
	Name     string    `yaml:"name"`
	LogSpecs []LogSpec `yaml:"logSpec"`
}

type LogSpec struct {
	AddLabels   map[string]string `yaml:"addLabels"`
	Interval    string            `yaml:"collectionInterval"`
	CollectFrom CollectFrom       `yaml:"collectFrom"`
}

type CollectFrom struct {
	Namespace   string            `yaml:"namespace"`
	PodSelector map[string]string `yaml:"podSelector"`
}

func ReadConfig(path *string) (*LogConfig, error) {
	configfile, err := ioutil.ReadFile(*path)
	if err != nil {
		return nil, err
	}

	var logConfig LogConfig

	err = yaml.Unmarshal([]byte(configfile), &logConfig)
	if err != nil {
		return nil, err
	}
	return &logConfig, nil
}
