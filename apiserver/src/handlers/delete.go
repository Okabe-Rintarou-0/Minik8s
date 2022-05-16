package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"minik8s/apiObject"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/util/topicutil"
	"net/http"
	"path"
)

func deleteSpecifiedNode(namespace, name string) (err error) {
	log("Node to delete is %s/%s", namespace, name)
	etcdNodeURL := path.Join(url.NodeURL, namespace, name)
	if err = etcd.Delete(etcdNodeURL); err == nil {
		etcdNodeStatusURL := path.Join(url.NodeURL, "status", namespace, name)
		err = etcd.Delete(etcdNodeStatusURL)
	}
	return
}

func deleteSpecifiedPod(namespace, name string) (pod *apiObject.Pod, err error) {
	log("Pod to delete is %s/%s", namespace, name)
	etcdPodURL := path.Join(url.PodURL, namespace, name)

	etcdPodStatusURL := path.Join(url.PodURL, "status", namespace, name)
	_ = etcd.Delete(etcdPodStatusURL)

	var raw string
	if raw, err = etcd.Get(etcdPodURL); err != nil {
		return nil, err
	}

	if err = json.Unmarshal([]byte(raw), &pod); err != nil {
		return nil, err
	}

	err = etcd.Delete(etcdPodURL)
	return
}

func HandleDeleteNode(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	if err := deleteSpecifiedNode(namespace, name); err != nil {
		c.String(http.StatusOK, err.Error())
	}
	c.String(http.StatusOK, "Delete successfully")
}

func HandleDeletePod(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	if podToDelete, err := deleteSpecifiedPod(namespace, name); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	} else {
		podDeleteMsg, _ := json.Marshal(entity.PodUpdate{
			Action: entity.DeleteAction,
			Target: *podToDelete,
		})
		listwatch.Publish(topicutil.SchedulerPodUpdateTopic(), podDeleteMsg)
	}
	c.String(http.StatusOK, "Delete successfully")
}
