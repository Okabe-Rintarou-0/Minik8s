package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"minik8s/apiObject"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/helper"
	"minik8s/apiserver/src/ipgen"
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
		Ip:         node.Ip,
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

	if helper.ExistsPod(pod.Namespace(), pod.Name()) {
		c.String(http.StatusOK, fmt.Sprintf("pod %s/%s already exists", pod.Namespace(), pod.Name()))
		return
	}

	pod.Metadata.UID = uidutil.New()
	if im, err := ipgen.New(url.PodIpGeneratorURL, url.PodIpBase); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	} else {
		if pod.Spec.ClusterIp, err = im.GetNext(); err != nil {
			c.String(http.StatusOK, err.Error())
			return
		}
	}
	log("receive pod %s/%s[ID = %v] %+v", pod.Namespace(), pod.Name(), pod.UID(), pod)
	log("receive pod %s/%s: %+v", pod.Namespace(), pod.Name(), pod)

	// Schedule first, then put the data to url: PodURL/node/namespace/name
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

	log("receive hpa[ID = %v]: %v", hpa.UID(), hpa)

	if err = addHPA(&hpa); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	c.String(http.StatusOK, "Apply successfully!")
}

func HandleApplyService(c *gin.Context) {
	service := apiObject.Service{}
	err := readAndUnmarshal(c.Request.Body, &service)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	if helper.ExistsService(service.Metadata.Namespace, service.Metadata.Name) {
		c.String(http.StatusOK, fmt.Sprintf("service %s/%s already exists", service.Metadata.Namespace, service.Metadata.Name))
		return
	}

	service.Metadata.UID = uidutil.New()
	if ig, err := ipgen.New(url.SvcIpGeneratorURL, url.ServiceIpBase); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	} else if service.Spec.ClusterIP, err = ig.GetNext(); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	log("receive service: %+v", service)

	serviceUpdate := entity.ServiceUpdate{
		Action: entity.CreateAction,
		Target: entity.ServiceTarget{
			Service:   service,
			Endpoints: make([]apiObject.Endpoint, 0),
		},
	}

	var serviceJson []byte
	if serviceJson, err = json.Marshal(service); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if err := etcd.Put(path.Join(url.ServiceURL, service.Metadata.Namespace, service.Metadata.Name), string(serviceJson)); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	for key, value := range service.Spec.Selector {
		if err := etcd.Put(path.Join(url.ServiceURL, key, value, service.Metadata.UID), string(serviceJson)); err != nil {
			c.String(http.StatusOK, err.Error())
			return
		}
		if endpoints, err := helper.GetEndpoints(key, value); err != nil {
			c.String(http.StatusOK, err.Error())
			return
		} else {
			serviceUpdate.Target.Endpoints = append(serviceUpdate.Target.Endpoints, endpoints...)
		}
	}

	serviceUpdateMsg, _ := json.Marshal(serviceUpdate)
	listwatch.Publish(topicutil.ServiceUpdateTopic(), serviceUpdateMsg)

	c.String(http.StatusOK, "Apply successfully!")
}

func HandleApplyDNS(c *gin.Context) {
	dns := apiObject.Dns{}
	err := readAndUnmarshal(c.Request.Body, &dns)
	if err != nil {
		c.String(http.StatusOK, err.Error())
	}
	log("receive dns: %+v", dns)
}

func HandleApplyGpuJob(c *gin.Context) {
	gpu := apiObject.GpuJob{}
	err := readAndUnmarshal(c.Request.Body, &gpu)
	if err != nil {
		c.String(http.StatusOK, err.Error())
	}
	gpu.Metadata.UID = uidutil.New()
	log("receive gpu job[ID = %v]: %v", gpu.UID(), gpu)

	GpuUpdateMsg, _ := json.Marshal(entity.GpuUpdate{
		Action: entity.CreateAction,
		Target: gpu,
	})

	listwatch.Publish(topicutil.GpuJobUpdateTopic(), GpuUpdateMsg)
	c.String(http.StatusOK, "Apply successfully!")
}
