package kpa

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/serverless/src/trigger"
	"minik8s/util/logger"
	"minik8s/util/topicutil"
	"minik8s/util/wait"
	"net/http"
	"sync"
	"time"
)

var logManager = logger.Log("KPA manager")
var bgCtx = context.Background()

const (
	scalePeriod        = time.Second * 10
	scaleTimeThreshold = time.Minute * 2
)

type functionReplicaSet struct {
	Function        string
	NumRequest      int
	NumReplicas     int
	LastRequestTime time.Time
}

type controller struct {
	workers               map[string]*worker
	workerLock            sync.Mutex
	scaleLock             sync.RWMutex
	functionReplicaSetMap map[string]*functionReplicaSet
}

func (c *controller) HandleTriggerFunc(ctx *gin.Context) {
	funcName := ctx.Param("function")
	dataRaw, _ := ioutil.ReadAll(ctx.Request.Body)
	result, err := c.TriggerFunc(funcName, entity.FunctionData(dataRaw))
	msg := entity.FunctionMsg{}
	if err != nil {
		msg.Error = err.Error()
		msg.Status = entity.TriggerError
	} else {
		msg.Data = result
		msg.Status = entity.TriggerSucceed
	}

	ctx.JSON(http.StatusOK, msg)
	msgJson, _ := json.Marshal(msg)
	logManager("Publish msg: %+v", msg)
	listwatch.Publish(topicutil.FunctionTriggerTopic(), msgJson)
}

func (c *controller) TriggerFunc(funcName string, data entity.FunctionData) (result entity.FunctionData, err error) {
	c.scaleLock.RLock()
	replicaSet := c.functionReplicaSetMap[funcName]
	c.scaleLock.RUnlock()
	if replicaSet == nil {
		return "", fmt.Errorf("no such function %s", funcName)
	}
	c.scaleLock.Lock()
	replicaSet.NumRequest += 1
	if replicaSet.NumRequest > replicaSet.NumReplicas {
		replicaSet.NumReplicas = replicaSet.NumRequest
		replicaSet.LastRequestTime = time.Now()
		c.scaleLock.Unlock()
		logManager("Received func %s trigger, scale to %d", funcName, replicaSet.NumReplicas)
		c.updateReplicaSetToApiServer(replicaSet)
	} else {
		c.scaleLock.Unlock()
	}

	result, err = trigger.Trigger(funcName, data)
	c.scaleLock.Lock()
	defer c.scaleLock.Unlock()
	replicaSet.NumRequest -= 1
	replicaSet.LastRequestTime = time.Now()
	logManager("func %s trigger ended, num request becomes to %d", funcName, replicaSet.NumRequest)

	return
}

func (c *controller) Run() {
	go listwatch.Watch(topicutil.FunctionUpdateTopic(), c.handleFunctionUpdate)
	go listwatch.Watch(topicutil.WorkflowUpdateTopic(), c.handleWorkflowUpdate)
	go wait.Period(scalePeriod, scalePeriod, c.scale)
}

func (c *controller) scale() {
	c.scaleLock.Lock()
	for _, funcReplicaSet := range c.functionReplicaSetMap {
		now := time.Now()
		if funcReplicaSet.NumRequest == 0 && funcReplicaSet.NumReplicas > 0 && now.Sub(funcReplicaSet.LastRequestTime) > scaleTimeThreshold {
			c.scaleToHalf(funcReplicaSet)
		}
	}
	c.scaleLock.Unlock()
}

type Controller interface {
	Run()
	HandleTriggerFunc(c *gin.Context)
}

func NewController() Controller {
	return &controller{
		workers:               make(map[string]*worker),
		functionReplicaSetMap: make(map[string]*functionReplicaSet),
	}
}
