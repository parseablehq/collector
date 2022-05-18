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

package main

import (
	"flag"
	"kube-collector/cmd"

	log "github.com/sirupsen/logrus"

	"os"
	"sync"
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
		go func(stream cmd.LogStream) {
			defer wg.Done()
			cmd.RunKubeCollector(config.Server, config.Username, config.Password, &stream)
		}(stream)
	}
	wg.Wait()
}
