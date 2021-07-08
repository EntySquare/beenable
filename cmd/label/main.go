package main

import (
	"beenable/lib"
	"bytes"
	"context"
	"fmt"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	kclient, err := lib.GetK8sClientSet()
	if err != nil {
		panic(err)
	}
	podName := os.Getenv("JOB_POD_NAME")
	nodeName := os.Getenv("JOB_NODE_NAME")

	var podLabelMap = make(map[string]string)
	var nodeLabelMap = make(map[string]string)

	pod, err := kclient.CoreV1().Pods("default").Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		fmt.Print("get pod error", err)
	}
	fmt.Printf("get pod success\n")

	node, err := kclient.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		fmt.Print("get node error", err)
	}
	fmt.Printf("get node success\n")

	addr := getUnLabelPod(node, kclient)
	if addr != "" {
		_ = os.Setenv("BEE_ADDRESS", addr)
		podLabelMap[addr] = "addr"
		pod.SetLabels(podLabelMap)

		fmt.Printf("label pod %v=addr\n", addr)
	} else {
		// require bee wallet address
		httpAddr := getBeeKey("http://10.1.66.146:8010/getAddressName")
		cmd := exec.Command("sh", "-c", "wget -P /home/bee/bee/file"+httpAddr+"http://10.1.66.146:8010/getAddressFile/"+httpAddr+".tar.gz && tar zxvf "+httpAddr+".tar.gz")
		_ = os.Setenv("BEE_ADDRESS", httpAddr)
		err := cmd.Start()
		time.Sleep(time.Second * 15)
		if err != nil {
			panic(err)
		}
		fmt.Printf("get addr success,%v\n", httpAddr)
		podLabelMap[httpAddr] = "addr"
		pod.SetLabels(podLabelMap)
		nodeLabelMap[httpAddr] = "addr"
		node.SetLabels(podLabelMap)

		fmt.Printf("label pod & node %v=addr\n", httpAddr)
	}
}

func getUnLabelPod(node *v1.Node, kclient *kubernetes.Clientset) string {
	nodeLabels := node.GetLabels()
	var pods *v1.PodList
	var err error
	for i, v := range nodeLabels {
		if v != "addr" {
			continue
		}
		for j := 0; j < 3; j++ {
			pods, err = kclient.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{
				LabelSelector: i + "=addr",
			})
			if err != nil {
				break
			}
		}
		if pods == nil {
			return i
		}
	}
	return ""
}

func getBeeKey(url string) string {

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	fmt.Printf("get address success %v\n", result.String())
	return result.String()
}
