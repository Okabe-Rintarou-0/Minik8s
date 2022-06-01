package entity

import (
	"encoding/json"
	"minik8s/apiObject"
	"time"
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

type FunctionTriggerStatus string

const (
	TriggerUnknown     FunctionTriggerStatus = "Unknown"
	TriggerSucceed     FunctionTriggerStatus = "Succeed"
	TriggerError       FunctionTriggerStatus = "Error"
	TriggerInterrupted FunctionTriggerStatus = "Interrupted"
)

type FunctionTriggerResult struct {
	WorkflowNamespace string
	WorkflowName      string
	Data              FunctionData
	Time              time.Time
	Status            FunctionTriggerStatus
	Error             string
	FinishedAll       bool
}

type FunctionStatus struct {
	Name      string
	Instances int
	CodePath  string
}
