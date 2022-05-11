package cache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"minik8s/entity"
)

// updateReplicaSetStatus updates the cache, to sync with api-server
// Incremental Synchronization
func (m *manager) updateReplicaSetStatus(msg *redis.Message) {
	replicaSetStatus := &entity.ReplicaSetStatus{}
	err := json.Unmarshal([]byte(msg.Payload), replicaSetStatus)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Received status %s of ReplicaSet[ID = %s]\n", replicaSetStatus.Lifecycle.String(), replicaSetStatus.ID)
	if replicaSetStatus.Lifecycle == entity.ReplicaSetDeleted {
		m.replicaSetStatusCache.Delete(replicaSetStatus.FullName())
	} else {
		m.replicaSetStatusCache.Update(replicaSetStatus.FullName(), replicaSetStatus)
	}
	//m.replicaSetStatusUpdateHook(replicaSetStatus)
}
