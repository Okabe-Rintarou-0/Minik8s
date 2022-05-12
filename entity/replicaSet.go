package entity

import (
	"minik8s/apiObject"
	"minik8s/apiObject/types"
	"time"
)

type ReplicaSetLifecycle byte

const (
	ReplicaSetScaling ReplicaSetLifecycle = iota
	ReplicaSetReady
	ReplicaSetError
	ReplicaSetDeleted
)

func (rsl *ReplicaSetLifecycle) String() string {
	switch *rsl {
	case ReplicaSetScaling:
		return "Scaling"
	case ReplicaSetReady:
		return "Ready"
	case ReplicaSetError:
		return "Error"
	case ReplicaSetDeleted:
		return "Deleted"
	}
	return "Unknown"
}

type ReplicaSetUpdate struct {
	Action ApiObjectUpdateAction
	Target apiObject.ReplicaSet
}

type ReplicaSetStatus struct {
	ID          types.UID
	Name        string
	Namespace   string
	Labels      apiObject.Labels
	Lifecycle   ReplicaSetLifecycle
	NumReplicas int
	NumReady    int
	CpuPercent  float64
	MemPercent  float64
	Error       string
	SyncTime    time.Time
}

func (rss *ReplicaSetStatus) FullName() string {
	return rss.Name + "_" + rss.Namespace
}
