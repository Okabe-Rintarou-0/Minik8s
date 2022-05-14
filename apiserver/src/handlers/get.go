package handlers

import (
	"github.com/gin-gonic/gin"
	"minik8s/entity"
	"minik8s/util/uidutil"
	"net/http"
	"os"
	"time"
)

// For test only
func getTestNodes() []*entity.NodeStatus {
	var nodes []*entity.NodeStatus
	node := &entity.NodeStatus{}
	node.Ip = "192.168.1.103"
	hostname, _ := os.Hostname()
	node.Hostname = hostname
	node.SyncTime = time.Now()
	node.Labels = map[string]string{
		"os": "windows",
	}

	node2 := &entity.NodeStatus{}
	node2.Ip = "192.168.1.103"
	node2.Hostname = "example"
	node2.SyncTime = time.Now()
	node2.Labels = map[string]string{
		"os": "linux",
	}

	nodes = append(nodes, node)
	nodes = append(nodes, node2)
	return nodes
}

// For test only
func getTestNode(name string) *entity.NodeStatus {
	for _, node := range getTestNodes() {
		if node.Hostname == name {
			return node
		}
	}
	return nil
}

func HandleGetNodes(c *gin.Context) {
	c.JSON(http.StatusOK, getTestNodes())
}

func HandleGetNode(c *gin.Context) {
	name := c.Param("name")
	c.JSON(http.StatusOK, getTestNode(name))
}

func podStatusForTest() *entity.PodStatus {
	return &entity.PodStatus{
		ID:        uidutil.New(),
		Name:      "example",
		Labels:    nil,
		Namespace: "default",
		Lifecycle: entity.PodContainerCreating,
		SyncTime:  time.Now(),
	}
}

func getPodStatusesForTest() []*entity.PodStatus {
	pod1 := podStatusForTest()
	pod2 := podStatusForTest()
	pod2.Name = "Example"
	return []*entity.PodStatus{pod1, pod2}
}

func getPodStatusForTest(name string) *entity.PodStatus {
	for _, pod := range getPodStatusesForTest() {
		if pod.Name == name {
			return pod
		}
	}
	return nil
}

func getPodDescriptionForTest(name string) *entity.PodDescription {
	var logs []entity.PodStatusLogEntry

	podStatus := podStatusForTest()
	logs = append(logs, entity.PodStatusLogEntry{
		Status: podStatus.Lifecycle,
		Time:   podStatus.SyncTime,
		Error:  podStatus.Error,
	})

	logs = append(logs, entity.PodStatusLogEntry{
		Status: entity.PodRunning,
		Time:   time.Now().Add(time.Minute * 30),
		Error:  "",
	})

	return &entity.PodDescription{
		CurrentStatus: *podStatusForTest(),
		Logs:          logs,
	}
}

func HandleGetPods(c *gin.Context) {
	c.JSON(http.StatusOK, getPodStatusesForTest())
}

func HandleGetPod(c *gin.Context) {
	name := c.Param("name")
	c.JSON(http.StatusOK, getPodStatusForTest(name))
}

func HandleDescribePod(c *gin.Context) {
	name := c.Param("name")
	c.JSON(http.StatusOK, getPodDescriptionForTest(name))
}
