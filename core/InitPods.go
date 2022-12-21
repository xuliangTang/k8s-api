package core

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"log"
	"reflect"
	"sync"
)

// PodMapImpl 全局对象，存储所有Pod
var PodMapImpl *PodMap

func init() {
	PodMapImpl = &PodMap{Data: new(sync.Map)}
}

type PodMap struct {
	Data *sync.Map // key:namespace value:[]*v1.Pod
}

// Add 添加
func (this *PodMap) Add(pod *v1.Pod) {
	if podList, ok := this.Data.Load(pod.Namespace); ok {
		newList := append(podList.([]*v1.Pod), pod)
		this.Data.Store(pod.Namespace, newList)
	} else {
		this.Data.Store(pod.Namespace, []*v1.Pod{pod})
	}
}

// Update 更新
func (this *PodMap) Update(pod *v1.Pod) error {
	if podList, ok := this.Data.Load(pod.Namespace); ok {
		podList := podList.([]*v1.Pod)
		for i, p := range podList {
			if p.Name == pod.Name {
				podList[i] = pod
				break
			}
		}
		return nil
	}

	return fmt.Errorf("replicaset [%s] not found", pod.Name)
}

// Delete 删除
func (this *PodMap) Delete(pod *v1.Pod) {
	if podList, ok := this.Data.Load(pod.Namespace); ok {
		podList := podList.([]*v1.Pod)
		for i, p := range podList {
			if p.Name == pod.Name {
				newRSList := append(podList[:i], podList[i+1:]...)
				this.Data.Store(pod.Namespace, newRSList)
				break
			}
		}
	}
}

// ListByLabels 根据标签获取Pod列表
func (this *PodMap) ListByLabels(ns string, labels map[string]string) ([]*v1.Pod, error) {
	ret := make([]*v1.Pod, 0)
	if podList, ok := this.Data.Load(ns); ok {
		podList := podList.([]*v1.Pod)
		for _, p := range podList {
			// 判断标签完全匹配
			if reflect.DeepEqual(p.Labels, labels) {
				ret = append(ret, p)
			}
		}

		return ret, nil
	}

	return nil, fmt.Errorf("pods not found")
}

// PodHandler informer实现
type PodHandler struct{}

func (this *PodHandler) OnAdd(obj interface{}) {
	PodMapImpl.Add(obj.(*v1.Pod))
}
func (this *PodHandler) OnUpdate(oldObj, newObj interface{}) {
	err := PodMapImpl.Update(newObj.(*v1.Pod))
	if err != nil {
		log.Println(err)
	}
}
func (this *PodHandler) OnDelete(obj interface{}) {
	PodMapImpl.Delete(obj.(*v1.Pod))
}
