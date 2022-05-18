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
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	DEFAULT_LOG_COLLECT_INTERVAL  = "1m"
	ENV_PARSEABLE_SERVER_URL      = "PARSEABLE_URL"
	ENV_PARSEABLE_SERVER_USERNAME = "PARSEABLE_USERNAME"
	ENV_PARSEABLE_SERVER_PASSWORD = "PARSEABLE_PASSWORD"
)

type CollectorConfig struct {
	Server     string
	Username   string
	Password   string
	LogStreams []LogStream `yaml:"logStreams"`
}

type LogStream struct {
	Name            string            `yaml:"name"`
	Tags            map[string]string `yaml:"tags"`
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
	if err := logConfig.SetCreds(); err != nil {
		return nil, err
	}

	return &logConfig, nil
}

func (logConfig *CollectorConfig) SetCreds() error {
	var ok bool
	logConfig.Server, ok = os.LookupEnv(ENV_PARSEABLE_SERVER_URL)
	if !ok {
		return fmt.Errorf("%s environment variable is not set", ENV_PARSEABLE_SERVER_URL)
	}
	logConfig.Username, ok = os.LookupEnv(ENV_PARSEABLE_SERVER_USERNAME)
	if !ok {
		logConfig.Password, ok = os.LookupEnv(ENV_PARSEABLE_SERVER_PASSWORD)
		if !ok {
			log.Info("Parseable credentials are not set as environment variables. Sending unauthenticated requests to Parseable server.")
		}
	}
	return nil
}

func (logConfig *CollectorConfig) ensureDefaults() {
	for _, logStream := range logConfig.LogStreams {
		if logStream.CollectInterval == "" {
			logStream.CollectInterval = DEFAULT_LOG_COLLECT_INTERVAL
		}
	}
}
