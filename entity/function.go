package entity

import "minik8s/apiObject"

type FunctionUpdate struct {
	Action ApiObjectUpdateAction
	Target apiObject.Function
}
