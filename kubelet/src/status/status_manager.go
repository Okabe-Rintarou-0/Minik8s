package status

import (
	"minik8s/apiObject"
	"minik8s/kubelet/src/runtime/cache"
	"minik8s/kubelet/src/types"
)

type Manager interface {
	GetPod(podUID types.UID) *apiObject.Pod
	UpdatePod(podUID types.UID, newPod *apiObject.Pod)
	DeletePod(podUID types.UID)
}

type manager struct {
	podCache cache.Cache
}

func (m *manager) GetPod(podUID types.UID) *apiObject.Pod {
	if podVal := m.podCache.Get(podUID); podVal != nil {
		return podVal.(*apiObject.Pod)
	}
	return nil
}

func (m *manager) UpdatePod(podUID types.UID, newPod *apiObject.Pod) {
	m.podCache.Update(podUID, newPod)
}

func (m *manager) DeletePod(podUID types.UID) {
	m.podCache.Delete(podUID)
}

func NewStatusManager() Manager {
	return &manager{
		podCache: cache.Default(),
	}
}
