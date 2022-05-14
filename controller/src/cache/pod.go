package cache

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"minik8s/entity"
)

// updatePodStatus updates the cache, to sync with api-server
// Incremental Synchronization
func (m *manager) updatePodStatus(msg *redis.Message) {
	podStatus := &entity.PodStatus{}
	err := json.Unmarshal([]byte(msg.Payload), podStatus)
	if err != nil {
		log(err.Error())
		return
	}
	log("Received status %s of Pod[ID = %s, cpu = %v, mem = %v]", podStatus.Lifecycle.String(), podStatus.ID, podStatus.CpuPercent, podStatus.MemPercent)
	switch podStatus.Lifecycle {
	case entity.PodDeleted:
		m.podStatusCache.Delete(podStatus.ID)
	default:
		if !m.podStatusCache.Exists(podStatus.ID) {
			m.podStatusCache.Add(podStatus.ID, podStatus)
		} else {
			m.podStatusCache.Update(podStatus.ID, podStatus)
		}
	}
	m.podStatusUpdateHook(podStatus)
}
