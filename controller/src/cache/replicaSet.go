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
	UID := replicaSetStatus.ID
	log("Received status %s of ReplicaSet[ID = %s]", replicaSetStatus.Lifecycle.String(), UID)
	if replicaSetStatus.Lifecycle == entity.ReplicaSetDeleted {
		m.replicaSetStatusCache.Delete(UID)
	} else {
		if m.replicaSetStatusCache.Exists(UID) {
			m.replicaSetStatusCache.Update(UID, replicaSetStatus)
		} else {
			m.replicaSetStatusCache.Add(UID, replicaSetStatus)
		}
	}
	//m.replicaSetStatusUpdateHook(replicaSetStatus)
}
