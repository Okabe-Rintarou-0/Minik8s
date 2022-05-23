package entity

import "minik8s/apiObject"

type GpuUpdate struct {
	Action ApiObjectUpdateAction
	Target apiObject.GpuJob
}
