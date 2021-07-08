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

	node, err := kclient.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		fmt.Print("get node error", err)
	}
	fmt.Printf("get node success\n")

	addr := getUnLabelPod(node, kclient)
	if addr != "" {
		cmd := exec.Command("sh", "-c", "export BEE_ADDRESS="+addr)
		err := cmd.Start()
		pod, err := kclient.CoreV1().Pods("default").Get(context.TODO(), podName, metav1.GetOptions{})
		if err != nil {
			fmt.Print("get pod error", err)
		}
		fmt.Printf("get pod success\n")
		pod.ObjectMeta.Labels[addr] = "addr"
		for i := 0; i < 3; i++ {
			_, err := kclient.CoreV1().Pods("default").Update(context.TODO(), pod, metav1.UpdateOptions{})
			if err != nil {
				fmt.Printf("labelpod Update pod with new label err:%v=addr , %v\n", addr, err)
				time.Sleep(time.Second * 5)
				pod, err = kclient.CoreV1().Pods("default").Get(context.TODO(), podName, metav1.GetOptions{})
				if err != nil {
					fmt.Print("get pod error", err)
				}
				fmt.Printf("get pod success\n")
				continue
			} else {
				break
			}
		}
		fmt.Printf("label pod %v=addr\n", addr)
	} else {
		// require bee wallet address
		httpAddr := getBeeKey("http://10.1.66.146:8010/getAddressName")
		cmd := exec.Command("sh", "-c", "wget -P /home/bee/bee/file/ http://10.1.66.146:8010/getAddressFile/"+httpAddr+".tar.gz && "+
			"tar zxvf /home/bee/bee/file/"+httpAddr+".tar.gz -C /home/bee/bee/file/ && export BEE_ADDRESS="+httpAddr)
		err := cmd.Start()
		if err != nil {
			panic(err)
		}
		fmt.Printf("get addr success,%v\n", httpAddr)
		pod, err := kclient.CoreV1().Pods("default").Get(context.TODO(), podName, metav1.GetOptions{})
		if err != nil {
			fmt.Print("get pod error", err)
		}
		fmt.Printf("get pod success\n")
		pod.ObjectMeta.Labels[httpAddr] = "addr"
		for i := 0; i < 3; i++ {
			_, err := kclient.CoreV1().Pods("default").Update(context.TODO(), pod, metav1.UpdateOptions{})
			if err != nil {
				fmt.Printf("labelpod Update pod with new label err:%v=addr , %v\n", httpAddr, err)
				time.Sleep(time.Second * 5)
				pod, err = kclient.CoreV1().Pods("default").Get(context.TODO(), podName, metav1.GetOptions{})
				if err != nil {
					fmt.Print("get pod error", err)
				}
				fmt.Printf("get pod success\n")
				continue
			} else {
				break
			}
		}
		node.ObjectMeta.Labels[httpAddr] = "addr"
		for i := 0; i < 3; i++ {
			node, err := kclient.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
			if err != nil {
				fmt.Print("get node error", err)
			}
			fmt.Printf("get node success\n")
			node.ObjectMeta.Labels[httpAddr] = "addr"
			_, err = kclient.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
			if err != nil {
				fmt.Printf("labelpod Update node with new label err:%v=addr , %v\n", httpAddr, err)
				time.Sleep(time.Second * 5)
				node, err = kclient.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
				if err != nil {
					fmt.Print("get node error", err)
				}
				fmt.Printf("get node success\n")
				continue
			} else {
				break
			}
		}

		fmt.Printf("label pod & node %v=addr\n", httpAddr)
	}
}

func getUnLabelPod(node *v1.Node, kclient *kubernetes.Clientset) string {
	nodeLabels := node.ObjectMeta.Labels
	var pods *v1.PodList
	var err error
	for i, v := range nodeLabels {
		fmt.Printf("nodelabel %v=%v\n", i, v)
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
			fmt.Printf("no podslist\n")
			return i
		} else if len(pods.Items) == 0 {
			fmt.Printf("node has keys rest\n")
			return i
		}
	}
	fmt.Printf("no keys rest\n")
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
