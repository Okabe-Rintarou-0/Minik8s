package kpa

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"minik8s/apiObject"
	"minik8s/entity"
	"minik8s/listwatch"
	"minik8s/serverless/src/function"
	"minik8s/util/logger"
	"minik8s/util/topicutil"
	"minik8s/util/wait"
	"sync"
	"time"
)

var logManager = logger.Log("KPA manager")

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
	scaleLock             sync.RWMutex
	functionReplicaSetMap map[string]*functionReplicaSet
}

func (c *controller) createFunction(apiFunc *apiObject.Function) error {
	if _, exists := c.functionReplicaSetMap[apiFunc.Name]; !exists {
		if err := function.InitFunction(apiFunc.Name, apiFunc.Name); err != nil {
			return err
		}
	}
	return nil
}

func (c *controller) handleFunctionUpdate(msg *redis.Message) {
	functionUpdate := entity.FunctionUpdate{}
	_ = json.Unmarshal([]byte(msg.Payload), &functionUpdate)

	apiFunc := functionUpdate.Target

	var err error
	switch functionUpdate.Action {
	case entity.CreateAction:
		err = c.createFunction(&apiFunc)
	}

	if err != nil {
		logger.Error(err.Error())
	}
}

func (c *controller) Run() {
	go listwatch.Watch(topicutil.FunctionUpdateTopic(), c.handleFunctionUpdate)

	go wait.Period(scalePeriod, scalePeriod, c.scale)
	wait.Forever()
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
}

func NewController() Controller {
	return &controller{}
}
