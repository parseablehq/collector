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

package client

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// k8s interface
type k8s interface {
	lister
	getter
}

// lister interface
type lister interface {
	ListPods(namespace, selector string) (*corev1.PodList, error)
}

// getter interface
type getter interface {
	GetPodLogs(pod corev1.Pod, podLogOptions corev1.PodLogOptions) ([]string, error)
}

var KubeClient k8s = getKubeClientset()

// k8s client struct
type client struct {
	*kubernetes.Clientset
}

// list pods method
func (c *client) ListPods(namespace, selector string) (*corev1.PodList, error) {
	pods, err := c.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{
		LabelSelector: selector,
	})
	if err != nil {
		return nil, err
	}

	return pods, nil
}

// get pods method
func (c *client) GetPodLogs(pod corev1.Pod, podLogOptions corev1.PodLogOptions) ([]string, error) {
	req := c.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOptions)
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return nil, err
	}

	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return nil, err
	}

	return strings.Split(buf.String(), "\n"), nil
}

func getKubeClientset() *client {

	var conf *rest.Config
	var err error

	if os.Getenv("RUN_LOCAL") == "true" {
		// for running locally
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		conf, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	} else {
		// creates the in-cluster config
		conf, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}
	cs, err := kubernetes.NewForConfig(conf)
	if err != nil {
		log.Printf("error in getting clientset from Kubeconfig: %v", err)
	}

	return &client{cs}
}
