package deployment

import (
	"context"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8sapi/core"
	"k8sapi/lib"
)

func RegHandlers(engine *gin.Engine) {
	// 修改副本数
	engine.Handle("POST", "/deployment/scale/update", incrReplicas)
	// 获取deployment列表
	engine.Handle("GET", "/api/deployments", getDeployments)
	// 获取deployment的pod列表
	engine.Handle("GET", "/api/deployment/pods", GetPodsByDeployment)
	// 获取pod json
	engine.Handle("GET", "/api/pod", FindPod)
	// 删除pod
	engine.Handle("DELETE", "/api/pod", DeletePodApi)
	// 创建deployment
	engine.Handle("POST", "/api/deployment", CreateDeploymentApi)
	// 删除deployment
	engine.Handle("DELETE", "/api/deployment", DeleteDeploymentApi)
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

func FindPod(gc *gin.Context) {
	ns := gc.DefaultQuery("ns", "default")
	podName := gc.DefaultQuery("pod", "default")

	pod := core.PodMapImpl.Find(ns, podName)

	gc.JSON(200, pod)
}

func DeletePodApi(gc *gin.Context) {
	ns := gc.DefaultQuery("ns", "default")
	podName := gc.DefaultQuery("pod", "default")

	DeletePod(ns, podName)
	lib.Success(gc, "success", nil)
}

func CreateDeploymentApi(gc *gin.Context) {
	req := &DeploymentCreateReq{}
	lib.CheckError(gc.ShouldBind(req))
	lib.CheckError(CreateDeployment(req))

	gc.Redirect(301, "/deployments")
}

func DeleteDeploymentApi(gc *gin.Context) {
	ns := gc.DefaultQuery("ns", "default")
	deploymentName := gc.DefaultQuery("deployment", "default")

	DeleteDeployment(ns, deploymentName)
	lib.Success(gc, "success", nil)
}
