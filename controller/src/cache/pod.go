package cache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"minik8s/entity"
)

// updatePodStatus updates the cache, to sync with api-server
// Incremental Synchronization
func (m *manager) updatePodStatus(msg *redis.Message) {
	podStatus := &entity.PodStatus{}
	err := json.Unmarshal([]byte(msg.Payload), podStatus)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	log("Received status %s of Pod[ID = %s]", podStatus.Lifecycle.String(), podStatus.ID)
	if podStatus.Lifecycle == entity.PodDeleted {
		m.podStatusCache.Delete(podStatus.ID)
	} else {
		m.podStatusCache.Update(podStatus.ID, podStatus)
	}
	m.podStatusUpdateHook(podStatus)
}
