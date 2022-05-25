package entity

import (
	"minik8s/apiObject"
	"minik8s/apiObject/types"
	"minik8s/kubelet/src/runtime/container"
	"time"
)

type PodLifecycle byte

const (
	PodContainerCreating PodLifecycle = iota
	PodCreated
	PodError
	PodRunning
	PodDeleted
	PodUnknown
	PodScheduled
)

func (pl *PodLifecycle) String() string {
	switch *pl {
	case PodScheduled:
		return "Scheduled"
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
	ID           types.UID
	Node         string
	Name         string
	Namespace    string
	Labels       apiObject.Labels
	Lifecycle    PodLifecycle
	CpuPercent   float64
	MemPercent   float64
	Error        string
	PortBindings container.PortBindings
	SyncTime     time.Time
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
