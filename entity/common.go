package entity

import "minik8s/apiObject"

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
	Node   string
	Target apiObject.Pod
}

type EndpointUpdate struct {
	Action ApiObjectUpdateAction
	Target EndpointTarget
}

type ServiceUpdate struct {
	Action ApiObjectUpdateAction
	Target ServiceTarget
}
