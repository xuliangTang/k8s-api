package deployment

import (
	"context"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8sapi/core"
	"k8sapi/lib"
)

func RegHandlers(engine *gin.Engine) {
	engine.Handle("POST", "/deployment/scale/update", incrReplicas)
	engine.Handle("GET", "/api/deployments", getDeployments)
	engine.Handle("GET", "/api/deployment/pods", GetPodsByDeployment)
}

func getDeployments(gc *gin.Context) {
	ns := gc.DefaultQuery("namespace", "default")
	ret := ListAll(ns)

	lib.Success(gc, "success", ret)
}

func incrReplicas(gc *gin.Context) {
	req := struct {
		NameSpace  string `json:"ns" binding:"required,min=1"`
		Deployment string `json:"deployment" binding:"required,min=1"`
		Incr       *bool  `json:"incr" binding:"required"`
	}{}

	lib.CheckError(gc.ShouldBindJSON(&req))

	ctx := context.Background()
	scale, err := lib.K8sClient.AppsV1().Deployments(req.NameSpace).GetScale(ctx, req.Deployment, v1.GetOptions{})
	lib.CheckError(err)

	if *req.Incr {
		scale.Spec.Replicas++
	} else {
		if scale.Spec.Replicas > 0 {
			scale.Spec.Replicas--
		} else {
			scale.Spec.Replicas = 0
		}
	}

	_, err = lib.K8sClient.AppsV1().Deployments(req.NameSpace).UpdateScale(ctx, req.Deployment, scale, v1.UpdateOptions{})
	lib.CheckError(err)

	lib.Success(gc, "操作成功", nil)
}

func GetPodsByDeployment(gc *gin.Context) {
	ns := gc.DefaultQuery("ns", "default")
	depname := gc.DefaultQuery("deployment", "default")

	// 获取deployment
	dep, err := core.DepMapImpl.Find(ns, depname)
	lib.CheckError(err)

	// 获取所有的ReplicaSet
	rsList, err := core.RSMapImpl.ListByNs(ns)
	lib.CheckError(err)

	// 获取labels
	labels, err := GetListWatchRsLabelByDeployment(dep, rsList)
	lib.CheckError(err)

	// 根据labels获取pod
	lib.Success(gc, "success", GetPodsByLabels(ns, labels))
}
