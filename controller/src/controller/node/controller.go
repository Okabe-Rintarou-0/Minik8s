package node

import (
	"minik8s/apiserver/src/url"
	"minik8s/controller/src/cache"
	"minik8s/entity"
	"minik8s/util/httputil"
	"minik8s/util/logger"
	"minik8s/util/wait"
	"strconv"
	"strings"
	"time"
)

const syncPeriod = time.Second * 5
const unhealthyTime = time.Second * 40
const unknownTime = time.Minute * 2

var log = logger.Log("Node Controller")

type Controller interface {
	Run()
}

func NewController(cacheManager cache.Manager) Controller {
	return &controller{
		cacheManager: cacheManager,
	}
}

type controller struct {
	cacheManager cache.Manager
}

func (c *controller) getNodeStatuses() []*entity.NodeStatus {
	return c.cacheManager.GetNodeStatuses()
}

func (c *controller) syncNodeStatusWithApiServer(status *entity.NodeStatus) {
	lifecycle := status.Lifecycle
	if lifecycle != entity.NodeUnknown && lifecycle != entity.NodeNotReady {
		return
	}

	URL := url.Prefix + strings.Replace(url.NodeStatusURLWithSpecifiedName, ":name", status.Hostname, 1)
	_ = httputil.PostForm(URL, map[string]string{
		"lifecycle": strconv.Itoa(int(lifecycle)),
	})

	//log("Sync with api-server and get resp: %s", resp)
}

func (c *controller) syncLoopIteration() {
	// Step 1: Get Node Statuses
	nodeStatuses := c.getNodeStatuses()

	// Step 2: Compute time delta, if not ready or unknown, sync with api-server
	for _, nodeStatus := range nodeStatuses {
		//log("Check node status[hostname = %v, lifecycle = %s]", nodeStatus.Hostname, nodeStatus.Lifecycle.String())
		timeDelta := time.Now().Sub(nodeStatus.SyncTime)
		//log("Time delta: %v", timeDelta.Seconds())
		switch {
		case timeDelta >= unknownTime:
			if nodeStatus.Lifecycle != entity.NodeUnknown {
				nodeStatus.Lifecycle = entity.NodeUnknown
				c.cacheManager.SetNodeStatus(nodeStatus.Hostname, nodeStatus)
				c.syncNodeStatusWithApiServer(nodeStatus)
			}
		case timeDelta >= unhealthyTime:
			if nodeStatus.Lifecycle != entity.NodeNotReady {
				nodeStatus.Lifecycle = entity.NodeNotReady
				c.cacheManager.SetNodeStatus(nodeStatus.Hostname, nodeStatus)
				c.syncNodeStatusWithApiServer(nodeStatus)
			}
		}
	}
}

func (c *controller) syncLoop() {
	wait.Period(syncPeriod, syncPeriod, c.syncLoopIteration)
}

func (c *controller) Run() {
	c.syncLoop()
}
