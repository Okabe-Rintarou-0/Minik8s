package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"minik8s/apiObject"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/util/logger"
	"minik8s/util/topicutil"
	"minik8s/util/uidutil"
	"net/http"
	"path"
	"time"
)

var log = logger.Log("Api-server")

func readAndUnmarshal(body io.ReadCloser, target interface{}) error {
	content, err := ioutil.ReadAll(body)
	defer body.Close()
	if err != nil {
		return err
	}
	if err = json.Unmarshal(content, target); err != nil {
		return err
	}
	return nil
}

func HandleApplyNode(c *gin.Context) {
	node := apiObject.Node{}
	err := readAndUnmarshal(c.Request.Body, &node)
	if err != nil {
		c.String(http.StatusOK, err.Error())
	}
	node.Metadata.UID = uidutil.New()
	log("receive node[ID = %v]: %v", node.UID(), node)

	var nodeJson []byte
	if nodeJson, err = json.Marshal(node); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	// exists?
	etcdNodeURL := path.Join(url.NodeURL, node.Namespace(), node.Name())
	if podJsonStr, err := etcd.Get(etcdNodeURL); err == nil {
		getNode := &apiObject.Node{}
		if err = json.Unmarshal([]byte(podJsonStr), getNode); err == nil {
			c.String(http.StatusOK, fmt.Sprintf("Node %s/%s already exists", getNode.Namespace(), getNode.Name()))
			return
		}
	}

	if err = etcd.Put(etcdNodeURL, string(nodeJson)); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	etcdNodeStatusURL := path.Join(url.NodeURL, "status", node.Namespace(), node.Name())
	var nodeStatusJson []byte
	if nodeStatusJson, err = json.Marshal(entity.NodeStatus{
		Hostname:   node.Name(),
		Ip:         "",
		Labels:     node.Labels(),
		Lifecycle:  entity.NodeUnknown,
		Error:      "",
		CpuPercent: 0,
		MemPercent: 0,
		NumPods:    0,
		SyncTime:   time.Now(),
	}); err == nil {
		_ = etcd.Put(etcdNodeStatusURL, string(nodeStatusJson))
	}

	c.String(http.StatusOK, "Apply successfully!")
}

func HandleApplyPod(c *gin.Context) {
	pod := apiObject.Pod{}
	err := readAndUnmarshal(c.Request.Body, &pod)
	if err != nil {
		c.String(http.StatusOK, err.Error())
	}
	pod.Metadata.UID = uidutil.New()
	log("receive pod[ID = %v]: %v", pod.UID(), pod)

	var podJson []byte
	if podJson, err = json.Marshal(pod); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	// exists?
	etcdPodURL := path.Join(url.PodURL, pod.Namespace(), pod.Name())
	if podJsonStr, err := etcd.Get(etcdPodURL); err == nil {
		getPod := &apiObject.Pod{}
		if err = json.Unmarshal([]byte(podJsonStr), getPod); err == nil {
			c.String(http.StatusOK, fmt.Sprintf("pod %s/%s already exists", getPod.Namespace(), getPod.Name()))
			return
		}
	}

	if err = etcd.Put(etcdPodURL, string(podJson)); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	etcdPodStatusURL := path.Join(url.PodURL, "status", pod.Namespace(), pod.Name())
	var podStatusJson []byte
	if podStatusJson, err = json.Marshal(entity.PodStatus{
		ID:         pod.UID(),
		Name:       pod.Name(),
		Namespace:  pod.Namespace(),
		Labels:     pod.Labels(),
		Lifecycle:  entity.PodUnknown,
		CpuPercent: 0,
		MemPercent: 0,
		Error:      "",
		SyncTime:   time.Now(),
	}); err == nil {
		_ = etcd.Put(etcdPodStatusURL, string(podStatusJson))
	}

	podUpdateMsg, _ := json.Marshal(entity.PodUpdate{
		Action: entity.CreateAction,
		Target: pod,
	})
	listwatch.Publish(topicutil.SchedulerPodUpdateTopic(), podUpdateMsg)
	c.String(http.StatusOK, "Apply successfully!")
}

func HandleApplyReplicaSet(c *gin.Context) {
	rs := apiObject.ReplicaSet{}
	err := readAndUnmarshal(c.Request.Body, &rs)
	if err != nil {
		c.String(http.StatusOK, err.Error())
	}
	rs.Metadata.UID = uidutil.New()
	log("receive rs[ID = %v]: %v", rs.UID(), rs)

	var rsJson []byte
	if rsJson, err = json.Marshal(rs); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	// exists?
	etcdReplicaSetURL := path.Join(url.ReplicaSetURL, rs.Namespace(), rs.Name())
	if rsJsonStr, err := etcd.Get(etcdReplicaSetURL); err == nil {
		getRs := &apiObject.ReplicaSet{}
		if err = json.Unmarshal([]byte(rsJsonStr), getRs); err == nil {
			c.String(http.StatusOK, fmt.Sprintf("replicaSet %s/%s already exists", getRs.Namespace(), getRs.Name()))
			return
		}
	}

	if err = etcd.Put(etcdReplicaSetURL, string(rsJson)); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	etcdReplicaSetStatusURL := path.Join(url.ReplicaSetURL, "status", rs.Namespace(), rs.Name())
	var replicaSetStatusJson []byte
	if replicaSetStatusJson, err = json.Marshal(entity.ReplicaSetStatus{
		ID:         rs.UID(),
		Name:       rs.Name(),
		Namespace:  rs.Namespace(),
		Labels:     rs.Labels(),
		Lifecycle:  entity.ReplicaSetUnknown,
		CpuPercent: 0,
		MemPercent: 0,
		Error:      "",
		SyncTime:   time.Now(),
	}); err == nil {
		_ = etcd.Put(etcdReplicaSetStatusURL, string(replicaSetStatusJson))
	}

	replicaSetUpdateMsg, _ := json.Marshal(entity.ReplicaSetUpdate{
		Action: entity.CreateAction,
		Target: rs,
	})

	listwatch.Publish(topicutil.ReplicaSetUpdateTopic(), replicaSetUpdateMsg)
	c.String(http.StatusOK, "Apply successfully!")
}

func HandleApplyHPA(c *gin.Context) {
	hpa := apiObject.HorizontalPodAutoscaler{}
	err := readAndUnmarshal(c.Request.Body, &hpa)
	if err != nil {
		c.String(http.StatusOK, err.Error())
	}
	hpa.Metadata.UID = uidutil.New()
	log("receive hpa[ID = %v]: %v", hpa.UID(), hpa)

	if err = addHPA(&hpa); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	c.String(http.StatusOK, "Apply successfully!")
}
