package entity

import (
	"minik8s/apiObject"
	"minik8s/apiObject/types"
	"time"
)

type HPALifecycle byte
type HPAMetrics string

const (
	HPACreated HPALifecycle = iota
	HPAReady
	HPADeleted
)

const (
	CpuMetrics     HPAMetrics = "Cpu"
	MemoryMetrics  HPAMetrics = "Memory"
	UnknownMetrics HPAMetrics = "Unknown"
)

func (hl *HPALifecycle) String() string {
	switch *hl {
	case HPACreated:
		return "Created"
	case HPAReady:
		return "Ready"
	case HPADeleted:
		return "Deleted"
	default:
		return "Unknown"
	}
}

type HPAUpdate struct {
	Action ApiObjectUpdateAction
	Target apiObject.HorizontalPodAutoscaler
}

type HPAStatus struct {
	ID          types.UID
	Name        string
	Namespace   string
	Labels      apiObject.Labels
	Lifecycle   HPALifecycle
	MinReplicas int
	MaxReplicas int
	NumReady    int
	NumTarget   int
	Metrics     HPAMetrics
	Benchmark   float64
	Error       string
	SyncTime    time.Time
}
