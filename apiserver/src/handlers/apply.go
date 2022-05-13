package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"minik8s/apiObject"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/util/logger"
	"minik8s/util/topicutil"
	"minik8s/util/uidutil"
	"net/http"
)

var log = logger.Log("Api-server")

func HandleApplyPod(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	pod := &apiObject.Pod{}
	if err = json.Unmarshal(body, pod); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	pod.Metadata.UID = uidutil.New()
	log("receive pod[ID = %v]: %v", pod.UID(), pod)
	var msg []byte
	msg, err = json.Marshal(entity.PodUpdate{
		Action: entity.CreateAction,
		Target: *pod,
	})
	listwatch.Publish(topicutil.SchedulerPodUpdateTopic(), msg)
	c.String(http.StatusOK, "Apply successfully!")
}

func HandleApplyReplicaSet(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	defer c.Request.Body.Close()
	rs := &apiObject.ReplicaSet{}
	if err = json.Unmarshal(body, rs); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	rs.Metadata.UID = uidutil.New()
	log("receive rs[ID = %v]: %v", rs.UID(), rs)
	var msg []byte
	msg, err = json.Marshal(entity.ReplicaSetUpdate{
		Action: entity.CreateAction,
		Target: *rs,
	})
	listwatch.Publish(topicutil.ReplicaSetUpdateTopic(), msg)
	c.String(http.StatusOK, "Apply successfully!")
}

func HandleApplyHPA(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	defer c.Request.Body.Close()
	hpa := &apiObject.HorizontalPodAutoscaler{}
	if err = json.Unmarshal(body, hpa); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if err = addHPA(hpa); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	c.String(http.StatusOK, "Apply successfully!")
}
