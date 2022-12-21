package deployment

import (
	"context"
	"fmt"
	v1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8sapi/lib"
	"strings"
)

func GetDepImages(deployment v1.Deployment) string {
	return GetPodImages(deployment.Spec.Template.Spec.Containers)
}

func GetPodImages(cs []coreV1.Container) string {
	var images strings.Builder
	images.WriteString(cs[0].Image)
	if lenImage := len(cs); lenImage > 1 {
		images.WriteString(fmt.Sprintf("和%d个镜像", lenImage))
	}

	return images.String()
}

func GetLabels(labels map[string]string) string {
	var labelStr strings.Builder

	for k, v := range labels {
		if labelStr.Len() != 0 {
			labelStr.WriteString(",")
		}
		labelStr.WriteString(fmt.Sprintf("%s=%s", k, v))
	}

	return labelStr.String()
}

// GetListWatchRsLabelByDeployment list-watch方式 根据deployment获取当前ReplicaSet的标签
func GetListWatchRsLabelByDeployment(deployment *v1.Deployment, rsList []*v1.ReplicaSet) ([]map[string]string, error) {
	labels := make([]map[string]string, 0)
	for _, rs := range rsList {
		if IsRsByDeployment(rs, deployment) {
			selector, err := metaV1.LabelSelectorAsMap(rs.Spec.Selector)
			if err != nil {
				return nil, err
			}
			labels = append(labels, selector)
		}
	}

	return labels, nil
}

// GetRsLabelByDeployment 根据deployment获取当前ReplicaSet的标签
func GetRsLabelByDeployment(deployment *v1.Deployment) string {
	dep, err := lib.K8sClient.AppsV1().Deployments("default").Get(context.Background(), deployment.Name, metaV1.GetOptions{})
	lib.CheckError(err)

	selector, err := metaV1.LabelSelectorAsSelector(dep.Spec.Selector)
	lib.CheckError(err)
	listOpt := metaV1.ListOptions{
		LabelSelector: selector.String(),
	}
	rsList, err := lib.K8sClient.AppsV1().ReplicaSets("default").List(context.Background(), listOpt)
	lib.CheckError(err)

	for _, rs := range rsList.Items {
		if IsCurrentRsByDeployment(&rs, deployment) {
			selector, err := metaV1.LabelSelectorAsSelector(rs.Spec.Selector)
			if err != nil {
				return ""
			}
			return selector.String()
		}
	}

	return ""
}

// IsRsByDeployment 判断rs是否属于当前deployment
func IsRsByDeployment(set *v1.ReplicaSet, deployment *v1.Deployment) bool {
	for _, rf := range set.OwnerReferences {
		if rf.Kind == "Deployment" && rf.Name == deployment.Name {
			return true
		}
	}

	return false
}

// IsCurrentRsByDeployment 判断rs是否属于当前deployment最新的一条
func IsCurrentRsByDeployment(set *v1.ReplicaSet, deployment *v1.Deployment) bool {
	if set.ObjectMeta.Annotations["deployment.kubernetes.io/revision"] != deployment.ObjectMeta.Annotations["deployment.kubernetes.io/revision"] {
		return false
	}

	return IsRsByDeployment(set, deployment)
}

// GetPodIsReady 评估Pod是否就绪
func GetPodIsReady(pod *coreV1.Pod) bool {
	// 所有容器是否就绪
	for _, condition := range pod.Status.Conditions {
		if condition.Type == "ContainersReady" && condition.Status != "True" {
			return false
		}
	}

	// readinessGates
	for _, rg := range pod.Spec.ReadinessGates {
		for _, condition := range pod.Status.Conditions {
			if condition.Type == rg.ConditionType && condition.Status != "True" {
				return false
			}
		}
	}

	return true
}
