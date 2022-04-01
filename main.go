package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// logMessage is the CRI internal log type.
type logMessage struct {
	timestamp time.Time
	stream    runtimeapi.LogStreamType
	log       []byte
}

func main() {
	clientset := GetKubeClientset()
	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), v1.ListOptions{
		LabelSelector: "app=cloudsql-proxy",
	})
	if err != nil {
		return
	}
	for _, p := range pods.Items {
		fmt.Println(p.GetName())
		getPodLogs(p)

	}

}

func getPodLogs(pod corev1.Pod) string {
	podLogOpts := corev1.PodLogOptions{
		Timestamps: true,
	}
	clientset := GetKubeClientset()

	req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)

	podLogs, err := req.Stream(context.TODO())

	if err != nil {
		return "error in opening stream"
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "error in copy information from podLogs to buf"
	}
	str := buf.String()

	newStr := strings.Split(str, "\n")
	//a := newStr[len(newStr)-1]
	fmt.Println(newStr[112])
	return str
}

func GetKubeClientset() *kubernetes.Clientset {
	var conf *rest.Config

	conf, err := clientcmd.BuildConfigFromFlags("", os.Getenv("HOME")+"/.kube/config")
	if err != nil {
		log.Printf("error in getting Kubeconfig: %v", err)
	}

	cs, err := kubernetes.NewForConfig(conf)
	if err != nil {
		log.Printf("error in getting clientset from Kubeconfig: %v", err)
	}

	return cs
}
