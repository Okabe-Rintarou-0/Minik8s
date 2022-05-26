package kpa

import (
	"context"
	"minik8s/listwatch"
	"minik8s/util/logger"
	"minik8s/util/topicutil"
	"minik8s/util/wait"
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
}

func NewController() Controller {
	return &controller{}
}
