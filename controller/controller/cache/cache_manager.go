package cache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"minik8s/apiObject/types"
	"minik8s/entity"
	"minik8s/kubelet/src/runtime/runtime"
	"minik8s/listwatch"
	"minik8s/util"
	utilcache "minik8s/util/cache"
)

type PodStatusUpdateHook func(podStatus *entity.PodStatus)

type Manager interface {
	Start()
	GetReplicaSetPodStatuses(rsUID types.UID) []*entity.PodStatus
	SetPodStatusUpdateHook(podStatusUpdateHook PodStatusUpdateHook)
}

type manager struct {
	podStatusCache      utilcache.Cache
	podStatusUpdateHook PodStatusUpdateHook
}

func (m *manager) SetPodStatusUpdateHook(podStatusUpdateHook PodStatusUpdateHook) {
	m.podStatusUpdateHook = podStatusUpdateHook
}

// updatePodStatus updates the cache, to sync with api-server
// Incremental Synchronization
func (m *manager) updatePodStatus(msg *redis.Message) {
	podStatus := &entity.PodStatus{}
	err := json.Unmarshal([]byte(msg.Payload), podStatus)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Received status %s of Pod[ID = %s]\n", podStatus.Status.String(), podStatus.ID)
	if podStatus.Status == entity.Deleted {
		m.podStatusCache.Delete(podStatus.ID)
	} else {
		m.podStatusCache.Update(podStatus.ID, podStatus)
	}
	m.podStatusUpdateHook(podStatus)
}

// Start starts the manager and listens to the pod status topic
// The pod status message is sent when the pod worker has created, deleted a pod
// After receiving the msg, updatePodStatus will be called to update the cache.
func (m *manager) Start() {
	topic := util.PodStatusTopic()
	listwatch.Watch(topic, m.updatePodStatus)
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
		podStatusCache: utilcache.Default(),
	}
}
