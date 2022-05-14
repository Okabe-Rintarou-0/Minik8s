package cache

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"minik8s/entity"
)

// updateNodeStatus updates the cache, to sync with api-server
// Incremental Synchronization
func (m *manager) updateNodeStatus(msg *redis.Message) {
	nodeStatus := &entity.NodeStatus{}
	err := json.Unmarshal([]byte(msg.Payload), nodeStatus)
	if err != nil {
		log(err.Error())
		return
	}
	log("Received status %s of Node[host = %s, cpu = %v, mem = %v, pods = %v]", nodeStatus.Lifecycle.String(), nodeStatus.Hostname, nodeStatus.CpuPercent, nodeStatus.MemPercent, nodeStatus.NumPods)
	hostName := nodeStatus.Hostname

	switch nodeStatus.Lifecycle {
	case entity.NodeCreated:
		m.nodeStatusCache.Add(hostName, nodeStatus)
	case entity.NodeUnknown:
		// ignore
		break
	case entity.NodeDeleted:
		m.nodeStatusCache.Delete(hostName)
	default:
		m.nodeStatusCache.Update(hostName, nodeStatus)
	}
	//m.nodeStatusUpdateHook(replicaSetStatus)
}
