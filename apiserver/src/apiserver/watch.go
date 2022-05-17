package apiserver

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"minik8s/util/logger"
	"path"
)

var logApiServer = logger.Log("api-server watching")

func syncNodeStatus(msg *redis.Message) {
	nodeStatus := &entity.NodeStatus{}
	err := json.Unmarshal([]byte(msg.Payload), nodeStatus)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logApiServer("Received status %s of Node[host = %s, cpu = %v, mem = %v, pods = %v]", nodeStatus.Lifecycle.String(), nodeStatus.Hostname, nodeStatus.CpuPercent, nodeStatus.MemPercent, nodeStatus.NumPods)

	var nodeStatusJson []byte
	if nodeStatusJson, err = json.Marshal(nodeStatus); err == nil {
		etcdURL := path.Join(url.NodeURL, "status", nodeStatus.Namespace, nodeStatus.Hostname)
		err = etcd.Put(etcdURL, string(nodeStatusJson))
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

func syncPodStatus(msg *redis.Message) {
	podStatus := &entity.PodStatus{}
	err := json.Unmarshal([]byte(msg.Payload), podStatus)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if podStatus.Lifecycle == entity.PodDeleted {
		return
	}

	logApiServer("Received status %s of Pod[name = %s, id = %v]", podStatus.Lifecycle.String(), podStatus.Name, podStatus.ID)

	var podStatusJson []byte
	if podStatusJson, err = json.Marshal(podStatus); err == nil {
		etcdURL := path.Join(url.PodURL, "status", podStatus.Namespace, podStatus.Name)
		err = etcd.Put(etcdURL, string(podStatusJson))
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

func syncReplicaSetStatus(msg *redis.Message) {
	replicaSetStatus := &entity.ReplicaSetStatus{}
	err := json.Unmarshal([]byte(msg.Payload), replicaSetStatus)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if replicaSetStatus.Lifecycle == entity.ReplicaSetDeleted {
		return
	}

	logApiServer("Received status %s of Rs[name = %s, id = %v]", replicaSetStatus.Lifecycle.String(), replicaSetStatus.Name, replicaSetStatus.ID)

	var replicaSetStatusJson []byte
	if replicaSetStatusJson, err = json.Marshal(replicaSetStatus); err == nil {
		etcdURL := path.Join(url.ReplicaSetURL, "status", replicaSetStatus.Namespace, replicaSetStatus.Name)
		err = etcd.Put(etcdURL, string(replicaSetStatusJson))
		if err != nil {
			logger.Error(err.Error())
		}
	}
}
