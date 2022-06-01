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
	"minik8s/util/httputil"
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

func deleteSpecifiedPod(namespace, name string) (pod *apiObject.Pod, node string, err error) {
	log("Pod to delete is %s/%s", namespace, name)

	etcdPodStatusURL := path.Join(url.PodURL, "status", namespace, name)
	_ = etcd.Delete(etcdPodStatusURL)

	var raw string
	nodes := helper.GetNodeHostnames()
	for _, node = range nodes {
		etcdPodURL := path.Join(url.PodURL, node, namespace, name)
		raw, err = etcd.Get(etcdPodURL)
		if err != nil || raw == "" {
			continue
		}
		if err = json.Unmarshal([]byte(raw), &pod); err != nil {
			return nil, "", err
		}

		// Delete endpoints
		// @TODO push to proxy
		log("Pod to delete is %+v", *pod)
		err = helper.DelEndpoints(*pod)
		if err = etcd.Delete(etcdPodURL); err == nil {
			log("Delete pod %s/%s successfully", namespace, name)
			break
		}
	}
	return
}

func deleteSpecifiedReplicaSet(namespace, name string) (rs *apiObject.ReplicaSet, err error) {
	log("Rs to delete is %s/%s", namespace, name)

	etcdReplicaSetStatusURL := path.Join(url.ReplicaSetURL, "status", namespace, name)
	_ = etcd.Delete(etcdReplicaSetStatusURL)

	var raw string
	etcdReplicaSetURL := path.Join(url.ReplicaSetURL, namespace, name)
	if raw, err = etcd.Get(etcdReplicaSetURL); err != nil {
		return nil, err
	}

	if err = json.Unmarshal([]byte(raw), &rs); err != nil {
		return nil, fmt.Errorf("no such replicaSet %s/%s", namespace, name)
	}

	err = etcd.Delete(etcdReplicaSetURL)
	return
}

func deleteSpecifiedGpuJob(namespace, name string) (gpu *apiObject.GpuJob, err error) {
	log("gpu to delete is %s/%s", namespace, name)

	var raw string

	_ = etcd.Delete(path.Join(url.GpuURL, "status", namespace, name))

	etcdURL := path.Join(url.GpuURL, namespace, name)
	if raw, err = etcd.Get(etcdURL); err != nil {
		return nil, err
	}

	if err = json.Unmarshal([]byte(raw), &gpu); err != nil {
		return nil, fmt.Errorf("no such gpu job %s/%s", namespace, name)
	}

	err = etcd.Delete(etcdURL)
	return
}

func deleteSpecifiedHPA(namespace, name string) (hpa *apiObject.HorizontalPodAutoscaler, err error) {
	log("hpa to delete is %s/%s", namespace, name)

	etcdHPAStatusURL := path.Join(url.HPAURL, "status", namespace, name)
	_ = etcd.Delete(etcdHPAStatusURL)

	var raw string
	etcdHPAURL := path.Join(url.HPAURL, namespace, name)
	if raw, err = etcd.Get(etcdHPAURL); err != nil {
		return nil, err
	}

	if err = json.Unmarshal([]byte(raw), &hpa); err != nil {
		return nil, fmt.Errorf("no such hpa %s/%s", namespace, name)
	}

	err = etcd.Delete(etcdHPAURL)
	return
}

func createAndPublishPodDeleteMsg(node string, pod *apiObject.Pod) {
	podDeleteMsg, _ := json.Marshal(entity.PodUpdate{
		Action: entity.DeleteAction,
		Target: *pod,
	})
	listwatch.Publish(topicutil.PodUpdateTopic(node), podDeleteMsg)
}

func HandleDeleteNode(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	if err := deleteSpecifiedNode(namespace, name); err != nil {
		c.String(http.StatusOK, err.Error())
	}
	c.String(http.StatusOK, "ok")
}

