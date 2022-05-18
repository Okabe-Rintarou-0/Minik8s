package node

import (
	"minik8s/apiserver/src/url"
	"minik8s/controller/src/cache"
	"minik8s/entity"
	"minik8s/util/httputil"
	"minik8s/util/logger"
	"minik8s/util/wait"
	"path"
	"strconv"
	"strings"
	"time"
)

const syncPeriod = time.Second * 5
const unhealthyTime = time.Second * 40
const unknownTime = time.Minute * 2
const deleteTime = time.Minute * 4

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

	URL := url.Prefix + path.Join(url.NodeURL, "status", status.Namespace, status.Hostname)
	_ = httputil.PutForm(URL, map[string]string{
		"lifecycle": strconv.Itoa(int(lifecycle)),
	})

	//log("Sync with api-server and get resp: %s", resp)
}

func (c *controller) deleteNodeAndAllNodePods(node, namespace, name string) {
	deleteNodeURL := url.Prefix + path.Join(url.NodeURL, namespace, name)
	httputil.DeleteWithoutBody(deleteNodeURL)
	deleteNodePodsURL := url.Prefix + strings.Replace(url.PodURLWithSpecifiedNode, ":node", node, 1)
	httputil.DeleteWithoutBody(deleteNodePodsURL)
}

func (c *controller) syncLoopIteration() {
	// Step 1: Get Node Statuses
	nodeStatuses := c.getNodeStatuses()

	// Step 2: Compute time delta, if not ready or unknown, sync with api-server
	for _, nodeStatus := range nodeStatuses {
		//log("Check node status[hostname = %v, lifecycle = %s]", nodeStatus.Hostname, nodeStatus.Lifecycle.String())
		timeDelta := time.Now().Sub(nodeStatus.SyncTime)
		//log("Time delta: %v", timeDelta.Seconds())
		fullName := path.Join(nodeStatus.Namespace, nodeStatus.Hostname)
		switch {
		case timeDelta >= deleteTime:
			log("Surpass delete time, should delete all the pods on the node %s", nodeStatus.Hostname)
			c.deleteNodeAndAllNodePods(nodeStatus.Hostname, nodeStatus.Namespace, nodeStatus.Hostname)
			c.cacheManager.DeleteNodeStatus(fullName)
		case timeDelta >= unknownTime:
			if nodeStatus.Lifecycle != entity.NodeUnknown {
				nodeStatus.Lifecycle = entity.NodeUnknown
				c.cacheManager.SetNodeStatus(fullName, nodeStatus)
				c.syncNodeStatusWithApiServer(nodeStatus)
			}
		case timeDelta >= unhealthyTime:
			if nodeStatus.Lifecycle != entity.NodeNotReady {
				nodeStatus.Lifecycle = entity.NodeNotReady
				c.cacheManager.SetNodeStatus(fullName, nodeStatus)
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
