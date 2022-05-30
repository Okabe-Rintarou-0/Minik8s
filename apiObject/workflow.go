package apiObject

import "minik8s/apiObject/types"

type NodeType string

const (
	ChoiceNode NodeType = "Choice"
	TaskNode   NodeType = "Task"
)

type Choice struct {
	Variable           string  `json:"variable"`
	Next               *string `json:"next"`

	NumericEquals      *int64  `json:"numericEquals,omitempty"`
	NumericNotEquals   *int64  `json:"numericNotEquals,omitempty"`
	NumericLessThan    *int64  `json:"numericLessThan,omitempty"`
	NumericGreaterThan *int64  `json:"numericGreaterThan,omitempty"`

	BooleanEquals      *bool   `json:"booleanEquals,omitempty"`

	StringEquals       *string `json:"stringEquals,omitempty"`
	StringNotEquals    *string `json:"stringNotEquals,omitempty"`
	StringLessThan     *string `json:"stringLessThan,omitempty"`
	StringGreaterThan  *string `json:"stringGreaterThan,omitempty"`
}

type Choices struct {
	Choices []Choice `json:"choices"`
}

type WorkflowNode struct {
	Type     NodeType `json:"type"`
	*Task    `json:",inline"`
	*Choices `json:",inline"`
}

type Task struct {
	Next *string `json:"next"`
}

type Workflow struct {
	Base    `json:",inline"`
	StartAt string                  `json:"startAt"`
	Params  map[string]interface{}  `json:"params"`
	Nodes   map[string]WorkflowNode `json:"nodes"`
}

func (wf *Workflow) UID() types.UID {
	return wf.Metadata.UID
}

func (wf *Workflow) Namespace() string {
	return wf.Metadata.Namespace
}

func (wf *Workflow) Name() string {
	return wf.Metadata.Name
}
