package cmd

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type LogConfig struct {
	LogStreams []struct {
		Name            string `yaml:"name"`
		CollectInterval int    `yaml:"collectInterval"`
		CollectFrom     []struct {
			Namespace   string            `yaml:"namespace"`
			PodSelector map[string]string `yaml:"podSelector"`
		} `yaml:"collectFrom"`
	} `yaml:"logStreams"`
}

func ReadConfig(path string) (*LogConfig, error) {
	configfile, err := ioutil.ReadFile(path)
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