func deletePod(namespace, name string) error {
	if podToDelete, node, err := deleteSpecifiedPod(namespace, name); err == nil {
		if podToDelete != nil {
			createAndPublishPodDeleteMsg(node, podToDelete)
			return nil
		} else {
			return fmt.Errorf("no such pod %s/%s", namespace, name)
		}
	} else {
		return err
	}
}

func HandleDeletePod(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	if err := deletePod(namespace, name); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	c.String(http.StatusOK, "ok")
}

func HandleDeleteReplicaSet(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	if replicaSetToDelete, err := deleteSpecifiedReplicaSet(namespace, name); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	} else {
		replicaSetDeleteMsg, _ := json.Marshal(entity.ReplicaSetUpdate{
			Action: entity.DeleteAction,
			Target: *replicaSetToDelete,
		})
		listwatch.Publish(topicutil.ReplicaSetUpdateTopic(), replicaSetDeleteMsg)
	}
	c.String(http.StatusOK, "ok")
}

func HandleDeleteGpuJob(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	if _, err := deleteSpecifiedGpuJob(namespace, name); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	c.String(http.StatusOK, "ok")
}

func HandleDeleteHPA(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	if hpaToDelete, err := deleteSpecifiedHPA(namespace, name); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	} else {
		hpaDeleteMsg, _ := json.Marshal(entity.HPAUpdate{
			Action: entity.DeleteAction,
			Target: *hpaToDelete,
		})
		listwatch.Publish(topicutil.HPAUpdateTopic(), hpaDeleteMsg)
	}
	c.String(http.StatusOK, "ok")
}

func HandleReset(c *gin.Context) {
	if err := etcd.DeleteAllKeys(); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	c.String(http.StatusOK, "ok")
}

func HandleDeleteNodePods(c *gin.Context) {
	node := c.Param("node")
	pods := helper.GetPodsApiObjectFromEtcd(node)
	for _, pod := range pods {
		etcdPodStatusURL := path.Join(url.PodURL, "status", pod.Namespace(), pod.Name())
		_ = etcd.Delete(etcdPodStatusURL)

		etcdPodURL := path.Join(url.PodURL, node, pod.Namespace(), pod.Name())
		if err := etcd.Delete(etcdPodURL); err == nil {
			createAndPublishPodDeleteMsg(node, pod)
		}
	}
}

func HandleDeleteService(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	if !helper.ExistsService(namespace, name) {
		c.String(http.StatusOK, fmt.Sprintf("service %s/%s does not exist", namespace, name))
		return
	}

	service := apiObject.Service{}
	if serviceJsonStr, err := etcd.Get(path.Join(url.ServiceURL, namespace, name)); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	} else {
		if err := json.Unmarshal([]byte(serviceJsonStr), &service); err != nil {
			c.String(http.StatusOK, err.Error())
			return
		}
	}

	serviceUpdate := entity.ServiceUpdate{
		Action: entity.DeleteAction,
		Target: entity.ServiceTarget{
			Service:   service,
			Endpoints: make([]apiObject.Endpoint, 0),
		},
	}
	for key, value := range service.Spec.Selector {
		if endpoints, err := helper.GetEndpoints(key, value); err != nil {
			c.String(http.StatusOK, err.Error())
			return
		} else {
			serviceUpdate.Target.Endpoints = append(serviceUpdate.Target.Endpoints, endpoints...)
		}
	}

	serviceDeleteMsg, _ := json.Marshal(serviceUpdate)
	listwatch.Publish(topicutil.ServiceUpdateTopic(), serviceDeleteMsg)

	for key, value := range service.Spec.Selector {
		if err := etcd.Delete(path.Join(url.ServiceURL, key, value, service.Metadata.UID)); err != nil {
			c.String(http.StatusOK, err.Error())
			return
		}
	}
	if err := etcd.Delete(path.Join(url.ServiceURL, namespace, name)); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	c.String(http.StatusOK, "ok")
}

