package cache

import (
	"minik8s/apiObject"
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

type ReplicaSetFullSyncAddHook func(replicaSetStatus *apiObject.ReplicaSet)
type ReplicaSetFullSyncDeleteHook func(replicaSetStatus *apiObject.ReplicaSet)

type HPAFullSyncAddHook func(hpa *apiObject.HorizontalPodAutoscaler)
type HPAFullSyncDeleteHook func(hpa *apiObject.HorizontalPodAutoscaler)

type Manager interface {
	Start()
	GetReplicaSetStatus(UID types.UID) *entity.ReplicaSetStatus
	GetReplicaSetPodStatuses(UID types.UID) []*entity.PodStatus
	GetNodeStatuses() []*entity.NodeStatus
	SetNodeStatus(fullName string, newStatus *entity.NodeStatus)
	DeleteNodeStatus(fullName string)
	SetPodStatusUpdateHook(podStatusUpdateHook PodStatusUpdateHook)
	SetReplicaSetFullSyncAddHook(replicaSetFullSyncAddHook ReplicaSetFullSyncAddHook)
	SetReplicaSetFullSyncDeleteHook(replicaSetFullSyncDeleteHook ReplicaSetFullSyncDeleteHook)
	SetHPAFullSyncAddHook(hpaFullSyncAddHook HPAFullSyncAddHook)
	SetHPAFullSyncDeleteHook(hpaFullSyncDeleteHook HPAFullSyncDeleteHook)
}

type manager struct {
	nodeStatusLock               sync.RWMutex
	podStatusCache               cacheutil.Cache
	replicaSetStatusCache        cacheutil.Cache
	nodeStatusCache              cacheutil.Cache
	hpaStatusCache               cacheutil.Cache
	podStatusUpdateHook          PodStatusUpdateHook
	replicaSetFullSyncAddHook    ReplicaSetFullSyncAddHook
	replicaSetFullSyncDeleteHook ReplicaSetFullSyncDeleteHook
	hpaFullSyncAddHook           HPAFullSyncAddHook
	hpaFullSyncDeleteHook        HPAFullSyncDeleteHook
}

func (m *manager) SetHPAFullSyncAddHook(hpaFullSyncAddHook HPAFullSyncAddHook) {
	m.hpaFullSyncAddHook = hpaFullSyncAddHook
}

func (m *manager) SetHPAFullSyncDeleteHook(hpaFullSyncDeleteHook HPAFullSyncDeleteHook) {
	m.hpaFullSyncDeleteHook = hpaFullSyncDeleteHook
}

func (m *manager) SetReplicaSetFullSyncAddHook(replicaSetFullSyncAddHook ReplicaSetFullSyncAddHook) {
	m.replicaSetFullSyncAddHook = replicaSetFullSyncAddHook
}

func (m *manager) SetReplicaSetFullSyncDeleteHook(replicaSetFullSyncDeleteHook ReplicaSetFullSyncDeleteHook) {
	m.replicaSetFullSyncDeleteHook = replicaSetFullSyncDeleteHook
}

func (m *manager) DeleteNodeStatus(fullName string) {
	m.nodeStatusLock.Lock()
	defer m.nodeStatusLock.Unlock()
	m.nodeStatusCache.Delete(fullName)
}

func (m *manager) SetNodeStatus(fullName string, newStatus *entity.NodeStatus) {
	m.nodeStatusLock.Lock()
	defer m.nodeStatusLock.Unlock()
	m.nodeStatusCache.Update(fullName, newStatus)
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

func (m *manager) getPodStatusesInternal() []*entity.PodStatus {
	values := m.podStatusCache.Values()
	var podStatus *entity.PodStatus
	var podStatuses []*entity.PodStatus
	for _, value := range values {
		podStatus = value.(*entity.PodStatus)
		podStatuses = append(podStatuses, podStatus)
	}
	return podStatuses
}

func (m *manager) getReplicaSetStatusesInternal() []*entity.ReplicaSetStatus {
	values := m.replicaSetStatusCache.Values()
	var replicaSetStatus *entity.ReplicaSetStatus
	var replicaSetStatuses []*entity.ReplicaSetStatus
	for _, value := range values {
		replicaSetStatus = value.(*entity.ReplicaSetStatus)
		replicaSetStatuses = append(replicaSetStatuses, replicaSetStatus)
	}
	return replicaSetStatuses
}

func (m *manager) getHPAStatusesInternal() []*entity.HPAStatus {
	values := m.hpaStatusCache.Values()
	var hpaStatus *entity.HPAStatus
	var hpaStatuses []*entity.HPAStatus
	for _, value := range values {
		hpaStatus = value.(*entity.HPAStatus)
		hpaStatuses = append(hpaStatuses, hpaStatus)
	}
	return hpaStatuses
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

func (m *manager) GetReplicaSetStatus(UID types.UID) *entity.ReplicaSetStatus {
	if v := m.replicaSetStatusCache.Get(UID); v != nil {
		return v.(*entity.ReplicaSetStatus)
	}
	return nil
}

func (m *manager) SetPodStatusUpdateHook(podStatusUpdateHook PodStatusUpdateHook) {
	m.podStatusUpdateHook = podStatusUpdateHook
}

// Start starts the manager and listens to the pod status topic
// The pod status message is sent when the pod worker has created, deleted a pod
// After receiving the msg, updatePodStatus will be called to update the cache.
func (m *manager) Start() {
	go listwatch.Watch(topicutil.PodStatusTopic(), m.updatePodStatus)
	go listwatch.Watch(topicutil.ReplicaSetStatusTopic(), m.updateReplicaSetStatus)
	go listwatch.Watch(topicutil.NodeStatusTopic(), m.updateNodeStatus)

	go wait.Period(time.Second*30, nodeStatusFullSyncPeriod, m.fullSyncNodeStatuses)
	go wait.Period(time.Second*30, podStatusFullSyncPeriod, m.fullSyncPodStatuses)
	go wait.Period(time.Second*30, replicaSetStatusFullSyncPeriod, m.fullSyncReplicaSetStatuses)
	go wait.Period(time.Second*30, hpaStatusFullSyncPeriod, m.fullSyncHPAStatuses)
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
		hpaStatusCache:        cacheutil.Default(),
	}
}
