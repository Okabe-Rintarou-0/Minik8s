package cache

import (
	"minik8s/apiObject/types"
	"minik8s/entity"
	"minik8s/kubelet/src/runtime/runtime"
	"minik8s/listwatch"
	cacheutil "minik8s/util/cache"
	"minik8s/util/logger"
	"minik8s/util/topicutil"
	"minik8s/util/wait"
	"sync"
	"time"
)

var log = logger.Log("Cache Manager")

type PodStatusUpdateHook func(podStatus *entity.PodStatus)

type ReplicaSetStatusUpdateHook func(replicaSetStatus *entity.ReplicaSetStatus)

type Manager interface {
	Start()
	RefreshReplicaSetStatus(fullName string)
	GetReplicaSetStatus(fullName string) *entity.ReplicaSetStatus
	GetReplicaSetPodStatuses(rsUID types.UID) []*entity.PodStatus
	GetNodeStatuses() []*entity.NodeStatus
	SetNodeStatus(hostname string, newStatus *entity.NodeStatus)
	SetPodStatusUpdateHook(podStatusUpdateHook PodStatusUpdateHook)
	SetReplicaSetStatusUpdateHook(replicaSetStatusUpdateHook ReplicaSetStatusUpdateHook)
}

type manager struct {
	nodeStatusLock             sync.RWMutex
	podStatusCache             cacheutil.Cache
	replicaSetStatusCache      cacheutil.Cache
	nodeStatusCache            cacheutil.Cache
	podStatusUpdateHook        PodStatusUpdateHook
	replicaSetStatusUpdateHook ReplicaSetStatusUpdateHook
}

func (m *manager) SetNodeStatus(hostname string, newStatus *entity.NodeStatus) {
	m.nodeStatusLock.Lock()
	defer m.nodeStatusLock.Unlock()
	m.nodeStatusCache.Update(hostname, newStatus)
}

func (m *manager) getNodeStatusesInternal() []*entity.NodeStatus {
	values := m.nodeStatusCache.Values()
	var nodeStatus *entity.NodeStatus
	var nodeStatuses []*entity.NodeStatus
	for _, value := range values {
		nodeStatus = value.(*entity.NodeStatus)
		nodeStatuses = append(nodeStatuses, nodeStatus)
	}
	return nodeStatuses
}

func (m *manager) GetNodeStatuses() []*entity.NodeStatus {
	m.nodeStatusLock.RLock()
	defer m.nodeStatusLock.RUnlock()
	return m.getNodeStatusesInternal()
}

func (m *manager) calcMetrics(podStatuses []*entity.PodStatus) (cpuPercent, memPercent float64) {
	cpuPercent = 0.0
	memPercent = 0.0
	for _, podStatus := range podStatuses {
		cpuPercent += podStatus.CpuPercent
		memPercent += podStatus.MemPercent
	}
	return
}

// RefreshReplicaSetStatus TODO deprecate it
func (m *manager) RefreshReplicaSetStatus(fullName string) {
	if v := m.replicaSetStatusCache.Get(fullName); v != nil {
		replicaSetStatus := v.(*entity.ReplicaSetStatus)
		podStatuses := m.GetReplicaSetPodStatuses(replicaSetStatus.ID)
		replicaSetStatus.CpuPercent, replicaSetStatus.MemPercent = m.calcMetrics(podStatuses)
	}
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
	go listwatch.Watch(topicutil.NodeStatusTopic(), m.updateNodeStatus)

	go wait.Period(time.Second*30, nodeStatusFullSyncPeriod, m.fullSyncNodeStatuses)
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
		podStatusCache:        cacheutil.Default(),
		replicaSetStatusCache: cacheutil.Default(),
		nodeStatusCache:       cacheutil.Default(),
	}
}