func HandleRemoveFunc(c *gin.Context) {
	name := c.Param("name")
	etcdURL := path.Join(url.FuncURL, name)
	apiFunc := apiObject.Function{}
	if raw, err := etcd.Get(etcdURL); err == nil {
		if err = json.Unmarshal([]byte(raw), &apiFunc); err == nil {
			if err = etcd.Delete(etcdURL); err != nil {
				c.String(http.StatusOK, err.Error())
				return
			}

			topic := topicutil.FunctionUpdateTopic()
			updateMsg, _ := json.Marshal(entity.FunctionUpdate{
				Action: entity.DeleteAction,
				Target: apiFunc,
			})

			log("Publish delete func %s msg", apiFunc.Name)
			listwatch.Publish(topic, updateMsg)
			c.String(http.StatusOK, "ok")
			return
		}
	}

	c.String(http.StatusOK, fmt.Sprintf("no such func %s", name))
	return
}

func HandleUpdateFunc(c *gin.Context) {
	name := c.Param("name")
	etcdURL := path.Join(url.FuncURL, name)
	apiFunc := apiObject.Function{}
	if raw, err := etcd.Get(etcdURL); err == nil {
		if err = json.Unmarshal([]byte(raw), &apiFunc); err == nil {
			if err = etcd.Delete(etcdURL); err != nil {
				c.String(http.StatusOK, err.Error())
				return
			}

			newFunc := apiObject.Function{}
			if err = httputil.ReadAndUnmarshal(c.Request.Body, &newFunc); err != nil {
				c.String(http.StatusOK, err.Error())
				return
			}

			topic := topicutil.FunctionUpdateTopic()
			updateMsg, _ := json.Marshal(entity.FunctionUpdate{
				Action: entity.UpdateAction,
				Target: newFunc,
			})

			log("Publish update func %s msg", newFunc.Name)
			listwatch.Publish(topic, updateMsg)
			c.String(http.StatusOK, "ok")
			return
		}
	}

	c.String(http.StatusOK, fmt.Sprintf("no such func %s", name))
	return
}

func HandleRemoveWorkflow(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	etcdWorkflowResultURL := path.Join(url.WorkflowURL, "result", namespace, name)
	_ = etcd.Delete(etcdWorkflowResultURL)

	etcdWorkflowURL := path.Join(url.WorkflowURL, namespace, name)
	wf := apiObject.Workflow{}
	raw, _ := etcd.Get(etcdWorkflowURL)
	if err := json.Unmarshal([]byte(raw), &wf); err != nil {
		c.String(http.StatusOK, fmt.Sprintf("no such workflow %s/%s", namespace, name))
		return
	}

	if err := etcd.Delete(etcdWorkflowURL); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	topic := topicutil.WorkflowUpdateTopic()
	updateMsg, _ := json.Marshal(entity.WorkflowUpdate{
		Action: entity.DeleteAction,
		Target: wf,
	})

	listwatch.Publish(topic, updateMsg)
	c.String(http.StatusOK, "ok")
}

func HandleDeleteDNS(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	if !helper.ExistsDNS(namespace, name) {
		c.String(http.StatusOK, fmt.Sprintf("dns %s/%s does not exist", namespace, name))
		return
	}

	dns := apiObject.Dns{}
	if dnsJsonStr, err := etcd.Get(path.Join(url.DNSURL, namespace, name)); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	} else {
		if err := json.Unmarshal([]byte(dnsJsonStr), &dns); err != nil {
			c.String(http.StatusOK, err.Error())
			return
		}
	}

	nm := nginx.New(dns.Metadata.UID)
	if err := nm.Shutdown(); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if err := dns2.New(path.Join(url.DNSDirPath, url.DNSHostsFileName)).DeleteIfExistEntry(dns.Spec.Host); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if err := etcd.Delete(path.Join(url.DNSURL, namespace, name)); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	c.String(http.StatusOK, "Delete successfully")
}
