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
	etcdURL := path.Join(url.PodURL, pod.Namespace(), pod.Name())
	if podJsonStr, err := etcd.Get(etcdURL); err == nil {
		getPod := &apiObject.Pod{}
		if err = json.Unmarshal([]byte(podJsonStr), getPod); err == nil {
			c.String(http.StatusOK, fmt.Sprintf("replicaSet %s/%s already exists", getPod.Namespace(), getPod.Name()))
			return
		}
	}

	if err = etcd.Put(etcdURL, string(podJson)); err != nil {
		c.String(http.StatusOK, err.Error())
		return
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
	etcdURL := path.Join(url.ReplicaSetURL, rs.Namespace(), rs.Name())
	if rsJsonStr, err := etcd.Get(etcdURL); err == nil {
		getRs := &apiObject.ReplicaSet{}
		if err = json.Unmarshal([]byte(rsJsonStr), getRs); err == nil {
			c.String(http.StatusOK, fmt.Sprintf("replicaSet %s/%s already exists", getRs.Namespace(), getRs.Name()))
			return
		}
	}

	if err = etcd.Put(etcdURL, string(rsJson)); err != nil {
		c.String(http.StatusOK, err.Error())
		return
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
