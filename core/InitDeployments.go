package core

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8sapi/lib"
	"log"
	"sync"
)

// DepMapImpl 全局对象，存储所有deployments
var DepMapImpl *DeploymentMap

func init() {
	DepMapImpl = &DeploymentMap{Data: new(sync.Map)}
}

type DeploymentMap struct {
	Data *sync.Map // key:namespace value:[]*v1.Deployments
}

// Add 添加
func (this *DeploymentMap) Add(deployment *v1.Deployment) {
	if depList, ok := this.Data.Load(deployment.Namespace); ok {
		depList = append(depList.([]*v1.Deployment), deployment)
		this.Data.Store(deployment.Namespace, depList)
	} else {
		this.Data.Store(deployment.Namespace, []*v1.Deployment{deployment})
	}
}

// ListByNs 获取列表
func (this *DeploymentMap) ListByNs(namespace string) ([]*v1.Deployment, error) {
	if depList, ok := this.Data.Load(namespace); ok {
		return depList.([]*v1.Deployment), nil
	}

	return nil, fmt.Errorf("record not found")
}

func (this *DeploymentMap) Find(ns string, depName string) (*v1.Deployment, error) {
	if depList, ok := this.Data.Load(ns); ok {
		for _, dep := range depList.([]*v1.Deployment) {
			if dep.Name == depName {
				return dep, nil
			}
		}
	}

	return nil, fmt.Errorf("record not found")
}

// Update 更新
func (this *DeploymentMap) Update(deployment *v1.Deployment) error {
	if depList, ok := this.Data.Load(deployment.Namespace); ok {
		depList := depList.([]*v1.Deployment)
		for i, dep := range depList {
			if dep.Name == deployment.Name {
				depList[i] = deployment
				break
			}
		}
		return nil
	}

	return fmt.Errorf("deployment [%s] not found", deployment.Name)
}

// Delete 删除
func (this *DeploymentMap) Delete(deployment *v1.Deployment) {
	if depList, ok := this.Data.Load(deployment.Namespace); ok {
		depList := depList.([]*v1.Deployment)
		for i, dep := range depList {
			if dep.Name == deployment.Name {
				newDepList := append(depList[:i], depList[i+1:]...)
				this.Data.Store(deployment.Namespace, newDepList)
				break
			}
		}
	}
}

// DepHandler informer实现
type DepHandler struct{}

func (this *DepHandler) OnAdd(obj interface{}) {
	DepMapImpl.Add(obj.(*v1.Deployment))
}
func (this *DepHandler) OnUpdate(oldObj, newObj interface{}) {
	err := DepMapImpl.Update(newObj.(*v1.Deployment))
	if err != nil {
		log.Println(err)
	}
}
func (this *DepHandler) OnDelete(obj interface{}) {
	DepMapImpl.Delete(obj.(*v1.Deployment))
}

// InitDeployments 执行监听
func InitDeployments() {
	informerFactory := informers.NewSharedInformerFactory(lib.K8sClient, 0)

	depInformer := informerFactory.Apps().V1().Deployments()
	depInformer.Informer().AddEventHandler(&DepHandler{})

	podInformer := informerFactory.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(&PodHandler{})

	replicaSetInformer := informerFactory.Apps().V1().ReplicaSets()
	replicaSetInformer.Informer().AddEventHandler(&ReplicaSetHandler{})

	informerFactory.Start(wait.NeverStop)
}
