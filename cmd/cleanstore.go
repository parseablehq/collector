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
	"collector/pkg/client"
	"collector/pkg/store"
	"collector/pkg/utils"
	"time"

	log "github.com/sirupsen/logrus"
)

func ExecCleanStore(namespace, selector string) {
	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		cleanStore(namespace, selector)
	}
}

func cleanStore(namespace, selector string) {

	pods, err := client.KubeClient.ListPods(namespace, selector)
	if err != nil {
		log.Error(err)
		return
	}

	// current state of pod names
	var currentPodNames []string

	for _, podName := range pods.Items {
		currentPodNames = append(currentPodNames, podName.GetName())
	}

	// range on store and if podName is not present in currentPodNames
	// delete podName from store
	for podName := range store.PoNameTime {
		if !utils.ContainsString(currentPodNames, podName) {
			store.DeletePodName(podName)
		}
	}
}
