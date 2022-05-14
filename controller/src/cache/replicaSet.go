package cache

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"minik8s/entity"
)

// updateReplicaSetStatus updates the cache, to sync with api-server
// Incremental Synchronization
func (m *manager) updateReplicaSetStatus(msg *redis.Message) {
	replicaSetStatus := &entity.ReplicaSetStatus{}
	err := json.Unmarshal([]byte(msg.Payload), replicaSetStatus)
	if err != nil {
		log(err.Error())
		return
	}
	log("Received status %s of ReplicaSet[ID = %s]", replicaSetStatus.Lifecycle.String(), replicaSetStatus.ID)
	if replicaSetStatus.Lifecycle == entity.ReplicaSetDeleted {
		m.replicaSetStatusCache.Delete(replicaSetStatus.FullName())
	} else {
		fullName := replicaSetStatus.FullName()
		if m.replicaSetStatusCache.Exists(fullName) {
			m.replicaSetStatusCache.Update(fullName, replicaSetStatus)
		} else {
			m.replicaSetStatusCache.Add(fullName, replicaSetStatus)
		}
	}
	//m.replicaSetStatusUpdateHook(replicaSetStatus)
}
