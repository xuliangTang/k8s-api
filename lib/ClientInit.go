package lib

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)

var K8sClient *kubernetes.Clientset

func init() {
	config := &rest.Config{
		Host:        "110.41.142.160:8009",
		BearerToken: "",
	}
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	K8sClient = c
}
