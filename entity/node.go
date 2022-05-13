package entity

import (
	"minik8s/apiObject"
	"time"
)

type NodeLifecycle byte

const (
	NodeReady NodeLifecycle = iota
	NodeUnknown
	NodeDeleted
)

func (nl *NodeLifecycle) String() string {
	switch *nl {
	case NodeReady:
		return "Ready"
	case NodeDeleted:
		return "Deleted"
	default:
		return "Unknown"
	}
}

type NodeStatus struct {
	Hostname  string
	Ip        string
	Labels    apiObject.Labels
	Lifecycle NodeLifecycle
	Error     string
	SyncTime  time.Time
}
