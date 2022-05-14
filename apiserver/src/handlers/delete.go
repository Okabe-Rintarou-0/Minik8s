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
)

func HandleDeletePod(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer c.Request.Body.Close()
	pod := &apiObject.Pod{}
	if err = json.Unmarshal(body, pod); err != nil {
		logger.Error(err.Error())
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
}

func HandleDeleteReplicaSet(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer c.Request.Body.Close()
	rs := &apiObject.ReplicaSet{}
	if err = json.Unmarshal(body, rs); err != nil {
		logger.Error(err.Error())
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
}

func HandleDeleteHPA(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer c.Request.Body.Close()
	hpa := &apiObject.HorizontalPodAutoscaler{}
	if err = json.Unmarshal(body, hpa); err != nil {
		logger.Error(err.Error())
		return
	}
	hpa.Metadata.UID = uidutil.New()
	log("receive hpa[ID = %v]: %v", hpa.UID(), hpa)
	var msg []byte
	msg, err = json.Marshal(entity.HPAUpdate{
		Action: entity.CreateAction,
		Target: *hpa,
	})
	listwatch.Publish(topicutil.HPAUpdateTopic(), msg)
}
