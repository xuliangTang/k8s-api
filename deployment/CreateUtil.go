package deployment

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8sapi/lib"
	"log"
)

func Create() {
	ctx := context.Background()

	listopt := metav1.ListOptions{}
	depList, err := lib.K8sClient.AppsV1().Deployments("myweb").List(ctx, listopt)

	for _, item := range depList.Items { //遍历所有deployment
		fmt.Println(item.Name)
	}
	ngxDep := &v1.Deployment{} //我们要创建的deployment
	b, _ := ioutil.ReadFile("yamls/nginx.yaml")
	ngxJson, _ := yaml.ToJSON(b)
	lib.CheckError(json.Unmarshal(ngxJson, ngxDep))

	createopt := metav1.CreateOptions{}

	_, err = lib.K8sClient.AppsV1().Deployments("myweb").
		Create(ctx, ngxDep, createopt)

	if err != nil {
		log.Fatal(err)
	}
}

type DeploymentCreateReq struct {
	Name  string `form:"name" binding:"required"`
	Image string `form:"image" binding:"required"`
}

func CreateDeployment(req *DeploymentCreateReq) error {
	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: "default",
		},
		Spec: v1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: genLabels(req.Name),
			},
			Template: coreV1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: genLabels(req.Name),
				},
				Spec: coreV1.PodSpec{
					Containers: genContainers(req),
				},
			},
		},
	}

	_, err := lib.K8sClient.AppsV1().Deployments("default").Create(context.Background(), deployment, metav1.CreateOptions{})

	return err
}

// 生成标签配置
func genLabels(name string) map[string]string {
	return map[string]string{
		"app": name,
	}
}

// 生成容器配置
func genContainers(req *DeploymentCreateReq) []coreV1.Container {
	containers := make([]coreV1.Container, 1)
	containers[0] = coreV1.Container{
		Name:  req.Name,
		Image: req.Image,
	}
	return containers
}
