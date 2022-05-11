package entity

import "minik8s/apiObject"

type HPAUpdate struct {
	Action ApiObjectUpdateAction
	Target apiObject.HorizontalPodAutoscaler
}
