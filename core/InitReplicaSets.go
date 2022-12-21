package core

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
	"log"
	"sync"
)

// RSMapImpl 全局对象，存储所有ReplicaSets
var RSMapImpl *ReplicaSetMap

func init() {
	RSMapImpl = &ReplicaSetMap{Data: new(sync.Map)}
}

type ReplicaSetMap struct {
	Data *sync.Map // key:namespace value:[]*v1.ReplicaSet
}

// Add 添加
func (this *ReplicaSetMap) Add(set *v1.ReplicaSet) {
	if rsList, ok := this.Data.Load(set.Namespace); ok {
		newList := append(rsList.([]*v1.ReplicaSet), set)
		this.Data.Store(set.Namespace, newList)
	} else {
		this.Data.Store(set.Namespace, []*v1.ReplicaSet{set})
	}
}

// Update 更新
func (this *ReplicaSetMap) Update(set *v1.ReplicaSet) error {
	if rsList, ok := this.Data.Load(set.Namespace); ok {
		rsList := rsList.([]*v1.ReplicaSet)
		for i, rs := range rsList {
			if rs.Name == set.Name {
				rsList[i] = set
				break
			}
		}
		return nil
	}

	return fmt.Errorf("replicaset [%s] not found", set.Name)
}

// Delete 删除
func (this *ReplicaSetMap) Delete(set *v1.ReplicaSet) {
	if rsList, ok := this.Data.Load(set.Namespace); ok {
		rsList := rsList.([]*v1.ReplicaSet)
		for i, rs := range rsList {
			if rs.Name == set.Name {
				newRSList := append(rsList[:i], rsList[i+1:]...)
				this.Data.Store(set.Namespace, newRSList)
				break
			}
		}
	}
}

// ListByNs 获取列表
func (this *ReplicaSetMap) ListByNs(ns string) ([]*v1.ReplicaSet, error) {
	if rsList, ok := this.Data.Load(ns); ok {
		return rsList.([]*v1.ReplicaSet), nil
	}

	return nil, fmt.Errorf("record not found")
}

// ReplicaSetHandler informer实现
type ReplicaSetHandler struct{}

func (this *ReplicaSetHandler) OnAdd(obj interface{}) {
	RSMapImpl.Add(obj.(*v1.ReplicaSet))
}
func (this *ReplicaSetHandler) OnUpdate(oldObj, newObj interface{}) {
	err := RSMapImpl.Update(newObj.(*v1.ReplicaSet))
	if err != nil {
		log.Println(err)
	}
}
func (this *ReplicaSetHandler) OnDelete(obj interface{}) {
	RSMapImpl.Delete(obj.(*v1.ReplicaSet))
}
