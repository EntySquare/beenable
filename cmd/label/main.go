package main

import (
	"beenable/lib"
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

func main() {
	kclient, err := lib.GetK8sClientSet()
	if err != nil {
		panic(err)
	}
	podName := os.Getenv("JOB_POD_NAME")

	pod, err := kclient.CoreV1().Pods("default").Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		fmt.Print("get pod error", err)
	}
	var podLabelMap = make(map[string]string)
	var nodeLabelMap = make(map[string]string)

	addr := os.Getenv("BEE_ADDRESS")
	podLabelMap[addr] = "address"
	pod.SetLabels(podLabelMap)
	node, err := kclient.CoreV1().Nodes().Get(context.TODO(), pod.Spec.NodeName, metav1.GetOptions{})
	if err != nil {
		fmt.Print("get node error", err)
	}
	nodeLabelMap[addr] = "address"
	node.SetLabels(nodeLabelMap)
}
