package k8s

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"strings"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
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

// k8s client struct
type k8sClient struct {
	*kubernetes.Clientset
}

// init k8s interface
var K8s k8s = getKubeClientset()

// list pods method
func (c *k8sClient) ListPods(namespace, selector string) (*corev1.PodList, error) {
	pods, err := c.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{
		LabelSelector: selector,
	})
	if err != nil {
		return nil, err
	}

	return pods, nil
}

// get pods method
func (c *k8sClient) GetPodLogs(pod corev1.Pod, podLogOptions corev1.PodLogOptions) ([]string, error) {
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

func getKubeClientset() *k8sClient {
	var conf *rest.Config

	conf, err := clientcmd.BuildConfigFromFlags("", os.Getenv("HOME")+"/.kube/config")
	if err != nil {
		log.Printf("error in getting Kubeconfig: %v", err)
	}

	cs, err := kubernetes.NewForConfig(conf)
	if err != nil {
		log.Printf("error in getting clientset from Kubeconfig: %v", err)
	}

	return &k8sClient{cs}
}
