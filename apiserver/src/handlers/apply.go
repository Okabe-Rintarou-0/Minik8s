package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"minik8s/apiObject"
	dns2 "minik8s/apiserver/src/dns"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/helper"
	"minik8s/apiserver/src/url"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/nginx"
	"minik8s/util/logger"
	"minik8s/util/topicutil"
	"minik8s/util/uidutil"
	"minik8s/util/weaveutil"
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
		return
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

	c.String(http.StatusOK, "ok")
}

func HandleApplyPod(c *gin.Context) {
	pod := apiObject.Pod{}
	err := readAndUnmarshal(c.Request.Body, &pod)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	if helper.ExistsPod(pod.Namespace(), pod.Name()) {
		c.String(http.StatusOK, fmt.Sprintf("pod %s/%s already exists", pod.Namespace(), pod.Name()))
		return
	}

	pod.Metadata.UID = uidutil.New()
	if pod.Spec.ClusterIp, err = helper.NewPodIp(); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	log("receive pod %s/%s[ID = %v] %+v", pod.Namespace(), pod.Name(), pod.UID(), pod)

	// Schedule first, then put the data to url: PodURL/node/namespace/name
	podUpdateMsg, _ := json.Marshal(entity.PodUpdate{
		Action: entity.CreateAction,
		Target: pod,
	})

	listwatch.Publish(topicutil.SchedulerPodUpdateTopic(), podUpdateMsg)
	c.String(http.StatusOK, "ok")
}

func HandleApplyReplicaSet(c *gin.Context) {
	rs := apiObject.ReplicaSet{}
	err := readAndUnmarshal(c.Request.Body, &rs)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
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
	c.String(http.StatusOK, "ok")
}

func HandleApplyHPA(c *gin.Context) {
	hpa := apiObject.HorizontalPodAutoscaler{}
	err := readAndUnmarshal(c.Request.Body, &hpa)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	log("receive hpa[ID = %v]: %v", hpa.UID(), hpa)

	if err = addHPA(&hpa); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	c.String(http.StatusOK, "ok")
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
	if service.Spec.ClusterIP, err = helper.NewServiceIp(); err != nil {
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

	c.String(http.StatusOK, "ok")
}

func HandleApplyDNS(c *gin.Context) {
	dns := apiObject.Dns{}
	err := readAndUnmarshal(c.Request.Body, &dns)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if helper.ExistsDNS(dns.Metadata.Namespace, dns.Metadata.Name) {
		c.String(http.StatusOK, fmt.Sprintf("dns %s/%s already exists", dns.Metadata.Namespace, dns.Metadata.Name))
		return
	}

	dns.Metadata.UID = uidutil.New()
	log("receive dns: %+v", dns)

	nm := nginx.New(dns.Metadata.UID)

	// Step 1: Apply mappings to nginx.conf
	servers := make([]nginx.Server, 1)
	servers[0].Port = 80
	for _, p := range dns.Spec.Paths {
		log("%#v", p)
		if !helper.ExistsService(dns.Metadata.Namespace, p.Service.Name) {
			c.String(http.StatusOK, fmt.Sprintf("service %s/%s does not exist", dns.Metadata.Namespace, p.Service.Name))
			return
		}
		service := apiObject.Service{}
		if serviceJsonStr, err := etcd.Get(path.Join(url.ServiceURL, dns.Metadata.Namespace, p.Service.Name)); err != nil {
			c.String(http.StatusOK, err.Error())
			return
		} else {
			log("%#v", serviceJsonStr)
			if err := json.Unmarshal([]byte(serviceJsonStr), &service); err != nil {
				c.String(http.StatusOK, err.Error())
			}
		}

		log("dns service: %+v", service)
		servers[0].Locations = append(servers[0].Locations, nginx.Location{
			Dest: service.Spec.ClusterIP + ":" + p.Service.Port,
			Addr: p.Path,
		})
	}
	if err := nm.Apply(servers); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	// Step 2: Start nginx container
	log("nginx starting..")
	if err = nm.Start(); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	// Step 3: Attach ip to nginx container
	var nginxIp string
	if nginxIp, err = helper.NewPodIp(); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if err = weaveutil.WeaveAttach(nm.GetName(), nginxIp+url.MaskStr); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	log("%#v", nginxIp)

	// Step 4: Modify dns configuration
	if err = dns2.New(path.Join(url.DNSDirPath, url.DNSHostsFileName)).AddEntry(dns.Spec.Host, nginxIp); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	// Step 5: Store to etcd
	if dnsJson, err := json.Marshal(dns); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	} else {
		if err := etcd.Put(path.Join(url.DNSURL, dns.Metadata.Namespace, dns.Metadata.Name), string(dnsJson)); err != nil {
			c.String(http.StatusOK, err.Error())
			return
		}
	}
	c.String(http.StatusOK, "Apply successfully!")
}

func HandleApplyGpuJob(c *gin.Context) {
	gpu := apiObject.GpuJob{}
	err := readAndUnmarshal(c.Request.Body, &gpu)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	gpu.Metadata.UID = uidutil.New()
	log("receive gpu job[ID = %v]: %v", gpu.UID(), gpu)

	etcdURL := path.Join(url.GpuURL, gpu.Namespace(), gpu.Name())
	if raw, err := etcd.Get(etcdURL); err == nil {
		old := apiObject.GpuJob{}
		if err := json.Unmarshal([]byte(raw), &old); err == nil {
			c.String(http.StatusOK, fmt.Sprintf("gpu job %s/%s already exists", gpu.Namespace(), gpu.Name()))
			return
		}
	}

	gpuJobJson, _ := json.Marshal(gpu)
	if err = etcd.Put(etcdURL, string(gpuJobJson)); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	GpuUpdateMsg, _ := json.Marshal(entity.GpuUpdate{
		Action: entity.CreateAction,
		Target: gpu,
	})

	listwatch.Publish(topicutil.GpuJobUpdateTopic(), GpuUpdateMsg)
	c.String(http.StatusOK, "ok")
}
