package core

import (
	"beenable/lib"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/kubernetes"
	"log"
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
	limitList["cpu"] = resource.MustParse("2000m")
	requestList["cpu"] = resource.MustParse("2000m")
	limitList["memory"] = resource.MustParse("25Gi")
	requestList["memory"] = resource.MustParse("8Gi")
	jbname := "entysquare-bee-job-" + "-" + rand.String(10)
	fmt.Println("run job : " + jbname)
	jb := lib.GetJob(jbname, 1, 10000, sf, limitList, requestList, s.SwapEndpoint, s.SwapEnable, s.SwapGas, s.SwapInitDeposit,
		s.DebugApiEnable, s.NetworkId, s.Mainnet, s.FullNode, s.Verbosity, s.ClefEnable, s.ImageName, s.Password, s.DataDir)
	_, err := s.client.BatchV1().Jobs("default").Create(context.TODO(), jb, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
}

// dynamic strategy
type DynamicStrategy struct {
}

func NewDynamicStrategy() *DynamicStrategy {
	return nil
}

func (s *DynamicStrategy) Run() error {
	return nil
}

//func (s *StaticStrategy) testPlot() {
//	sf := lib.StartTestAffinity()
//	limitList := corev1.ResourceList{}
//	requestList := corev1.ResourceList{}
//	limitList["cpu"] = resource.MustParse("2000m")
//	requestList["cpu"] = resource.MustParse("2000m")
//	limitList["memory"] = resource.MustParse("25Gi")
//	requestList["memory"] = resource.MustParse("8Gi")
//	farmer := s.FarmerKey[:8]
//	jbname := "entysquare-k-" + s.K + "-job-plot-farmer-" + farmer + "-" + rand.String(5)
//	fmt.Println("run job : " + jbname)
//	jb := lib.GetJob(jbname, 1, 10000, sf, limitList, requestList, s.FarmerKey,
//		s.PoolKey, s.UserDir, s.ImageName, s.K, s.ReportIp, s.ReportPort)
//	_, err := s.client.BatchV1().Jobs("default").Create(context.TODO(), jb, metav1.CreateOptions{})
//	if err != nil {
//		log.Fatal(err)
//	}
//}
