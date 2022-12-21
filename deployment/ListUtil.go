package deployment

import (
	"k8sapi/core"
	"k8sapi/lib"
	"k8sapi/model"
)

// ListAll 显示所有
func ListAll(namespace string) (ret []*model.DeploymentModel) {
	depList, err := core.DepMapImpl.ListByNs(namespace)
	lib.CheckError(err)

	for _, item := range depList { //遍历所有deployment
		ret = append(ret, &model.DeploymentModel{
			Name:     item.Name,
			Replicas: [3]int32{item.Status.Replicas, item.Status.AvailableReplicas, item.Status.UnavailableReplicas},
			Images:   GetDepImages(*item),
		})
	}

	return
}
