package deployment

import (
	"context"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8sapi/core"
	"k8sapi/lib"
	"k8sapi/model"
	"strings"
)

func GetPodMessage(pod *corev1.Pod) string {
	var message strings.Builder
	for _, condition := range pod.Status.Conditions {
		if condition.Status != "True" {
			message.WriteString(condition.Message)
		}
	}
	return message.String()
}

// GetPodsByLabels 获取pods DTO 把原生的 pod 对象转换为自己的实体对象
func GetPodsByLabels(ns string, labels map[string]string) (pods []*model.PodModel) {
	podList, err := core.PodMapImpl.ListByLabels(ns, labels)
	lib.CheckError(err)

	pods = make([]*model.PodModel, len(podList))

	for i, pod := range podList {
		pods[i] = &model.PodModel{
			Name:     pod.Name,
			NodeName: pod.Spec.NodeName,
			Images:   GetPodImages(pod.Spec.Containers),
			Phase:    string(pod.Status.Phase),
			IsReady:  GetPodIsReady(pod),
			// Message:   GetPodMessage(pod),
			Message:   core.EventMapImpl.GetMessage(pod.Namespace, "Pod", pod.Name),
			CreatedAt: pod.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
	}

	return
}

func GetPodsByDep(namespace string, dep *v1.Deployment) (pods []*model.PodModel) {
	ctx := context.Background()
	listOpt := metav1.ListOptions{
		//LabelSelector: GetLabels(dep.Spec.Selector.MatchLabels),
		LabelSelector: GetRsLabelByDeployment(dep),
	}

	podList, err := lib.K8sClient.CoreV1().Pods(namespace).List(ctx, listOpt)
	lib.CheckError(err)

	pods = make([]*model.PodModel, len(podList.Items))

	for i, pod := range podList.Items {
		pods[i] = &model.PodModel{
			Name:      pod.Name,
			NodeName:  pod.Spec.NodeName,
			Images:    GetPodImages(pod.Spec.Containers),
			CreatedAt: pod.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
	}

	return pods
}

func DepDetail(namespace string, name string) (ret model.DeploymentModel) {
	ctx := context.Background()
	getOpt := metav1.GetOptions{}
	dep, err := lib.K8sClient.AppsV1().Deployments(namespace).Get(ctx, name, getOpt)
	lib.CheckError(err)

	ret.Name = dep.Name
	ret.Images = GetDepImages(*dep)
	ret.NameSpace = dep.Namespace
	ret.CreatedAt = dep.CreationTimestamp.Format("2006-01-02 15:04:05")
	ret.Pods = GetPodsByDep(namespace, dep)
	ret.Replicas = [3]int32{dep.Status.Replicas, dep.Status.AvailableReplicas, dep.Status.UnavailableReplicas}

	return
}
