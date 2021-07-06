package core

import (
	"beenable/lib"
	"bytes"
	"context"
	"fmt"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/kubernetes"
	"log"
	"net/http"
	"os"
	"time"
)

// static strategy
type StaticStrategy struct {
	SwapEndpoint    string
	SwapEnable      string
	SwapGas         string
	SwapInitDeposit string
	DebugApiEnable  string
	NetworkId       string
	Mainnet         string
	FullNode        string
	Verbosity       string
	ClefEnable      string
	ImageName       string
	Password        string
	DataDir         string

	client *kubernetes.Clientset
}

type BeeAddr struct {
	Address []string
}

func NewStaticStrategy(swapEndpoint, swapEnable, swapGas, swapInitDeposit, debugApiEnable, networkId, mainnet, fullNode, verbosity, clefEnable, imageName, password, dataDir string) *StaticStrategy {

	kclient, err := lib.GetK8sClientSet()
	if err != nil {
		panic(err)
	}

	var s = &StaticStrategy{
		SwapEndpoint:    swapEndpoint,
		SwapEnable:      swapEnable,
		SwapGas:         swapGas,
		SwapInitDeposit: swapInitDeposit,
		DebugApiEnable:  debugApiEnable,
		NetworkId:       networkId,
		Mainnet:         mainnet,
		FullNode:        fullNode,
		Verbosity:       verbosity,
		ClefEnable:      clefEnable,
		ImageName:       imageName,
		Password:        password,
		DataDir:         dataDir,
		client:          kclient,
	}
	return s
}

func (s *StaticStrategy) Run() error {
	s.start()
	return nil
}

func (s *StaticStrategy) start() {
	sf := lib.StartBeeAffinity()
	limitList := corev1.ResourceList{}
	requestList := corev1.ResourceList{}
	//limitList["cpu"] = resource.MustParse("10%")
	//requestList["cpu"] = resource.MustParse("10%")
	limitList["memory"] = resource.MustParse("1Gi")
	requestList["memory"] = resource.MustParse("200m")
	// require bee wallet address
	addr := getBeeKey("http://192.168.2.12/getAddressName")
	//labelKey := key[0:10]
	jbname := "entysquare-bee-job-" + addr + "-" + rand.String(10)
	fmt.Println("run job : " + jbname)
	// random port
	port1 := generateRangeNum(10001, 20000)
	// make address file
	//f, err := os.OpenFile("bee_address", os.O_RDWR|os.O_TRUNC, 0600)
	//defer f.Close()
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	_, err = f.Write([]byte(addr))
	//}
	port2 := port1 + 1
	port3 := port2 + 1
	jb := lib.GetJob(jbname, 1, 10000, sf, limitList, requestList, s.SwapEndpoint, s.SwapEnable, s.SwapGas, s.SwapInitDeposit,
		s.DebugApiEnable, s.NetworkId, s.Mainnet, s.FullNode, s.Verbosity, s.ClefEnable, s.ImageName, s.Password, s.DataDir, addr, port1, port2, port3)
	_ = os.Setenv("address", addr)
	_, err := s.client.BatchV1().Jobs("default").Create(context.TODO(), jb, metav1.CreateOptions{})
	if err != nil {
		log.Fatal("create job error", err)
	}
}

func (s *StaticStrategy) Restart() error {
	s.start()
	return nil
}

func generateRangeNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max - min)
	randNum = randNum + min
	fmt.Printf("rand is %v\n", randNum)
	return randNum
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
	return result.String()
}
