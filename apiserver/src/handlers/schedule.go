package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"minik8s/apiObject"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/helper"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"net/http"
	"path"
	"time"
)

func HandleSchedulePod(c *gin.Context) {
	node := c.Param("node")

	pod := apiObject.Pod{}
	err := readAndUnmarshal(c.Request.Body, &pod)
	if err != nil {
		c.String(http.StatusOK, err.Error())
	}

	if helper.ExistsPod(pod.Namespace(), pod.Name()) {
		c.String(http.StatusOK, fmt.Sprintf("pod %s/%s already exists", pod.Namespace(), pod.Name()))
		return
	}

	var podJson []byte
	if podJson, err = json.Marshal(pod); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	etcdPodURL := path.Join(url.PodURL, node, pod.Namespace(), pod.Name())
	if err := etcd.Put(etcdPodURL, string(podJson)); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	// Store pod's endpoints into etcd
	// @TODO push to proxy
	if err = helper.AddEndpoints(pod); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	etcdPodStatusURL := path.Join(url.PodURL, "status", pod.Namespace(), pod.Name())
	var podStatusJson []byte
	if podStatusJson, err = json.Marshal(entity.PodStatus{
		ID:         pod.UID(),
		Node:       node,
		Name:       pod.Name(),
		Namespace:  pod.Namespace(),
		Labels:     pod.Labels(),
		Lifecycle:  entity.PodScheduled,
		CpuPercent: 0,
		MemPercent: 0,
		Error:      "",
		SyncTime:   time.Now(),
	}); err == nil {
		_ = etcd.Put(etcdPodStatusURL, string(podStatusJson))
	}

	log("Schedule pod %s/%s to node %s", pod.Namespace(), pod.Name(), node)

	c.String(http.StatusOK, "Schedule successfully!")
}
