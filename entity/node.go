package entity

import (
	"minik8s/apiObject"
	"time"
)

type NodeLifecycle byte

const (
	NodeReady NodeLifecycle = iota
	NodeUnknown
	NodeNotReady
	NodeCreated
	NodeDeleted
)

func (nl *NodeLifecycle) String() string {
	switch *nl {
	case NodeCreated:
		return "Created"
	case NodeReady:
		return "Ready"
	case NodeDeleted:
		return "Deleted"
	case NodeNotReady:
		return "NotReady"
	default:
		return "Unknown"
	}
}

type NodeStatus struct {
	Hostname   string
	Ip         string
	Labels     apiObject.Labels
	Lifecycle  NodeLifecycle
	Error      string
	CpuPercent float64
	MemPercent float64
	NumPods    int
	SyncTime   time.Time
}
