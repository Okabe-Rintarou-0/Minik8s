package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"minik8s/apiObject"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/helper"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"minik8s/util/logger"
	"net/http"
	"path"
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

func getReplicaSetStatusFromEtcd(namespace, name string) (replicaSet *entity.ReplicaSetStatus) {
	etcdURL := path.Join(url.ReplicaSetURL, "status", namespace, name)
	if raw, err := etcd.Get(etcdURL); err == nil {
		if err = json.Unmarshal([]byte(raw), &replicaSet); err == nil {
			return replicaSet
		}
		logger.Error(err.Error())
	}
	return nil
}

func getReplicaSetStatusesFromEtcd() (replicaSets []*entity.ReplicaSetStatus) {
	etcdURL := path.Join(url.ReplicaSetURL, "status")
	raws, err := etcd.GetAll(etcdURL)
	for _, raw := range raws {
		replicaSet := &entity.ReplicaSetStatus{}
		if err = json.Unmarshal([]byte(raw), &replicaSet); err == nil {
			replicaSets = append(replicaSets, replicaSet)
		}
	}
	return
}

func getHPAStatusFromEtcd(namespace, name string) (hpa *entity.HPAStatus) {
	etcdURL := path.Join(url.HPAURL, "status", namespace, name)
	if raw, err := etcd.Get(etcdURL); err == nil {
		if err = json.Unmarshal([]byte(raw), &hpa); err == nil {
			return hpa
		}
		logger.Error(err.Error())
	}
	return nil
}

func getHPAStatusesFromEtcd() (hpas []*entity.HPAStatus) {
	etcdURL := path.Join(url.HPAURL, "status")
	raws, err := etcd.GetAll(etcdURL)
	for _, raw := range raws {
		hpa := &entity.HPAStatus{}
		if err = json.Unmarshal([]byte(raw), &hpa); err == nil {
			hpas = append(hpas, hpa)
		}
	}
	return
}

func getPodApiObjectFromEtcd(node, namespace, name string) (pod *apiObject.Pod) {
	etcdURL := path.Join(url.PodURL, node, namespace, name)
	if raw, err := etcd.Get(etcdURL); err == nil {
		if err = json.Unmarshal([]byte(raw), &pod); err == nil {
			return pod
		}
	}
	return nil
}

func getReplicaSetApiObjectFromEtcd(namespace, name string) (replicaSet *apiObject.ReplicaSet) {
	etcdURL := path.Join(url.ReplicaSetURL, namespace, name)
	if raw, err := etcd.Get(etcdURL); err == nil {
		if err = json.Unmarshal([]byte(raw), &replicaSet); err == nil {
			return replicaSet
		}
	}
	return nil
}

func getHPAApiObjectFromEtcd(namespace, name string) (hpa *apiObject.HorizontalPodAutoscaler) {
	etcdURL := path.Join(url.HPAURL, namespace, name)
	if raw, err := etcd.Get(etcdURL); err == nil {
		if err = json.Unmarshal([]byte(raw), &hpa); err == nil {
			return hpa
		}
	}
	return nil
}

func HandleGetNodeStatus(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	c.JSON(http.StatusOK, getNodeStatusFromEtcd(namespace, name))
}

func HandleGetNodeStatuses(c *gin.Context) {
	c.JSON(http.StatusOK, getNodeStatusesFromEtcd())
}

func HandleGetPodStatus(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	c.JSON(http.StatusOK, getPodStatusFromEtcd(namespace, name))
}

func HandleGetPodStatuses(c *gin.Context) {
	c.JSON(http.StatusOK, getPodStatusesFromEtcd())
}

func HandleGetReplicaSetStatus(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	c.JSON(http.StatusOK, getReplicaSetStatusFromEtcd(namespace, name))
}

func HandleGetReplicaSetStatuses(c *gin.Context) {
	c.JSON(http.StatusOK, getReplicaSetStatusesFromEtcd())
}

func HandleDescribePod(c *gin.Context) {
	name := c.Param("name")
	c.JSON(http.StatusOK, getPodDescriptionForTest(name))
}

func HandleGetPodApiObject(c *gin.Context) {
	node := c.Param("node")
	namespace := c.Param("namespace")
	name := c.Param("name")
	c.JSON(http.StatusOK, getPodApiObjectFromEtcd(node, namespace, name))
}

func HandleGetPodsApiObject(c *gin.Context) {
	node := c.Param("node")
	c.JSON(http.StatusOK, helper.GetPodsApiObjectFromEtcd(node))
}

func HandleGetReplicaSetApiObject(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	c.JSON(http.StatusOK, getReplicaSetApiObjectFromEtcd(namespace, name))
}

func HandleGetHPAApiObject(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	c.JSON(http.StatusOK, getHPAApiObjectFromEtcd(namespace, name))
}

func HandleGetHPAStatus(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	c.JSON(http.StatusOK, getHPAStatusFromEtcd(namespace, name))
}

func HandleGetHPAStatuses(c *gin.Context) {
	c.JSON(http.StatusOK, getHPAStatusesFromEtcd())
}
