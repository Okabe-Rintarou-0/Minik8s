package entity

import (
	"minik8s/apiObject"
	"time"
)

const (
	CreateAction ApiObjectUpdateAction = iota
	DeleteAction
	UpdateAction
)

func (action *ApiObjectUpdateAction) String() string {
	switch *action {
	case CreateAction:
		return "Create"
	case DeleteAction:
		return "Delete"
	case UpdateAction:
		return "Update"
	}
	return "Unknown"
}

type ApiObjectUpdateAction byte

type PodUpdate struct {
	Action ApiObjectUpdateAction
	Target apiObject.Pod
}

type ReplicaSetUpdate struct {
	Action ApiObjectUpdateAction
	Target apiObject.ReplicaSet
}

type Status byte

const (
	ContainerCreating Status = iota
	Error
	Running
	Deleted
)

func (s *Status) String() string {
	switch *s {
	case ContainerCreating:
		return "ContainerCreating"
	case Error:
		return "Error"
	case Running:
		return "Running"
	case Deleted:
		return "Deleted"
	}
	return "Unknown"
}

type PodStatus struct {
	ID         string
	Name       string
	Labels     apiObject.Labels
	Namespace  string
	Status     Status
	CpuPercent float64
	MemPercent float64
	Error      string
	SyncTime   time.Time
}

type PodStatusLogEntry struct {
	Status Status
	Time   time.Time
	Error  string
}

type PodDescription struct {
	CurrentStatus PodStatus
	Logs          []PodStatusLogEntry
}
