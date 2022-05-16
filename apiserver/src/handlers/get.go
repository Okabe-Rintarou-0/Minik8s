package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"minik8s/apiObject"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"minik8s/util/logger"
	"net/http"
	"path"
	"time"
)

func getNodeStatusesFromEtcd() (nodes []*entity.NodeStatus) {
	etcdURL := path.Join(url.NodeURL, "status")
	raws, err := etcd.GetAll(etcdURL)
	for _, raw := range raws {
		node := &entity.NodeStatus{}
		if err = json.Unmarshal([]byte(raw), &node); err == nil {
			nodes = append(nodes, node)
		}
	}
	return
}

func getNodeStatusFromEtcd(namespace, name string) (node *entity.NodeStatus) {
	etcdURL := path.Join(url.NodeURL, "status", namespace, name)
	if raw, err := etcd.Get(etcdURL); err == nil {
		if err = json.Unmarshal([]byte(raw), &node); err == nil {
			return node
		}
		logger.Error(err.Error())
	}
	return nil
}

func getPodStatusFromEtcd(namespace, name string) (pod *entity.PodStatus) {
	etcdURL := path.Join(url.PodURL, "status", namespace, name)
	if raw, err := etcd.Get(etcdURL); err == nil {
		if err = json.Unmarshal([]byte(raw), &pod); err == nil {
			return pod
		}
		logger.Error(err.Error())
	}
	return nil
}

func getPodStatusesFromEtcd() (pods []*entity.PodStatus) {
	etcdURL := path.Join(url.PodURL, "status")
	raws, err := etcd.GetAll(etcdURL)
	for _, raw := range raws {
		pod := &entity.PodStatus{}
		if err = json.Unmarshal([]byte(raw), &pod); err == nil {
			pods = append(pods, pod)
		}
	}
	return
}

func getPodApiObjectFromEtcd(namespace, name string) (pod *apiObject.Pod) {
	etcdURL := path.Join(url.PodURL, namespace, name)
	if raw, err := etcd.Get(etcdURL); err == nil {
		if err = json.Unmarshal([]byte(raw), &pod); err == nil {
			return pod
		}
	}
	return nil
}

func HandleGetNodeStatuses(c *gin.Context) {
	c.JSON(http.StatusOK, getNodeStatusesFromEtcd())
}

func HandleGetNodeStatus(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	c.JSON(http.StatusOK, getNodeStatusFromEtcd(namespace, name))
}

func getPodDescriptionForTest(name string) *entity.PodDescription {
	var logs []entity.PodStatusLogEntry

	podStatus := podStatusForTest()
	logs = append(logs, entity.PodStatusLogEntry{
		Status: podStatus.Lifecycle,
		Time:   podStatus.SyncTime,
		Error:  podStatus.Error,
	})

	logs = append(logs, entity.PodStatusLogEntry{
		Status: entity.PodRunning,
		Time:   time.Now().Add(time.Minute * 30),
		Error:  "",
	})

	return &entity.PodDescription{
		CurrentStatus: *podStatusForTest(),
		Logs:          logs,
	}
}

func HandleGetPodStatuses(c *gin.Context) {
	c.JSON(http.StatusOK, getPodStatusesFromEtcd())
}

func HandleGetPodStatus(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	c.JSON(http.StatusOK, getPodStatusFromEtcd(namespace, name))
}

func HandleDescribePod(c *gin.Context) {
	name := c.Param("name")
	c.JSON(http.StatusOK, getPodDescriptionForTest(name))
}

func HandleGetPod(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	c.JSON(http.StatusOK, getPodApiObjectFromEtcd(namespace, name))
}
