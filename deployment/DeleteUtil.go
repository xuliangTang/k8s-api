package deployment

import (
	"context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8sapi/lib"
)

func DeletePod(ns, name string) {
	err := lib.K8sClient.CoreV1().Pods(ns).Delete(context.Background(), name, v1.DeleteOptions{})
	lib.CheckError(err)
}

func DeleteDeployment(ns, name string) {
	err := lib.K8sClient.AppsV1().Deployments(ns).Delete(context.Background(), name, v1.DeleteOptions{})
	lib.CheckError(err)
}
