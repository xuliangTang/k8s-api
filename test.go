package main

import (
	"context"
	"fmt"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8sapi/lib"
)

func main() {
	dep, _ := lib.K8sClient.AppsV1().Deployments("default").Get(context.Background(), "two-myngx", v1.GetOptions{})

	selector, _ := v1.LabelSelectorAsSelector(dep.Spec.Selector)

	listOpt := v1.ListOptions{
		LabelSelector: selector.String(),
	}
	rsList, _ := lib.K8sClient.AppsV1().ReplicaSets("default").List(context.Background(), listOpt)

	for _, rs := range rsList.Items {
		fmt.Println(rs.Name)
		fmt.Println(IsCurrentRsByDep(&rs, dep))
	}
}

// pod-template-hash=54bfc97676
// pod-template-hash=54bfc97676

func IsCurrentRsByDep(set *appsV1.ReplicaSet, deployment *appsV1.Deployment) bool {
	if set.ObjectMeta.Annotations["deployment.kubernetes.io/revision"] != deployment.ObjectMeta.Annotations["deployment.kubernetes.io/revision"] {
		return false
	}

	for _, rf := range set.OwnerReferences {
		if rf.Kind == "Deployment" && rf.Name == deployment.Name {
			return true
		}
	}

	return false
}
