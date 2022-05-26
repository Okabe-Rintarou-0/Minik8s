package apiserver

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/helper"
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

	etcdURL := path.Join(url.NodeURL, "status", nodeStatus.Namespace, nodeStatus.Hostname)
	var oldNodeStatusStr string
	if oldNodeStatusStr, err = etcd.Get(etcdURL); err == nil {
		oldNodeStatus := entity.NodeStatus{}
		if err = json.Unmarshal([]byte(oldNodeStatusStr), &oldNodeStatus); err == nil {
			nodeStatus.Ip = oldNodeStatus.Ip
			nodeStatus.Labels = oldNodeStatus.Labels
		}
	}

	var nodeStatusJson []byte
	if nodeStatusJson, err = json.Marshal(nodeStatus); err == nil {
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

	if podStatus.Lifecycle == entity.PodDeleted || !helper.ExistsPod(podStatus.Namespace, podStatus.Name) {
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

func syncHPAStatus(msg *redis.Message) {
	hpaStatus := &entity.HPAStatus{}
	err := json.Unmarshal([]byte(msg.Payload), hpaStatus)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if hpaStatus.Lifecycle == entity.HPADeleted {
		return
	}

	logApiServer("Received status %s of HPA[name = %s, id = %v]", hpaStatus.Lifecycle.String(), hpaStatus.Name, hpaStatus.ID)

	var hpaStatusJson []byte
	if hpaStatusJson, err = json.Marshal(hpaStatus); err == nil {
		etcdURL := path.Join(url.HPAURL, "status", hpaStatus.Namespace, hpaStatus.Name)
		err = etcd.Put(etcdURL, string(hpaStatusJson))
		if err != nil {
			logger.Error(err.Error())
		}
	}
}
