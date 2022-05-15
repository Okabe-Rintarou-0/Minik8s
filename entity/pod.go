package entity

import (
	"minik8s/apiObject"
	"minik8s/apiObject/types"
	"time"
)

type PodLifecycle byte

const (
	PodContainerCreating PodLifecycle = iota
	PodCreated
	PodError
	PodRunning
	PodDeleted
)

func (pl *PodLifecycle) String() string {
	switch *pl {
	case PodContainerCreating:
		return "ContainerCreating"
	case PodError:
		return "Error"
	case PodCreated:
		return "Created"
	case PodRunning:
		return "Running"
	case PodDeleted:
		return "Deleted"
	}
	return "Unknown"
}

type PodStatus struct {
	ID         types.UID
	Name       string
	Namespace  string
	Labels     apiObject.Labels
	Lifecycle  PodLifecycle
	CpuPercent float64
	MemPercent float64
	Error      string
	SyncTime   time.Time
}

type PodStatusLogEntry struct {
	Status PodLifecycle
	Time   time.Time
	Error  string
}

type PodDescription struct {
	CurrentStatus PodStatus
	Logs          []PodStatusLogEntry
}
