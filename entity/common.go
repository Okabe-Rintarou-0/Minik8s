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
	Target apiObject.Pod
}
