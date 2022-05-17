package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"minik8s/apiObject"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/util/topicutil"
	"net/http"
	"path"
	"strconv"
)

func HandleSetNodeStatus(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	form := getPutForm(c.Request.Body)
	if form == nil {
		c.String(http.StatusOK, "must have a form!")
		return
	}
	lifecycleInt64, _ := strconv.Atoi(form["lifecycle"])
	lifecycle := entity.NodeLifecycle(lifecycleInt64)
	log("Received lifecycle %v from %v", lifecycle.String(), name)

	etcdURL := path.Join(url.NodeURL, "status", namespace, name)
	raw, err := etcd.Get(etcdURL)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	nodeStatus := &entity.NodeStatus{}
	if err = json.Unmarshal([]byte(raw), nodeStatus); err != nil {
		c.String(http.StatusOK, fmt.Sprintf("no such node %s/%s", namespace, name))
		return
	}

	nodeStatus.Lifecycle = lifecycle
	nodeStatusJson, _ := json.Marshal(nodeStatus)

	if err = etcd.Put(etcdURL, string(nodeStatusJson)); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("Set node %s/%s status successfully", namespace, name))
}

func HandleSetReplicaSet(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	form := getPutForm(c.Request.Body)
	if form == nil {
		c.String(http.StatusOK, "must have a form!")
		return
	}

	replicas, _ := strconv.Atoi(form["replicas"])
	log("Received replicas %v from %s/%s", replicas, namespace, name)

	etcdURL := path.Join(url.ReplicaSetURL, namespace, name)
	raw, err := etcd.Get(etcdURL)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	replicaSet := &apiObject.ReplicaSet{}
	if err = json.Unmarshal([]byte(raw), replicaSet); err != nil {
		c.String(http.StatusOK, fmt.Sprintf("no such rs %s/%s", namespace, name))
		return
	}

	replicaSet.SetReplicas(replicas)
	replicaSetJson, _ := json.Marshal(replicaSet)

	if err = etcd.Put(etcdURL, string(replicaSetJson)); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	replicaSetUpdateMsg, _ := json.Marshal(entity.ReplicaSetUpdate{
		Action: entity.UpdateAction,
		Target: *replicaSet,
	})
	listwatch.Publish(topicutil.ReplicaSetUpdateTopic(), replicaSetUpdateMsg)
	c.String(http.StatusOK, fmt.Sprintf("Set replicaSet %s/%s num replicas to %v successfully", namespace, name, replicas))
}
