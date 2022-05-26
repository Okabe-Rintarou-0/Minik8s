package entity

import "minik8s/apiObject"

type WorkflowUpdate struct {
	Action ApiObjectUpdateAction
	Target apiObject.Workflow
}
