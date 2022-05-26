package entity

import "minik8s/apiObject"

type FunctionUpdate struct {
	Action ApiObjectUpdateAction
	Target apiObject.Function
}

type FunctionData map[string]interface{}

type FunctionStatus string

const (
	Succeed     FunctionStatus = "Succeed"
	Error       FunctionStatus = "Error"
	Interrupted FunctionStatus = "Interrupted"
)

type FunctionMsg struct {
	Data   FunctionData
	Status FunctionStatus
	Error  string
}
