package entity

import (
	"encoding/json"
	"minik8s/apiObject"
)

type FunctionUpdate struct {
	Action ApiObjectUpdateAction
	Target apiObject.Function
}

type FunctionData string

func NewFunctionData(data map[string]interface{}) FunctionData {
	jsonRaw, _ := json.Marshal(data)
	return FunctionData(jsonRaw)
}

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
