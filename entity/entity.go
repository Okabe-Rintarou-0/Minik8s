package entity

import "minik8s/apiObject"

const (
	CreateAction PodUpdateAction = iota
	DeleteAction
	UpdateAction
)

type PodUpdateAction byte

type PodUpdate struct {
	Action PodUpdateAction
	Target apiObject.Pod
}
