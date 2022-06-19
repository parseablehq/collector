package cmd

import (
	"kube-collector/pkg/client"
	"kube-collector/pkg/store"
	"kube-collector/pkg/utils"
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
