package lib

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

func GetK8sClientSet() (*kubernetes.Clientset, error) {
	k8s_config_in_cluster := os.Getenv("FILTAB_K8S_CONFIG_IN_CLUSTER")

	if k8s_config_in_cluster == "true" {
		// creates the in-cluster config
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		// creates the clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}

		return clientset, err
	} else {
		home := homeDir()
		kubeconfig := filepath.Join(home, ".kube", "config")
		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}

		// create the clientset
		return kubernetes.NewForConfig(config)
	}
}
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
