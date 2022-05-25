package kpa

import (
	"minik8s/controller/src/cache"
	"minik8s/util/logger"
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
	cacheManager          cache.Manager
	scaleLock             sync.RWMutex
	functionReplicaSetMap map[string]functionReplicaSet
}

func (c *controller) Run() {
	wait.Period(scalePeriod, scalePeriod, c.scale)
}

func (c *controller) scale() {
	c.scaleLock.Lock()
	for _, funcReplicaSet := range c.functionReplicaSetMap {
		now := time.Now()
		if funcReplicaSet.NumRequest == 0 && funcReplicaSet.NumReplicas > 0 && now.Sub(funcReplicaSet.LastRequestTime) > scaleTimeThreshold {
			c.scaleToHalf(&funcReplicaSet)
		}
	}
	c.scaleLock.Unlock()
}

type Controller interface {
	Run()
}

func NewController(cacheManager cache.Manager) Controller {
	return &controller{
		cacheManager: cacheManager,
	}
}
