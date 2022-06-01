package entity

import (
	"minik8s/apiObject"
	"time"
)

type GpuUpdate struct {
	Action ApiObjectUpdateAction
	Target apiObject.GpuJob
}

type GpuJobStatus struct {
	Namespace    string
	Name         string
	State        string
	LastSyncTime time.Time
}
