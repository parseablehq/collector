package cmd

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type LogConfig struct {
	LogStreams []LogStream `yaml:"logStreams"`
}

type LogStream struct {
	Name            string            `yaml:"name"`
	AddLabels       map[string]string `yaml:"addLabels"`
	CollectInterval int               `yaml:"collectInterval"`
	CollectFrom     []struct {
		Namespace   string            `yaml:"namespace"`
		PodSelector map[string]string `yaml:"podSelector"`
	} `yaml:"collectFrom"`
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
