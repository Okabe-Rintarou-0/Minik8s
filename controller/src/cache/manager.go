package cache

import (
	"minik8s/apiObject/types"
	"minik8s/entity"
	"minik8s/kubelet/src/runtime/runtime"
	"minik8s/listwatch"
	utilcache "minik8s/util/cache"
	"minik8s/util/logger"
	"minik8s/util/topicutil"
)

var log = logger.Log("Cache Manager")

type PodStatusUpdateHook func(podStatus *entity.PodStatus)

type ReplicaSetStatusUpdateHook func(replicaSetStatus *entity.ReplicaSetStatus)

type Manager interface {
	Start()
	GetReplicaSetStatus(fullName string) *entity.ReplicaSetStatus
	GetReplicaSetPodStatuses(rsUID types.UID) []*entity.PodStatus
	SetPodStatusUpdateHook(podStatusUpdateHook PodStatusUpdateHook)
	SetReplicaSetStatusUpdateHook(replicaSetStatusUpdateHook ReplicaSetStatusUpdateHook)
}

type manager struct {
	podStatusCache             utilcache.Cache
	replicaSetStatusCache      utilcache.Cache
	podStatusUpdateHook        PodStatusUpdateHook
	replicaSetStatusUpdateHook ReplicaSetStatusUpdateHook
}

func (m *manager) GetReplicaSetStatus(fullName string) *entity.ReplicaSetStatus {
	if v := m.replicaSetStatusCache.Get(fullName); v != nil {
		return v.(*entity.ReplicaSetStatus)
	}
	return nil
}

func (m *manager) SetPodStatusUpdateHook(podStatusUpdateHook PodStatusUpdateHook) {
	m.podStatusUpdateHook = podStatusUpdateHook
}

func (m *manager) SetReplicaSetStatusUpdateHook(replicaSetStatusUpdateHook ReplicaSetStatusUpdateHook) {
	m.replicaSetStatusUpdateHook = replicaSetStatusUpdateHook
}

// Start starts the manager and listens to the pod status topic
// The pod status message is sent when the pod worker has created, deleted a pod
// After receiving the msg, updatePodStatus will be called to update the cache.
func (m *manager) Start() {
	go listwatch.Watch(topicutil.PodStatusTopic(), m.updatePodStatus)
	go listwatch.Watch(topicutil.ReplicaSetStatusTopic(), m.updateReplicaSetStatus)
}

func (m *manager) GetReplicaSetPodStatuses(rsUID types.UID) []*entity.PodStatus {
	values := m.podStatusCache.Values()
	var podStatus *entity.PodStatus
	var podStatuses []*entity.PodStatus
	for _, value := range values {
		podStatus = value.(*entity.PodStatus)
		if podRsUID, exists := podStatus.Labels[runtime.KubernetesReplicaSetUIDLabel]; exists && podRsUID == rsUID {
			podStatuses = append(podStatuses, podStatus)
		}
	}
	return podStatuses
}

func NewManager() Manager {
	return &manager{
		podStatusCache:        utilcache.Default(),
		replicaSetStatusCache: utilcache.Default(),
	}
}
