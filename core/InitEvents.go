package core

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"sync"
)

// EventMapImpl 全局对象，存储所有Event
var EventMapImpl *EventMap

func init() {
	EventMapImpl = &EventMap{Data: new(sync.Map)}
}

type EventMap struct {
	Data *sync.Map // key:namespace_kind_name value: *v1.Event
}

func (this *EventMap) GetKey(event *v1.Event) string {
	key := fmt.Sprintf("%s_%s_%s", event.Namespace, event.InvolvedObject.Kind, event.InvolvedObject.Name)
	return key
}

// Add 添加
func (this *EventMap) Add(event *v1.Event) {
	EventMapImpl.Data.Store(this.GetKey(event), event)
}

// Delete 删除
func (this *EventMap) Delete(event *v1.Event) {
	EventMapImpl.Data.Delete(this.GetKey(event))
}

func (this *EventMap) GetMessage(ns string, kind string, name string) string {
	key := fmt.Sprintf("%s_%s_%s", ns, kind, name)
	if v, ok := this.Data.Load(key); ok {
		return v.(*v1.Event).Message
	}

	return ""
}

// EventHandler informer实现
type EventHandler struct{}

func (this *EventHandler) OnAdd(obj interface{}) {
	EventMapImpl.Add(obj.(*v1.Event))
}
func (this *EventHandler) OnUpdate(oldObj, newObj interface{}) {
	EventMapImpl.Add(newObj.(*v1.Event))
}
func (this *EventHandler) OnDelete(obj interface{}) {
	EventMapImpl.Delete(obj.(*v1.Event))
}
